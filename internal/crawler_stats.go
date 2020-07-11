package internal

import "fmt"

type CrawlerStats interface {
	GetTotalCrawledCount() int64
	GetFailedCrawledCount() int64
	GetSuccessCrawledCount() int64
	GetFailedParseCount() int64
	GetSucessParseCount() int64
	GetFailedFetchCount() int64
	GetSucessFetchCount() int64
	GetUnCrawlableLinkCount() int64
	ToString() string
}

type crawlerStatsImpl struct {
	sucessFetchCount int64
	failedFetchCount int64
	sucessParseCount int64
	failedParseCount int64
	uncrawlableLink int64
}

// TODO change it to return Depth wise metrics. At each Depth what is the stats.

func NewCrawlerStats(sucessFetchCount int64, failedFetchCount int64, sucessParseCount int64, failedParseCount int64, uncrawlableLink int64) CrawlerStats {
	crawlerStats := &crawlerStatsImpl{sucessFetchCount: sucessFetchCount,
		failedFetchCount: failedFetchCount,
		sucessParseCount: sucessParseCount,
		failedParseCount: failedParseCount,
		uncrawlableLink:  uncrawlableLink,
	}
	return crawlerStats
}

func (c *crawlerStatsImpl) GetUnCrawlableLinkCount() int64 {
	return c.uncrawlableLink
}

func (c *crawlerStatsImpl) GetTotalCrawledCount() int64 {
	return c.GetSuccessCrawledCount() + c.GetFailedCrawledCount()
}

func (c *crawlerStatsImpl) GetFailedCrawledCount() int64 {
	return c.failedFetchCount + c.failedParseCount
}

func (c *crawlerStatsImpl) GetSuccessCrawledCount() int64 {
	return c.sucessParseCount
}

func (c *crawlerStatsImpl) GetFailedParseCount() int64 {
	return c.failedParseCount
}

func (c *crawlerStatsImpl) GetSucessParseCount() int64 {
	return c.sucessParseCount
}

func (c *crawlerStatsImpl) GetFailedFetchCount() int64 {
	return c.failedFetchCount
}

func (c *crawlerStatsImpl) GetSucessFetchCount() int64 {
	return c.sucessFetchCount
}

func (c *crawlerStatsImpl) ToString() string {
	return fmt.Sprintf("Total - %d, Success - %d, Failed - %d, Uncrawlable - %d",
		c.GetTotalCrawledCount(), c.GetSuccessCrawledCount(), c.GetFailedCrawledCount(), c.GetUnCrawlableLinkCount())
}