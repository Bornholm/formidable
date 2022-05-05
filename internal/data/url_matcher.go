package data

import "net/url"

type URLMatcher interface {
	Match(url *url.URL) bool
}
