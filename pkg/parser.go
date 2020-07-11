package pkg

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
	log "github.com/sirupsen/logrus"
	"net/url"
	"path"
	"strings"
	"web-crawler/internal"
)

type ParseResponse struct {
	urlNode *internal.UrlNode
	doc		*goquery.Document
}

type Parser interface {
	parse(request *FetchResponse) (result []*url.URL, doc *goquery.Document, err error)
}

type DefaultParser struct {}

var (
	aHrefMatcher    = cascadia.MustCompile("a[href]")
	baseHrefMatcher = cascadia.MustCompile("base[href]")
)

// Scrape the document's content to gather all links
// Implementation taken from github.com/PuerkitoBio/gocrawl
func (c *DefaultParser) parse(parseRequest *FetchResponse) (result []*url.URL, doc *goquery.Document, err error) {
	defer func() {
		if err == nil {
			parseRequest.resp.Body.Close()
		}
	}()

	doc, err = goquery.NewDocumentFromReader(parseRequest.resp.Body)
	if (err != nil) {
		return nil, nil, err
	}
	doc.Url = parseRequest.urlNode.ActualUrl
	baseURL, _ := doc.FindMatcher(baseHrefMatcher).Attr("href")
	urls := doc.FindMatcher(aHrefMatcher).Map(func(_ int, s *goquery.Selection) string {
		val, exist := s.Attr("href")
		var err error
		if exist && baseURL != "" {
			val, err = c.handleBaseTag(doc.Url, baseURL, val)
			if err != nil {
				log.Errorf("Error handling base Tage for %s, bsae Url %s, error - %s", val, baseURL, err)
			}
		}
		return val
	})
	for _, s := range urls {
		// If href starts with "#", then it points to this same exact URL, ignore (will fail to parse anyway)
		if len(s) > 0 && !strings.HasPrefix(s, "#") {
			if parsed, e := url.Parse(s); e == nil {
				parsed = doc.Url.ResolveReference(parsed)
				result = append(result, parsed)
			} else {
				log.Errorf("ignore on unparsable policy %s: %s", s, e.Error())
			}
		}
	}
	return
}

func (c *DefaultParser) handleBaseTag(root *url.URL, baseHref string, aHref string) (string, error) {
	resolvedBase, err := root.Parse(baseHref)
	if err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(aHref)
	if err != nil {
		return "", err
	}
	// If a[href] starts with a /, it overrides the base[href]
	if parsedURL.Host == "" && !strings.HasPrefix(aHref, "/") {
		aHref = path.Join(resolvedBase.Path, aHref)
	}

	resolvedURL, err := resolvedBase.Parse(aHref)
	if err != nil {
		return "", err
	}
	return resolvedURL.String(), nil
}