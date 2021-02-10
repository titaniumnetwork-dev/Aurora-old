package rewrites

import (
	"golang.org/x/net/html"
	//	"bytes"
	//	"io/ioutil"
	"io"
	"log"
	"regexp"
	"strings"
)

// This would have to be modified in the future when path support is added
func ProxyUrl(url string) string {
	re := regexp.MustCompile(`(\:\/)([^\/])`)
	url = re.ReplaceAllString(url, "$1/$2")

	return url
}

func Header(key string, val []string) []string {
	// TODO: Continue adding more header rewrites
	valStr := strings.Join(val, "; ")
	switch key {
	case "Location":
		// TODO: Change the global config of the status code once global config is added
	case "Set-Cookie":
		re1 := regexp.MustCompile(`Domain=(.*?);`)
		// TODO: Insert data once configuration is supported
		valStr = re1.ReplaceAllString(valStr, "Domain=(insert hostname);")
		re2 := regexp.MustCompile(`Path=(.*?);`)
		// TODO: Insert data once configuration is supported
		valStr = re2.ReplaceAllString(valStr, "Path=(insert proxy path);")
	}

	val = strings.Split(valStr, "; ")

	return val
}

// TODO: Add html parser rewrites (almost done)
/*
func Html(body io.ReadCloser) io.ReadCloser {
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tokenType {
		case html.ErrorToken:
			log.Fatal(err)
		case html.StartTagToken:
			for _, elm := range token.Attr {
				if elm.Key == "href" || elm.Key == "src" || elm.Key == "poster" || elm.Key == "data" || elm.Key == "action" || elm.Key == "srcset" || elm.Key == "data-src" || elm.Key == "data-href" {
					if strings.HasPrefix(elm.Val, "/") {
						// TODO: Insert data once configuration is supported
						elm.Val = "(insert proxy url)" + elm.Val
					}
				}
			}
		}
	}
	// TODO: Return io.ReadCloser body
	return body
}
*/

// TODO: Add css rewrites
// Use a css parser

// TODO: Add xml rewrites (for external entities)
// Use https://golang.org/pkg/encoding/xml/

// TODO: Add js injection
/*
func Js(body io.ReadCloser) io.ReadCloser {
	// Needs to read bytes instead
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	file, err := ioutil.ReadFile("inject.js")

	// Don't know if this formatting will work
	bodyBytes := append(file, buf)
	// TODO: Convert bodyBytes to io.ReadCloser

	return body
}
*/
