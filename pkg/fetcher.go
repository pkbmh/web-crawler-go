package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"web-crawler/config"
	"web-crawler/internal"
)

type FetchResponse struct {
	resp *http.Response
	urlNode *internal.UrlNode
}

type Fetcher interface {
	fetch(urlNode *internal.UrlNode) (*FetchResponse, error)
}

type defaultFetcher struct {
	crawlerConfig *config.CrawlerConfig
	httpClient    *http.Client
}

func NewDefaultFetcher(crawlerConfig *config.CrawlerConfig) Fetcher {
	requestTimeOut := time.Duration(crawlerConfig.RequestTimeOutInMillS)
	httpClient := &http.Client{
		Timeout:requestTimeOut * time.Millisecond,
	}

	fetcher := &defaultFetcher{
		crawlerConfig:crawlerConfig,
		httpClient:httpClient,
	}
	return fetcher
}

func (c *defaultFetcher) fetch(urlNode *internal.UrlNode) (*FetchResponse, error) {
	req, _ := http.NewRequest("GET", urlNode.NormalizedUrl.String(), nil)
	req.Header.Set("User-Agent", c.crawlerConfig.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Errorf("Error req %v %s", req, err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	parsingRequest := &FetchResponse{
		resp:resp,
		urlNode:urlNode,
	}

	return parsingRequest, nil
}