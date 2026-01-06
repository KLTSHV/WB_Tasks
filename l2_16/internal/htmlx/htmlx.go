package htmlx

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"mirror/internal/urlmap"
)

type LinkKind int

const (
	KindPage LinkKind = iota
	KindResource
)

type DiscoveredLink struct {
	AbsoluteURL *url.URL
	Normalized  string
	Kind        LinkKind
}

// ParseRewriteAndCollect делает три вещи за один проход:
// 1) парсит HTML в DOM,
// 2) переписывает ссылки src/href на локальные относительные пути (только same-host),
// 3) собирает список найденных ссылок (страницы и ресурсы)
// depth/maxDepth нужны, чтобы не переписыват page-ссылки, которые мы не собираемся скачивать глубже
func ParseRewriteAndCollect(
	htmlBytes []byte,
	documentURL *url.URL,
	documentLocalPath string,
	mapper *urlmap.Mapper,
	depth int,
	maxDepth int,
) ([]byte, []DiscoveredLink, error) {
	root, err := html.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}
	// Определяем базовый URL для относительных ссылок
	baseURL := documentURL
	if baseHref := findBaseHref(root); baseHref != "" {
		if parsedBase, parseErr := documentURL.Parse(baseHref); parseErr == nil {
			baseURL = parsedBase
		}
	}

	var links []DiscoveredLink

	// Обходим DOM и обрабатываем нужные теги
	var walk func(node *html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "a":
				rewriteAttr(node, "href", KindPage, baseURL, documentLocalPath, mapper, &links, depth, maxDepth)
			case "img":
				rewriteAttr(node, "src", KindResource, baseURL, documentLocalPath, mapper, &links, depth, maxDepth)
			case "script":
				rewriteAttr(node, "src", KindResource, baseURL, documentLocalPath, mapper, &links, depth, maxDepth)
			case "link":
				rewriteAttr(node, "href", KindResource, baseURL, documentLocalPath, mapper, &links, depth, maxDepth)
			case "iframe":
				rewriteAttr(node, "src", KindResource, baseURL, documentLocalPath, mapper, &links, depth, maxDepth)
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(root)

	var out bytes.Buffer
	if err := html.Render(&out, root); err != nil {
		return nil, nil, fmt.Errorf("ошибка рендера HTML: %w", err)
	}

	return out.Bytes(), links, nil
}

func rewriteAttr(
	node *html.Node,
	attrKey string,
	kind LinkKind,
	baseURL *url.URL,
	documentLocalPath string,
	mapper *urlmap.Mapper,
	links *[]DiscoveredLink,
	depth int,
	maxDepth int,
) {
	for i := range node.Attr {
		if !strings.EqualFold(node.Attr[i].Key, attrKey) {
			continue
		}

		rawValue := strings.TrimSpace(node.Attr[i].Val)
		if rawValue == "" || isSkippable(rawValue) {
			return
		}

		absoluteURL, normalizedKey, err := mapper.Normalize(rawValue, baseURL)
		if err != nil || !mapper.Allowed(absoluteURL) {
			// Внешние ссылки не трогаем!!
			return
		}

		// Если это ссылка на страницу и мы уже на пределе глубины —
		// оставляем оригинальный URL, чтобы не вести пользователя на нескачанную локальную страницу
		if kind == KindPage && depth >= maxDepth {
			return
		}

		isPage := kind == KindPage

		localPath, err := mapper.LocalPath(absoluteURL, isPage)
		if err != nil {
			return
		}

		relativeLink, err := urlmap.RelativeFromFile(documentLocalPath, localPath)
		if err != nil {
			return
		}

		// Переписываем ссылку на локальную.
		node.Attr[i].Val = relativeLink

		*links = append(*links, DiscoveredLink{
			AbsoluteURL: absoluteURL,
			Normalized:  normalizedKey,
			Kind:        kind,
		})
		return
	}
}

func findBaseHref(root *html.Node) string {
	var href string

	var walk func(n *html.Node)
	walk = func(n *html.Node) {
		if href != "" {
			return
		}
		if n.Type == html.ElementNode && n.Data == "base" {
			for _, attr := range n.Attr {
				if strings.EqualFold(attr.Key, "href") {
					href = strings.TrimSpace(attr.Val)
					return
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(root)
	return href
}

func isSkippable(link string) bool {
	low := strings.ToLower(strings.TrimSpace(link))
	return strings.HasPrefix(low, "mailto:") ||
		strings.HasPrefix(low, "tel:") ||
		strings.HasPrefix(low, "javascript:") ||
		strings.HasPrefix(low, "data:")
}
