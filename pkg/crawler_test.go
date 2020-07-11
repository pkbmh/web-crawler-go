package pkg

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/magiconair/properties/assert"
	"net/url"
	"sort"
	"testing"
	"web-crawler/config"
	"web-crawler/internal"
)

type mockObject struct{
	visitedUrls []string
}

func (c *mockObject) fetch(urlNode *internal.UrlNode) (*FetchResponse, error) {
	return &FetchResponse{urlNode: urlNode}, nil
}

func (c *mockObject) 	parse(request *FetchResponse) (result []*url.URL, doc *goquery.Document, err error) {
	var urlStrings []string
	if (request.urlNode.Depth == 0) {
		urlStrings = []string{"http://localhost:8080/aboutme"}
	} else if (request.urlNode.Depth == 1){
		urlStrings = []string{"http://localhost:8080/aboutme", "http://localhost:8080/contact"}
	} else {
		urlStrings = []string{"http://localhost:8080/aboutme", "http://localhost:8080/contact", "http://localhost:8080/xyz"}
	}

	urlObjects := make([]*url.URL, 0, len(urlStrings))

	for _, urlString := range urlStrings {
		urlObject, _ := internal.NormalizedStringUrl(urlString)
		urlObjects = append(urlObjects, urlObject)
	}
	return urlObjects, nil, nil
}

func (c *mockObject) visit(response *ParseResponse) {
	//log.Info(response.urlNode.actualUrl.String())
	c.visitedUrls = append(c.visitedUrls, response.urlNode.ActualUrl.String())
}

func Test_Crawler_Crawl(t *testing.T) {
	seedUrl := "http://localhost:8080/"
	crawlerConfig := &config.CrawlerConfig {
		Parallelism:           config.DefaultParallelism,
		MaxDepth:              4,
		IgnoreRobotsTxt:       config.DefaultIgnoreRobotsTxt,
		CrawlRate:             config.DefaultCrawlRate,
		UserAgent:             config.DefaultUserAgent,
		SeedUrl:               seedUrl,
		RequestTimeOutInMillS: config.DefaultRequestTimeOutInMillS,
		QueueSize:             config.DefaultQueueSize,
	}
	visitedUrls := make([]string, 0)
	m := &mockObject{
		visitedUrls:visitedUrls,
	}
	statsCounter := NewStatsCounter()

	storage := NewInMemoryStore()
	validator,err := NewDefaultValidator(storage, crawlerConfig)
	if (err != nil) {
		panic(err)
	}
	crawler, err := NewDefaultCrawler(crawlerConfig, m, m, statsCounter, m, validator, storage)
	if err == nil {
		_ = crawler.Crawl()
	}
	wantedCrawledUrl := []string{"http://localhost:8080", "http://localhost:8080/aboutme", "http://localhost:8080/contact", "http://localhost:8080/xyz"}
	gotCrawledUrl := m.visitedUrls
	sort.Strings(gotCrawledUrl)
	assert.Equal(t, gotCrawledUrl, wantedCrawledUrl)
}

func Test_Crawler_Crawl_MaxDepth(t *testing.T) {
	seedUrl := "http://localhost:8080/"
	crawlerConfig := &config.CrawlerConfig {
		Parallelism:           config.DefaultParallelism,
		MaxDepth:              3,
		IgnoreRobotsTxt:       config.DefaultIgnoreRobotsTxt,
		CrawlRate:             config.DefaultCrawlRate,
		UserAgent:             config.DefaultUserAgent,
		SeedUrl:               seedUrl,
		RequestTimeOutInMillS: config.DefaultRequestTimeOutInMillS,
		QueueSize:             config.DefaultQueueSize,
	}

	visitedUrls := make([]string, 0)
	m := &mockObject{
		visitedUrls:visitedUrls,
	}
	statsCounter := NewStatsCounter()

	storage := NewInMemoryStore()
	validator,err := NewDefaultValidator(storage, crawlerConfig)
	if (err != nil) {
		panic(err)
	}
	crawler, err := NewDefaultCrawler(crawlerConfig, m, m, statsCounter, m, validator, storage)
	if err == nil {
		_ = crawler.Crawl()
	}

	wantedCrawledUrl := []string{"http://localhost:8080", "http://localhost:8080/aboutme", "http://localhost:8080/contact"}
	gotCrawledUrl := m.visitedUrls
	sort.Strings(gotCrawledUrl)
	assert.Equal(t, gotCrawledUrl, wantedCrawledUrl)
}

func Test_Crawler_GetMetrics(t *testing.T) {
	seedUrl := "http://localhost:8080/"
	crawlerConfig := &config.CrawlerConfig {
		Parallelism:           config.DefaultParallelism,
		MaxDepth:              4,
		IgnoreRobotsTxt:       config.DefaultIgnoreRobotsTxt,
		CrawlRate:             config.DefaultCrawlRate,
		UserAgent:             config.DefaultUserAgent,
		SeedUrl:               seedUrl,
		RequestTimeOutInMillS: config.DefaultRequestTimeOutInMillS,
		QueueSize:             config.DefaultQueueSize,
	}
	visitedUrls := make([]string, 0)
	m := &mockObject{
		visitedUrls:visitedUrls,
	}
	statsCounter := NewStatsCounter()

	storage := NewInMemoryStore()
	validator,err := NewDefaultValidator(storage, crawlerConfig)
	if (err != nil) {
		panic(err)
	}
	crawler, err := NewDefaultCrawler(crawlerConfig, m, m, statsCounter, m, validator, storage)
	if err == nil {
		_ = crawler.Crawl()
	}
	gotSnapShot := statsCounter.GetSnapshot()
	wantSnapShot := internal.NewCrawlerStats(4, 0, 4, 0, 0)
	assert.Equal(t, gotSnapShot, wantSnapShot)
}
