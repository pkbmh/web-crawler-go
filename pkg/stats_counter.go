package pkg

import (
	"sync/atomic"
	"web-crawler/internal"
)

type StatsCounter interface {
	RecordFetchSuccess(count int64);
	RecordFetchFailure(count int64);
	RecordParseSuccess(count int64);
	RecordParseFailure(count int64);
	RecordUncrawlableLink(count int64);
	GetSnapshot() internal.CrawlerStats;
}

type statsCounterImpl struct {
	fetchSuccess int64
	fetchFailed int64
	parseSucess int64
	parseFailed int64
	uncrawlableLink int64
}

func NewStatsCounter() StatsCounter{
	return &statsCounterImpl{}
}

func (c *statsCounterImpl) RecordFetchSuccess(count int64) {
	atomic.AddInt64(&c.fetchSuccess, count)
}

func (c *statsCounterImpl) RecordFetchFailure(count int64) {
	atomic.AddInt64(&c.fetchFailed, count)
}

func (c *statsCounterImpl) RecordParseSuccess(count int64) {
	atomic.AddInt64(&c.parseSucess, count)
}

func (c *statsCounterImpl) RecordParseFailure(count int64) {
	atomic.AddInt64(&c.parseFailed, count)
}

func (c *statsCounterImpl) RecordUncrawlableLink(count int64) {
	atomic.AddInt64(&c.uncrawlableLink, count)
}

func (c *statsCounterImpl) GetSnapshot() internal.CrawlerStats {
	return internal.NewCrawlerStats(c.fetchSuccess, c.fetchFailed, c.parseSucess, c.parseFailed, c.uncrawlableLink)
}
