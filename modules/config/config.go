package config

import "net/url"

type y struct {
	BlockedDomains    []string `yaml:"BlockedDomains"`
	BlockedHeaders    []string `yaml:"BlockedHeaders"`
	BlockedUserAgents []string `yaml:"BlockedUserAgents"`

	// TODO: Support subdomains too
	HTTPPrefix string `yaml:"HTTPPrefix"`
	WSPrefix   string `yaml:"WSPrefix"`

	/*
	   var WRTCPrefixExists bool
	   var WRTCPrefix String
	*/

	Port string `yaml:"Port"`

	Cert string `yaml:"Cert"`
	Key  string `yaml:"Key"`

	Cap int64 `yaml:"Cap"`
}

var (
	Scheme   string
	URL      *url.URL
	ProxyURL *url.URL

	YAML = y{}
)
