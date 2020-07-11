package pkg

import (
	log "github.com/sirupsen/logrus"
)

// Do your stuff with crawler here!
type Visitor interface {
	visit(parseRespose *ParseResponse);
}

type DefaultVisitor struct {}

func (c *DefaultVisitor) visit(parseRespose *ParseResponse) {
	urlNode := parseRespose.urlNode
	title := parseRespose.doc.Find("title").Text()
	log.Infof("Visit URL - %s, Depth - %d, title - %s", urlNode.NormalizedUrl, urlNode.Depth, title)
}