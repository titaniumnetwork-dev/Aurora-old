package global

import (
	"io/ioutil"
	"log"
	"net/url"
)

// New wip json based config

// TODO: Bind to an empty interface and have getters and setters functions 

val, err := ioutil.ReadFile("./config.json")
if err != nil {
	log.Fatal(err)
}

type Config struct {}

json.Unmarshal(val, &config)

if !Config.Port {
	log.Fatal("You need to specify a port")
} else if !Config.Prefix {
	log.Fatal("You need to specify a prefix")
}

// Old environment variable based config

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
