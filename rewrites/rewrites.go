package rewrites

import (
	// "golang.org/x/net/html"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func ProxyUrl(url string) string {
	re := regexp.MustCompile(`(\:\/)([^\/])`)
	url = re.ReplaceAllString(url, "$1/$2")

	return url
}

func Header(key string, val []string) []string {
	// TODO: Continue adding more header rewrites
	switch {
	case key == "Location":
		// TODO: Change the global config of the status code once global config is added
	case key == "Set-Cookie":
		re1 := regexp.MustCompile(`Domain=(.*?);`)
		// TODO: Insert data once global config is added
		val = re1.ReplaceAllString(strings.Join(val, "; "), "Domain=(insert hostname);") 
		re2 := regexp.MustCompile(`Path=(.*?);`)
		// TODO: Insert data once global config is added
		val = re2.ReplaceAllString(strings.Join(val, "; "), "Path=(insert proxy path);")
	}

	val = strings.Split(val, "; ")

	return val
}

// TODO: Add html parser rewrites
func HTML(body io.ReadCloser) io.ReadCloser {
	// TODO: Figure out how to actually save the changes
	tokenizer := html.NewTokenizer(body)
	for {
		if tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			log.Fatal(tokenizer.Error())
		case html.StartTagToken:
			token := tokenizer.Token()
			for _, element := range token.Attr {
				if element.Key == "href" || element.Key == "src" || element.Key == "poster" || element.Key == "data" || element.Key == "action" || element.Key == "srcset" ||  element.Key == "data-src" || element.Key == "data-href"  {
					if strings.HasPrefix(element.Val, "/") {
						// TODO: Insert data once configuration is supported
						eleent.Val = "(insert proxy url)" +  element.Val
					}
				}
			}
		}
	}

}

// TODO: Add js injection
func Js(body io.ReadCloser) io.ReadCloser {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	file, err := ioutil.ReadFile("inject.js")

	bodyBytes := append(file, buf)
	// Convert bodyBytes to io.ReadCloser

	return body
}
