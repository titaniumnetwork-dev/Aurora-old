package config

import "net/url"

type yaml struct {
	BlockedDomains    []string `yaml:"BlockedDomains"`
	BlockedUserAgents []string `yaml:"BlockedUserAgents"`

	// TODO: Support subdomains too
	HTTPPrefix string `yaml:"HTTPPrefix"`
	WSPrefix   string `yaml:"WSPrefix"`

	/*
	   var WRTCPrefixExists bool
	   var WRTCPrefix String
	*/

	Port string `yaml:"Port"`

	SSLOverProxy bool   `yaml:"SSLOverProxy"`
	Cert         string `yaml:"Cert"`
	Key          string `yaml:"Key"`

	Cap int64 `yaml:"Cap"`
}

var (
	Scheme   string
	URL      *url.URL
	ProxyURL *url.URL

	YAML = yaml{}
)
