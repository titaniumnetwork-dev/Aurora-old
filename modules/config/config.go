package config

import (
	"net/url"
)

// Maybe remove the exists and instead see if they are nil?
// For this use os.GetEnv

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

	SSLOverProxy       bool
	SSLOverProxyExists bool
	SSLCertExists      bool
	SSLKey             string
	SSLKeyExists       bool
)
