**Web Crawler**

A simple web crawler to crawl single url. Can be easily extended to multiple starting Urls.

Currently it supports - 
* Max Depth
* Parallel Crawling
* Max Crawl Rate

To run:

dep ensure -v
go run init/main.go -seedurl https://www.bloomberg.com/

TODO:

* Handle Robots.txt
* Check Get HEAD before fetching url.# web-crawler-go
