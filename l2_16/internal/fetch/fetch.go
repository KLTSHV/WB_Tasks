package fetch

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DownloadResult struct {
	FinalURL    string // итоговый URL после редиректов
	ContentType string // чистый media type
	Bytes       int64  // сколько байт реально записали
	StatusCode  int
}

// Transport низкоуровневая часть HTTP-клиента (пулы соединений, keep-alive, лимиты)
func NewTransport(parallelism int) *http.Transport {
	if parallelism <= 0 {
		parallelism = 4
	}
	return &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: maxInt(10, parallelism),
		MaxConnsPerHost:     maxInt(10, parallelism),
		IdleConnTimeout:     30 * time.Second,
	}
}

// DownloadOne скачивает URL и сохраняет ответ в destPath
func DownloadOne(
	ctx context.Context,
	client *http.Client,
	sourceURL string,
	destPath string,
	maxBytes int64,
	userAgent string,
) (DownloadResult, error) {
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return DownloadResult{}, fmt.Errorf("не удалось создать запрос: %w", err)
	}
	if userAgent != "" {
		request.Header.Set("User-Agent", userAgent)
	}

	response, err := client.Do(request)
	if err != nil {
		return DownloadResult{}, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer response.Body.Close()

	result := DownloadResult{
		FinalURL:    response.Request.URL.String(),
		StatusCode:  response.StatusCode,
		ContentType: "",
		Bytes:       0,
	}

	// Зеркалирование имеет смысл только на успешных ответах
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		snippet, _ := io.ReadAll(io.LimitReader(response.Body, 2048))
		return result, fmt.Errorf("плохой статус %d: %s", response.StatusCode, strings.TrimSpace(string(snippet)))
	}

	// Достаём Content-Type и отбрасываем параметры  чтобы проще сравнивать
	contentTypeHeader := response.Header.Get("Content-Type")
	if contentTypeHeader != "" {
		mediaType, _, parseErr := mime.ParseMediaType(contentTypeHeader)
		if parseErr == nil {
			result.ContentType = mediaType
		} else {
			// На случай кривого заголовка сохраняем как есть
			result.ContentType = contentTypeHeader
		}
	}

	// Создаём директории заранее
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return result, fmt.Errorf("не удалось создать папки: %w", err)
	}
	// Пишем в временный файл и только потом переименовываем
	// так не останется битых файлов при ошибке в середине загрузки
	tempPath := destPath + ".part"
	file, err := os.Create(tempPath)
	if err != nil {
		return result, fmt.Errorf("не удалось создать файл: %w", err)
	}
	fileClosed := false
	defer func() {
		if !fileClosed {
			_ = file.Close()
		}
	}()

	var reader io.Reader = response.Body
	if maxBytes > 0 {
		reader = io.LimitReader(response.Body, maxBytes+1)
	}

	writtenBytes, err := io.Copy(file, reader)
	if err != nil {
		_ = file.Close()
		fileClosed = true
		_ = os.Remove(tempPath)
		return result, fmt.Errorf("ошибка записи тела ответа: %w", err)
	}

	if maxBytes > 0 && writtenBytes > maxBytes {
		_ = file.Close()
		fileClosed = true
		_ = os.Remove(tempPath)
		return result, errors.New("ответ слишком большой (превышен maxBytes)")
	}

	if err := file.Sync(); err != nil {
		_ = file.Close()
		fileClosed = true
		_ = os.Remove(tempPath)
		return result, fmt.Errorf("sync не удался: %w", err)
	}
	if err := file.Close(); err != nil {
		fileClosed = true
		_ = os.Remove(tempPath)
		return result, fmt.Errorf("close не удался: %w", err)
	}
	fileClosed = true

	if err := os.Rename(tempPath, destPath); err != nil {
		_ = os.Remove(tempPath)
		return result, fmt.Errorf("не удалось переименовать временный файл: %w", err)
	}

	result.Bytes = writtenBytes
	return result, nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
