package Config

import (
	"io/ioutil"
	"log"
	"net/url"
)

// Types are ignored when mapping to interface

// See if it is possible to make this is a singleton
func Init() {
	val, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	var config map[string]interface{}

	json.Unmarshal(val, &config)

	// TODO: Adapt this
	if config["port"] {
		log.Fatal("You need to specify a port")
	} else if config["prefix"] {
		log.Fatal("You need to specify a prefix")
	}
}

func Get(val string) string  {
	val := config[val]

	return val, nil
}

func Set(val string) string {
	config[val] = val

	return val, nil
}