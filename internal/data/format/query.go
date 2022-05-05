package format

import "net/url"

func MatchURLQueryFormat(url *url.URL, format string) bool {
	return url.Query().Get("format") == format
}
