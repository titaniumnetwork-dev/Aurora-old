package config

import (
	"net/url"
)

var Cookie string
var CookieExists bool

var URL url.URL
var ProxyURL url.URL

var HTTPPrefix string
var HTTPPrefixExists bool
var WSPrefix string
var WSPrefixExists bool
var Port string
var PortExists bool
var SSLCert string
var SSLCertExists bool
var SSLKey string
var SSLKeyExists bool
