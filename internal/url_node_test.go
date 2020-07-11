package internal

import (
	"github.com/magiconair/properties/assert"
	"net/url"
	"testing"
)

func TestUrlToUrlNode(t *testing.T) {
	urlString := "http://localhost:8080/aboutme/../contact"
	urlObject, err := url.Parse(urlString)
	assert.Equal(t, err, nil)
	urlNormalized, err := url.Parse("http://localhost:8080/contact")
	assert.Equal(t, err, nil)

	gotUrlNode, err := UrlToUrlNode(urlObject, nil)
	assert.Equal(t, err, nil)

	wantUrlNode := UrlNode{
		ActualUrl:     urlObject,
		NormalizedUrl: urlNormalized,
		Depth:         0,
		ParentNode:    nil,
	}

	assert.Equal(t, gotUrlNode.NormalizedUrl.String(), wantUrlNode.NormalizedUrl.String())
	assert.Equal(t, gotUrlNode.NormalizedUrl, wantUrlNode.NormalizedUrl)
	assert.Equal(t, gotUrlNode.ActualUrl, wantUrlNode.ActualUrl)
	assert.Equal(t, gotUrlNode.ParentNode, wantUrlNode.ParentNode)
	assert.Equal(t, gotUrlNode.Depth, 0)

	urlString2 := "http://localhost:8080/aboutme/"
	urlObject2, err := url.Parse(urlString2)
	urlNormalized2, err := url.Parse("http://localhost:8080/aboutme")
	assert.Equal(t, err, nil)

	wantUrlNode2 := UrlNode{
		ActualUrl:     urlObject2,
		NormalizedUrl: urlNormalized2,
		Depth:         1,
		ParentNode:    gotUrlNode,
	}

	assert.Equal(t, err, nil)
	gotUrlNode2, err := UrlToUrlNode(urlObject2, gotUrlNode)
	assert.Equal(t, err, nil)
	assert.Equal(t, gotUrlNode2.NormalizedUrl.String(), wantUrlNode2.NormalizedUrl.String())
	assert.Equal(t, gotUrlNode2.NormalizedUrl, wantUrlNode2.NormalizedUrl)
	assert.Equal(t, gotUrlNode2.ActualUrl, wantUrlNode2.ActualUrl)
	assert.Equal(t, gotUrlNode2.ParentNode, wantUrlNode2.ParentNode)
	assert.Equal(t, gotUrlNode2.Depth, 1)
}
