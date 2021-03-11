package config

import (
	"net/url"
)

var (
	BlockedDomains    = [0]string{}
	BlockedHeaders    = [4]string{"Content-Security-Policy", "Content-Security-Policy-Report-Only", "Strict-Transport-Security", "X-Frame-Options"}
	BlockedUserAgents = [0]string{}

	Scheme   string
	URL      *url.URL
	ProxyURL *url.URL

	HTTPPrefix       string
	HTTPPrefixExists bool
	WSPrefix         string
	WSPrefixExists   bool

	/*
	   var WRTCPrefixExists bool
	   var WRTCPrefix String
	*/

	Port          string
	PortExists    bool
	SSLCert       string
	SSLCertExists bool
	SSLKey        string
	SSLKeyExists  bool
)
