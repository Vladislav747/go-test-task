package title

import (
	"golang.org/x/net/html"
	"io"
	"log"
)

/*
Проверка что элемент title
*/
func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

/*
Ф-ция проход по html
*/
func traverse(n *html.Node) (string, bool) {
	if n == nil {
		return "", false
	}

	if n.FirstChild == nil {
		return "", false
	}

	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

/*
Главная ф-ция получения title из html
*/
func GetHtmlTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		log.Printf("failed to parse HTML Title")
	}

	return traverse(doc)
}
