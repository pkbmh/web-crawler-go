package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"time"
	"web-crawler/config"
	"web-crawler/pkg"
)

var (
	seedUrl = flag.String("seedurl", "", "Seed Url")
	parallelism = flag.Int("parallelism", config.DefaultParallelism, "number of parallel workers")
	ignoreRobotsTxt = flag.Bool("ignore_robots_txt", config.DefaultIgnoreRobotsTxt, "Ignore Robots.Txt")
	crawlRate = flag.Int("crawl_rate", config.DefaultCrawlRate, "Crawl QPS")
	userAgent = flag.String("user_agent", config.DefaultUserAgent, "User Agent for request header")
	maxDepth = flag.Int("max_depth", config.DefaultMaxDepth, "Max Depth To Crawl")
	requestTimeOut = flag.Int64("request_timeout_ms", config.DefaultRequestTimeOutInMillS, "Http(s) Request Timeout in Millisecond")
	queueSize = flag.Int("queue_size", config.DefaultQueueSize, "Specify the queue size responsible for holding the fetched urls")
)

func main() {
	flag.Parse()

	if *seedUrl == "" {
		panic("Please Input Seed Url!")
	}

	crawlerConfig := &config.CrawlerConfig{
		Parallelism: *parallelism,
		SeedUrl:     *seedUrl,
		CrawlRate:   *crawlRate,
		IgnoreRobotsTxt: *ignoreRobotsTxt,
		UserAgent: *userAgent,
		MaxDepth: *maxDepth,
		RequestTimeOutInMillS: *requestTimeOut,
		QueueSize:*queueSize,
	}

	storage := pkg.NewInMemoryStore()
	statsCounter := pkg.NewStatsCounter()
	fetcher := pkg.NewDefaultFetcher(crawlerConfig)
	validator,err := pkg.NewDefaultValidator(storage, crawlerConfig)
	if (err != nil) {
		panic(err)
	}
	visitor := &pkg.DefaultVisitor{}
	parser := &pkg.DefaultParser{}

	//TODO use builder pattern to build the crawler.
	crawler, err := pkg.NewDefaultCrawler(crawlerConfig, fetcher, parser, statsCounter, visitor, validator, storage)
	if (err != nil) {
		panic(err)
	}

	// Emit metrics every 2 second.
	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for {
			select {
			case <- ticker.C:
				metrics := statsCounter.GetSnapshot()
				log.Infof("Metrics - %s", metrics.ToString())
			}
		}
	}()

	err = crawler.Crawl()
	if err != nil {
		panic(err)
	}
	ticker.Stop()
	metrics := statsCounter.GetSnapshot()
	log.Infof("Final Metrics - %s", metrics.ToString())
}
