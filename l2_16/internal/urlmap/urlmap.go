package urlmap

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Mapper struct {
	StartURL  *url.URL // стартовый URL (как база для относительных ссылок)
	OutputDir string   // корневая папка, куда всё сохраняем
	startHost string   // домен стартового URL (same-host политика)
}

func New(startRawURL, outputDir string) (*Mapper, error) {
	parsed, err := url.Parse(startRawURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось распарсить стартовый URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, errors.New("поддерживаются только http/https")
	}
	if parsed.Host == "" {
		return nil, errors.New("в URL отсутствует host")
	}
	return &Mapper{
		StartURL:  parsed,
		OutputDir: outputDir,
		startHost: strings.ToLower(parsed.Hostname()),
	}, nil
}

// Allowed — простая политика: скачиваем только тот же домен, что у стартового URL.
func (m *Mapper) Allowed(absoluteURL *url.URL) bool {
	return absoluteURL != nil && strings.ToLower(absoluteURL.Hostname()) == m.startHost
}

// Normalize:
// - превращает относительную ссылку в абсолютную (через baseURL)
// - убирает #fragment
// - приводит host к lower-case
// - чистит путь (убирает ../ и ./)
// - возвращает нормализованный URL и строку-ключ для дедупликации
func (m *Mapper) Normalize(rawLink string, baseURL *url.URL) (*url.URL, string, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawLink))
	if err != nil {
		return nil, "", err
	}

	var absolute *url.URL
	if baseURL != nil {
		absolute = baseURL.ResolveReference(parsed)
	} else {
		absolute = m.StartURL.ResolveReference(parsed)
	}

	if absolute.Scheme != "http" && absolute.Scheme != "https" {
		return nil, "", fmt.Errorf("неподдерживаемая схема: %s", absolute.Scheme)
	}

	absolute.Fragment = ""
	absolute.Host = strings.ToLower(absolute.Host)

	absolute.Path = path.Clean(absolute.Path)
	if absolute.Path == "." {
		absolute.Path = "/"
	}

	normalizedKey := absolute.String()
	return absolute, normalizedKey, nil
}

// LocalPath строит детерминированный путь на диске для URL
// isPage=true: сохраняем как страницу (для /dir или /dir/ -> /dir/index.html)
// isPage=false: сохраняем как ресурс (оставляем расширения как есть)
func (m *Mapper) LocalPath(absoluteURL *url.URL, isPage bool) (string, error) {
	if absoluteURL == nil {
		return "", errors.New("absoluteURL = nil")
	}

	hostDir := sanitizeSegment(absoluteURL.Hostname())

	urlPath := absoluteURL.Path
	if urlPath == "" {
		urlPath = "/"
	}

	extension := path.Ext(urlPath)

	if isPage {
		// Страницы без расширения трактуем как директории: /blog -> /blog/index.html
		if strings.HasSuffix(urlPath, "/") || extension == "" {
			urlPath = strings.TrimSuffix(urlPath, "/")
			if urlPath == "" {
				urlPath = "/index.html"
			} else {
				urlPath = urlPath + "/index.html"
			}
		}
	} else {
		// Ресурсы: если вдруг путь оканчивается на '/', тоже сохраняем как index.html (редко, но безопасно).
		if strings.HasSuffix(urlPath, "/") {
			urlPath = strings.TrimSuffix(urlPath, "/")
			if urlPath == "" {
				urlPath = "/index.html"
			} else {
				urlPath = urlPath + "/index.html"
			}
		}
	}

	urlPath = strings.TrimPrefix(urlPath, "/")
	urlPath = sanitizePath(urlPath)

	// Если есть query — добавляем хэш, чтобы не перезатирать файлы (например /img?id=1 и /img?id=2).
	if absoluteURL.RawQuery != "" {
		queryHash := shortHash(absoluteURL.RawQuery)
		ext2 := path.Ext(urlPath)
		if ext2 != "" {
			base := strings.TrimSuffix(urlPath, ext2)
			urlPath = base + "__q" + queryHash + ext2
		} else {
			urlPath = urlPath + "__q" + queryHash
		}
	}

	return filepath.Join(m.OutputDir, hostDir, filepath.FromSlash(urlPath)), nil
}

// RelativeFromFile считает относительный путь от одного локального файла до другого
// Это нужно, чтобы в HTML писать "../assets/app.js", а не абсолютные пути с диска
func RelativeFromFile(fromFile, toFile string) (string, error) {
	fromDir := filepath.Dir(fromFile)
	rel, err := filepath.Rel(fromDir, toFile)
	if err != nil {
		return "", err
	}
	// В HTML всегда используем '/', даже на Windows.
	return filepath.ToSlash(rel), nil
}

func shortHash(text string) string {
	sum := sha1.Sum([]byte(text))
	return hex.EncodeToString(sum[:])[:10]
}

var reUnsafe = regexp.MustCompile(`[<>:"\\|?*\x00-\x1F]`)

func sanitizeSegment(segment string) string {
	segment = strings.TrimSpace(segment)
	segment = reUnsafe.ReplaceAllString(segment, "_")
	if segment == "" {
		return "_"
	}
	return segment
}

func sanitizePath(p string) string {
	parts := strings.Split(p, "/")
	out := make([]string, 0, len(parts))
	for _, seg := range parts {
		seg = sanitizeSegment(seg)
		if seg == "" || seg == "." || seg == ".." {
			continue
		}
		out = append(out, seg)
	}
	if len(out) == 0 {
		return "index.html"
	}
	return strings.Join(out, "/")
}
