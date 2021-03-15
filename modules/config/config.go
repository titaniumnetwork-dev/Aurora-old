package config

import "net/url"

type yaml struct {
	BlockedDomains    []string
	BlockedUserAgents []string

	// TODO: Support subdomains too
	HTTPPrefix string
	WSPrefix   string

	/*
	   var WRTCPrefixExists bool
	   var WRTCPrefix String
	*/

	Port string

	SSLOverProxy bool
	Cert         string
	Key          string

	Cap int64
}

var (
	Scheme   string
	URL      *url.URL
	ProxyURL *url.URL

	YAML = yaml{}
)
