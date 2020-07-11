package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"
	"net/url"
	"sync"
	"web-crawler/config"
	"web-crawler/internal"
)

type Crawler interface {
	Crawl() error
	GetMetrics() internal.CrawlerStats
}

type defaultCrawler struct {
	CrawlerConfig *config.CrawlerConfig
	wg            *sync.WaitGroup
	fetchQueue    chan *internal.UrlNode
	parsingQueue  chan *FetchResponse
	fetcher       Fetcher
	parser        Parser
	seedUrlNode   *internal.UrlNode
	statsCounter  StatsCounter
	rateLimiter   ratelimit.Limiter
	visitor       Visitor
	validator     Validator
	storage       Storage
}

func NewDefaultCrawler(crawlerConfig *config.CrawlerConfig, fetcher Fetcher, parser Parser, statsCounter StatsCounter, visitor Visitor, validator Validator, storage Storage) (Crawler, error) {
	if crawlerConfig == nil || crawlerConfig.SeedUrl == "" {
		return nil, fmt.Errorf("SeedUrl Can't be Nil")
	}
	rateLimiter := ratelimit.New(crawlerConfig.CrawlRate)
	fetchQueue := make(chan *internal.UrlNode, crawlerConfig.QueueSize)
	parsingQueue := make(chan *FetchResponse)
	wg := &sync.WaitGroup{}

	c := defaultCrawler{
		CrawlerConfig:crawlerConfig,
		wg:wg,
		fetcher:fetcher,
		parser:parser,
		fetchQueue:fetchQueue,
		parsingQueue:parsingQueue,
		rateLimiter:rateLimiter,
		statsCounter:statsCounter,
		visitor:visitor,
		validator:validator,
		storage: storage,
	}

	seedUrl := crawlerConfig.SeedUrl
	var err error
	c.seedUrlNode, err = internal.UrlStringToUrlNode(seedUrl,nil)
	if err != nil {
		log.Errorf("Error converting Url String to Url Node %s, %s",seedUrl, err)
		return nil, err
	}
	return &c, nil
}

func (c *defaultCrawler) Crawl() error {
	c.runFetcher()
	c.runParser()
	err := c.addToFetchQueue(c.seedUrlNode)
	if err != nil {
		return err
	}
	c.wg.Wait()
	return nil
}

func (c *defaultCrawler) fetchWorker() {
	for s := range c.fetchQueue {
		_ = c.rateLimiter.Take()
		parseRequest, err := c.fetcher.fetch(s)
		if err != nil {
			c.wg.Done()
			log.Errorf("Fetch Error %s for URL %s ->  " , err, s.NormalizedUrl)
			c.statsCounter.RecordFetchFailure(1)
		} else {
			c.statsCounter.RecordFetchSuccess(1)
			c.parsingQueue <- parseRequest
		}
	}
}

func (c *defaultCrawler) parseWorker() {
	for parseRequest := range c.parsingQueue {
		links, doc, err := c.parser.parse(parseRequest)
		if err == nil {
			c.visitor.visit(&ParseResponse{urlNode: parseRequest.urlNode, doc:doc})
			c.collectUrls(links, parseRequest.urlNode)
			c.statsCounter.RecordParseSuccess(1)
			err1 := c.storage.Store(parseRequest.urlNode.NormalizedUrl.String(), internal.VISITED)
			if err1 != nil {
				log.Errorf("Error storing the url %s value VISITED, error - %s", parseRequest.urlNode.NormalizedUrl.String(), err1)
			}
		} else {
			log.Errorf("Error Parsing URL %s error - %s ", parseRequest.urlNode.NormalizedUrl, err)
			c.statsCounter.RecordParseFailure(1)
			err = c.storage.Store(parseRequest.urlNode.NormalizedUrl.String(), internal.PARSE_ERROR)
			if (err != nil) {
				log.Errorf("Error storing the url %s value PARSE_ERROR, error - %s", parseRequest.urlNode.NormalizedUrl.String(), err)
			}
		}
		c.wg.Done()
	}
}

func (c *defaultCrawler) runFetcher() {
	for i := 0; i < c.CrawlerConfig.Parallelism; i++ {
		go c.fetchWorker()
	}
}

func (c *defaultCrawler) runParser() {
	for i := 0; i < c.CrawlerConfig.Parallelism; i++ {
		go c.parseWorker()
	}
}

func (c *defaultCrawler) collectUrls(links []*url.URL, parentNode *internal.UrlNode) {

	if (links == nil || len(links) == 0) {
		return
	}

	for _, link := range links {
		linkNode,err := internal.UrlToUrlNode(link, parentNode)
		if (err != nil) {
			log.Errorf("Error Converting Url %s to UrlNode, error - %s", link.String(), err)
		}
		if (c.validator.isValid(linkNode)) {
			err = c.addToFetchQueue(linkNode)
			if (err != nil) {
				log.Errorf("Error Adding Url %s to FetchQueue, error - %s", link.String(), err)
			}
		}
	}
}

func (c *defaultCrawler) addToFetchQueue(urlNode *internal.UrlNode) error {
	c.wg.Add(1)
	err := c.storage.Store(urlNode.NormalizedUrl.String(), internal.QUEUED)
	if err != nil {
		return err
	}
	c.fetchQueue <- urlNode
	return nil
}

func (c *defaultCrawler) GetMetrics() internal.CrawlerStats {
	return c.statsCounter.GetSnapshot()
}