package pkg

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"web-crawler/config"
	"web-crawler/internal"
)

type Validator interface {
	isValid(node *internal.UrlNode) bool
	filter(nodes []*internal.UrlNode) []*internal.UrlNode
}

type defaultValidator struct {
	storage           Storage
	normalizedSeedUrl *url.URL
	crawlerConfig     *config.CrawlerConfig
}

func NewDefaultValidator(storage Storage, crawlerConfig *config.CrawlerConfig) (Validator, error) {
	seedUrl := crawlerConfig.SeedUrl
	normalizedSeedUrl, err := internal.NormalizedStringUrl(seedUrl)
	if err != nil {
		log.Errorf("Error converting Url String to Normalized Url %s, %s",seedUrl, err)
		return nil, err
	}

	return &defaultValidator{
		storage:storage,
		normalizedSeedUrl:normalizedSeedUrl,
		crawlerConfig: crawlerConfig,
	}, nil
}

// TODO add robots.txt validation
func (c *defaultValidator) isValid(urlNode *internal.UrlNode) bool {
	if (c.crawlerConfig.MaxDepth != 0 && urlNode.Depth >= c.crawlerConfig.MaxDepth) {
		return false
	}
	if (c.normalizedSeedUrl.Host != urlNode.NormalizedUrl.Host) {
		return false
	}
	isPresent := c.storage.isPresent(urlNode.NormalizedUrl.String())
	return !isPresent
}

func (c *defaultValidator) filter(urlNodes []*internal.UrlNode) []*internal.UrlNode {
	validUrls := make([]*internal.UrlNode, 0, len(urlNodes))

	for _, linkUrlNode := range urlNodes {
		if c.isValid(linkUrlNode) {
			validUrls = append(validUrls, linkUrlNode)
			err := c.storage.Store(linkUrlNode.NormalizedUrl.String(), internal.QUEUED)
			if err != nil {
				log.Errorf("Error Storing Url %s error - %s", linkUrlNode.NormalizedUrl.String(), err)
			}
		}
	}
	return validUrls
}