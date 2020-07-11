package internal

import (
	"github.com/PuerkitoBio/purell"
	"net/url"
)

type UrlNode struct {
	// ActualUrl is the actual url.
	ActualUrl     *url.URL
	NormalizedUrl *url.URL
	ParentNode    *UrlNode
	Depth         int
}

type UrlState int32

const (
	UNVISITED	UrlState = iota
	QUEUED
	VISITED
	PARSE_ERROR
)

func UrlToUrlNode(url *url.URL, parentNode *UrlNode) (*UrlNode, error) {
	normalizedUrl, err := NormalizedUrl(url)

	if err != nil {
		return nil, err
	}
	depth := 0
	if parentNode != nil {
		depth = parentNode.Depth + 1
	}

	return &UrlNode{
		ActualUrl:     url,
		NormalizedUrl: normalizedUrl,
		ParentNode:    parentNode,
		Depth:         depth,
	}, nil
}

func UrlStringToUrlNode(urlString string, parentNode *UrlNode) (*UrlNode, error) {
	urlZ, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	return UrlToUrlNode(urlZ, parentNode)
}

func NormalizedUrl(urlZ *url.URL) (*url.URL, error) {
	if (urlZ == nil) {
		return nil, nil
	}
	normalizedUrlString := purell.NormalizeURL(urlZ, purell.FlagsAllGreedy)
	normalizedUrl, err := urlZ.Parse(normalizedUrlString)
	return normalizedUrl, err
}

func NormalizedStringUrl(urlString string) (*url.URL, error) {
	if (urlString == "") {
		return nil, nil
	}
	normalizedUrlString, err := purell.NormalizeURLString(urlString, purell.FlagsAllGreedy)
	if err != nil {
		return nil, err
	}
	normalizedUrl, err := url.Parse(normalizedUrlString)
	return normalizedUrl, err
}