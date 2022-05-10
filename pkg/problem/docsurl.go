package problem

import "net/url"

// DocsHost is the host name of the documentation page. Use in various helper
// functions as a fallback if no host is provided.
var DocsHost = "wharf.iver.com"

// ConvertURLToAbsDocsURL adds schema and sets the host if that has not been set.
func ConvertURLToAbsDocsURL(u url.URL) *url.URL {
	if !u.IsAbs() {
		u.Scheme = "https"
		u.Host = DocsHost
	}
	if u.Fragment == "" && u.Host == DocsHost {
		u.Fragment = u.Path
		u.Path = "/"
	}
	return &u
}
