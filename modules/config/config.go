package config

import (
	"net/url"
)

var (
	BlockedDomains    = [0]string{}
	BlockedHeaders    = [4]string{"Content-Security-Policy", "Content-Security-Policy-Report-Only", "Strict-Transport-Security", "X-Frame-Options"}
	BlockedUserAgents = [0]string{}

	HTTPScheme       string
	HTTPURL          *url.URL
	HTTPProxyURL     *url.URL
	HTTPPrefix       string
	HTTPPrefixExists bool
	WSScheme         string
	WSURL            *url.URL
	WSProxyUrl       *url.URL
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
