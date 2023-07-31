package store

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func TraverseAndExtract(s *goquery.Selection) {
	nodeName := goquery.NodeName(s)

	// Ignore script and style tags
	if nodeName == "header" || nodeName == "footer" || nodeName == "nav" || nodeName == "aside" || nodeName == "comments" || nodeName == "script" || nodeName == "style" {
		return
	}

	//class, _ := s.Attr("class")
	//id, _ := s.Attr("id")

	switch nodeName {
	case "h1", "h2", "h3", "h4", "h5", "h6", "p", "li", "td", "span", "blockquote", "pre":
		fmt.Println(nodeName, ": ", s.Text())
	case "a":
		href, _ := s.Attr("href")
		fmt.Println(nodeName, ": ", href)
	case "img":
		src, exists := s.Attr("src")
		if !exists {
			// I want to handle lazy loading images too
			src, _ = s.Attr("data-src")
		}
		fmt.Println(nodeName, ": ", src)
	//case "article", "div", "section", "main":
	//	if containsAny(class, "content", "article", "post", "story", "news", "blog") ||
	//		containsAny(id, "content", "article", "post", "story", "news", "blog") {
	//		fmt.Println(nodeName, ": ", s.Text())
	//		// Recursively traverse children
	//		s.Children().Each(func(i int, child *goquery.Selection) {
	//			TraverseAndExtract(child)
	//		})
	//	}
	case "table":
		s.Find("tr").Each(func(i int, tr *goquery.Selection) {
			cells := ""
			tr.Find("th, td").Each(func(j int, cell *goquery.Selection) {
				cells += cell.Text() + "\t"
			})
			fmt.Println("table row: ", cells)
		})
	case "meta":
		if name, exists := s.Attr("name"); exists && name == "author" {
			content, _ := s.Attr("content")
			fmt.Println("Author: ", content)
		}
	}
}
