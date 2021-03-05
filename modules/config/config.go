package global

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
)

// TODO: Switch from environment variable to json

val, err := ioutil.ReadFile("../../config.json")
if err != nil {
	log.Fatal(err)
}

err := json.Unmarshal(val) 
if err != nil {
	log.Fatal(err)
}

// Old Code - Ignore for now

var Cookie string
var CookieExists bool

var URL url.URL
var ProxyURL url.URL

var Prefix string
var PrefixExists bool
var Port string
var PortExists bool
var SSLCert string
var SSLCertExists bool
var SSLKey string
var SSLKeyExists bool
