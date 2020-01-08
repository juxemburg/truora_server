package htmlinfo

import (
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)

/*PageInfo represents selected info from an HTML page*/
type PageInfo struct {
	Title   string
	IconURL string
}

func invalidHTMLPageInfo() *PageInfo {
	return &PageInfo{Title: "Invalid HTML page", IconURL: "Invalid HTML page"}
}

/*GetHTMLPageInfo gets the html title, and icon, of the page URL provided*/
func GetHTMLPageInfo(domain string) *PageInfo {
	pageURL := getPageURL(domain)
	resp, err := http.Get(pageURL)
	if err != nil {
		return invalidHTMLPageInfo()
	}

	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return invalidHTMLPageInfo()
	}

	headNode := extractHTMLNode(doc, func(node *html.Node) bool {
		return node.Type == html.ElementNode && node.Data == "head"
	})
	titleNode := extractHTMLNode(headNode, func(node *html.Node) bool {
		return node.Type == html.ElementNode && node.Data == "title"
	})
	iconNode := extractHTMLNode(headNode, func(node *html.Node) bool {
		if node.Type != html.ElementNode || node.Data != "link" {
			return false
		}
		for _, attr := range node.Attr {
			if attr.Key == "type" && attr.Val == "image/x-icon" {
				return true
			}
			if attr.Key == "rel" && (attr.Val == "shortcut icon" || attr.Val == "apple-touch-icon") {
				return true
			}
		}
		return false
	})
	var pageTitle = "Page title not found"
	if titleNode != nil {
		pageTitle = titleNode.FirstChild.Data
	}
	var iconURL = fmt.Sprintf("http://www.google.com/s2/favicons?domain=%v", domain)
	if iconNode != nil {
		for _, attr := range iconNode.Attr {
			if attr.Key == "href" {
				iconURL = attr.Val
				break
			}
		}
	}

	return &PageInfo{Title: pageTitle, IconURL: iconURL}
}

func getPageURL(url string) string {
	uriRegexp := regexp.MustCompile(`^https:\/\/`)
	if !uriRegexp.MatchString(url) {
		return fmt.Sprintf(`https://%v`, url)
	}
	return url
}

func extractHTMLNode(root *html.Node, extractionFn func(node *html.Node) bool) *html.Node {
	if extractionFn(root) {
		return root
	}

	for child := root.FirstChild; child != nil; child = child.NextSibling {
		result := extractHTMLNode(child, extractionFn)
		if result != nil {
			return result
		}
	}

	return nil
}
