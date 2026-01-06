package crawl

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"

	"mirror/internal/fetch"
	"mirror/internal/htmlx"
	"mirror/internal/urlmap"
)

// Settings — настройки процесса зеркалирования.
type Settings struct {
	MaxDepth  int   // максимальная глубина рекурсии по страницам (a[href])
	Parallel  int   // количество воркеров
	MaxBytes  int64 // лимит размера одного ответа
	UserAgent string

	HTTPClient *http.Client
	URLMap     *urlmap.Mapper
}

// Job — одна задача на скачивание.
// IsPage=true: считаем URL страницей (HTML) и допускаем рекурсию по ссылкам внутри.
type Job struct {
	SourceURL  string
	Normalized string
	LocalPath  string

	Depth  int
	IsPage bool
}

func Run(ctx context.Context, settings Settings, startRawURL string) error {
	if settings.URLMap == nil || settings.HTTPClient == nil {
		return errors.New("не заданы URLMap и/или HTTPClient")
	}
	if settings.Parallel <= 0 {
		settings.Parallel = 4
	}

	startAbs, startKey, err := settings.URLMap.Normalize(startRawURL, nil)
	if err != nil {
		return err
	}
	if !settings.URLMap.Allowed(startAbs) {
		return errors.New("стартовый URL не проходит политику same-host")
	}

	startLocalPath, err := settings.URLMap.LocalPath(startAbs, true)
	if err != nil {
		return err
	}

	jobs := make(chan Job, 4096)

	// visited нужен, чтобы не скачивать один и тот же URL повторно и не зациклиться.
	var (
		visitedMu sync.Mutex
		visited   = make(map[string]string) // normalized URL -> local path
	)

	// inflight показывает, сколько задач сейчас в “обработке + очереди”.
	// Когда inflight = 0, значит больше нечего делать и можно закрывать канал.
	var inflight sync.WaitGroup

	enqueueIfNew := func(job Job) {
		visitedMu.Lock()
		if _, exists := visited[job.Normalized]; exists {
			visitedMu.Unlock()
			return
		}
		visited[job.Normalized] = job.LocalPath
		visitedMu.Unlock()

		inflight.Add(1)
		jobs <- job
	}

	// Стартовая задача — всегда “страница”.
	enqueueIfNew(Job{
		SourceURL:  startAbs.String(),
		Normalized: startKey,
		LocalPath:  startLocalPath,
		Depth:      0,
		IsPage:     true,
	})

	// Worker pool
	var workers sync.WaitGroup
	for i := 0; i < settings.Parallel; i++ {
		workers.Add(1)
		go func(workerID int) {
			defer workers.Done()
			for job := range jobs {
				processOne(ctx, settings, job, enqueueIfNew)
				inflight.Done()
			}
		}(i)
	}

	// Закрываем очередь, когда задач больше не осталось.
	go func() {
		inflight.Wait()
		close(jobs)
	}()

	workers.Wait()
	return nil
}

func processOne(ctx context.Context, settings Settings, job Job, enqueueIfNew func(Job)) {
	downloadResult, err := fetch.DownloadOne(
		ctx,
		settings.HTTPClient,
		job.SourceURL,
		job.LocalPath,
		settings.MaxBytes,
		settings.UserAgent,
	)
	if err != nil {
		log.Printf("FAIL %s: %v", job.SourceURL, err)
		return
	}

	log.Printf("OK   %s -> %s (%s, %d bytes)", job.SourceURL, job.LocalPath, downloadResult.ContentType, downloadResult.Bytes)

	// Ресурсы (картинки/css/js) мы просто сохраняем и не рекурсируем.
	if !job.IsPage {
		return
	}

	// Рекурсия делается только по HTML.
	if downloadResult.ContentType != "text/html" {
		return
	}

	// Берём сохранённый HTML с диска (это удобно: переписываем именно то, что реально скачали).
	htmlBytes, err := os.ReadFile(job.LocalPath)
	if err != nil {
		log.Printf("FAIL read %s: %v", job.LocalPath, err)
		return
	}

	documentURL, _, err := settings.URLMap.Normalize(downloadResult.FinalURL, nil)
	if err != nil {
		log.Printf("FAIL normalize final url %s: %v", downloadResult.FinalURL, err)
		return
	}

	rewrittenHTML, discoveredLinks, err := htmlx.ParseRewriteAndCollect(
		htmlBytes,
		documentURL,
		job.LocalPath,
		settings.URLMap,
		job.Depth,
		settings.MaxDepth,
	)
	if err != nil {
		log.Printf("FAIL parse/rewrite %s: %v", job.SourceURL, err)
		return
	}

	// Перезаписываем HTML уже с локальными ссылками.
	if err := os.WriteFile(job.LocalPath, rewrittenHTML, 0o644); err != nil {
		log.Printf("FAIL write %s: %v", job.LocalPath, err)
		return
	}

	// Добавляем найденные ссылки в очередь.
	for _, link := range discoveredLinks {
		isPage := link.Kind == htmlx.KindPage

		nextDepth := job.Depth
		if isPage {
			nextDepth = job.Depth + 1
			if nextDepth > settings.MaxDepth {
				continue
			}
		}

		localPath, err := settings.URLMap.LocalPath(link.AbsoluteURL, isPage)
		if err != nil {
			continue
		}

		enqueueIfNew(Job{
			SourceURL:  link.AbsoluteURL.String(),
			Normalized: link.Normalized,
			LocalPath:  localPath,
			Depth:      nextDepth,
			IsPage:     isPage,
		})
	}
}
