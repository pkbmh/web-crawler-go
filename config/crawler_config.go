package config

// Default options
const (
	DefaultParallelism		  			int				= 5
	DefaultMaxDepth			  			int				= 3
	DefaultIgnoreRobotsTxt	  			bool			= false
	DefaultCrawlRate		 			int			 	= 100
	DefaultUserAgent          			string          = `Mozilla/5.0 (Windows NT 6.1; rv:15.0) example-webcrawler/1.0 Gecko/20120716 Firefox/15.0a2`
	DefaultRequestTimeOutInMillS		int64	    	= 5000
	DefaultQueueSize		  			int				= 500000
)

type CrawlerConfig struct {
	Parallelism     int
	// Depth Index start from 0. So for MaxDepth=3, it will crawl all links upto depth index 2.
	MaxDepth        int
	// TODO implement honoring robots.txt
	IgnoreRobotsTxt bool
	CrawlRate       int
	UserAgent       string
	SeedUrl         string
	RequestTimeOutInMillS int64
	QueueSize		int
}

func NewDefaultCrawlerConfig(seedUrl string) *CrawlerConfig {
	return &CrawlerConfig {
		DefaultParallelism,
		DefaultMaxDepth,
		DefaultIgnoreRobotsTxt,
		DefaultCrawlRate,
		DefaultUserAgent,
		seedUrl,
		DefaultRequestTimeOutInMillS,
		DefaultQueueSize,
	}
}