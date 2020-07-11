package pkg

import (
	"github.com/magiconair/properties/assert"
	"testing"
	"web-crawler/config"
	"web-crawler/internal"
)

func Test_defaultValidator_isValid(t *testing.T) {
	storage := NewInMemoryStore()
	seedUrl := "http://localhost:8080/"
	crawlerConfig := &config.CrawlerConfig {
		Parallelism:           config.DefaultParallelism,
		MaxDepth:              2,
		IgnoreRobotsTxt:       config.DefaultIgnoreRobotsTxt,
		CrawlRate:             config.DefaultCrawlRate,
		UserAgent:             config.DefaultUserAgent,
		SeedUrl:               seedUrl,
		RequestTimeOutInMillS: config.DefaultRequestTimeOutInMillS,
		QueueSize:             config.DefaultQueueSize,
	}
	validator, err := NewDefaultValidator(storage, crawlerConfig)
	assert.Equal(t, err, nil)

	urlString1 := "http://localhost:8080/aboutme"
	urlString2 := "http://facebook.com/aboutme"

	urlNode1, err := internal.UrlStringToUrlNode(urlString1, nil)
	urlNode2, err := internal.UrlStringToUrlNode(urlString2, nil)

	got1 := validator.isValid(urlNode1)
	assert.Equal(t, got1, true, "Valid Url Check 1")

	got2 := validator.isValid(urlNode2)
	assert.Equal(t, got2, false, "Valid Url Check 2")

	// After already in storage
	_ = storage.Store(urlString1, internal.QUEUED)
	got3 := validator.isValid(urlNode1)
	assert.Equal(t, got3, false, "Valid Url Check 3")

	// Validate Max Depth Check.
	urlString3 := "http://localhost:8080/contact"
	urlNode1, err = internal.UrlStringToUrlNode(urlString1, urlNode2)
	urlNode3, err := internal.UrlStringToUrlNode(urlString3, urlNode1)
	got4 := validator.isValid(urlNode3)
	assert.Equal(t, got4, false, "Valid Url Check 4")
}