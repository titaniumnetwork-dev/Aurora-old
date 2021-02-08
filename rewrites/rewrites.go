package rewrites

import (
	//	golang.org/x/net/html"
	//	"bytes"
	//	"io/ioutil"
	//	"io"
	"regexp"
	//	"strings"
)

// This would have to be modified in the future when path support is added
func ProxyUrl(url string) string {
	re := regexp.MustCompile(`(\:\/)([^\/])`)
	url = re.ReplaceAllString(url, "$1/$2")

	return url
}

/*
func Header(key string, val []string) []string {
	// TODO: Continue adding more header rewrites
	switch key {
	case "Location":
		// TODO: Change the global config of the status code once global config is added
	case "Set-Cookie":
		re1 := regexp.MustCompile(`Domain=(.*?);`)
		// TODO: Insert data once configuration is supported
		val := re1.ReplaceAllString(strings.Join(val, "; "), "Domain=(insert hostname);")
		re2 := regexp.MustCompile(`Path=(.*?);`)
		// TODO: Insert data once configuration is supported
		val = re2.ReplaceAllString(strings.Join(val, "; "), "Path=(insert proxy path);")
	}

	// I don't know if this would work
	val = strings.Split(valString, "; ")

	return val
}
*/

// TODO: Add html parser rewrites
/*
func HTML(body io.ReadCloser) io.ReadCloser {
	// TODO: Actually save the changes
	tokenizer := html.NewTokenizer(body)
	// TODO: Switch to using a while loop that ends when the file is done being parsed
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			log.Fatal(tokenizer.Error())
		case html.StartTagToken:
			token := tokenizer.Token()
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
}
*/

// TODO: Add css rewrites
// Use a css parser

// TODO: Add xml rewrites
// Use https://golang.org/pkg/encoding/xml/

// TODO: Add js rewrites
// Use a js parser

// TODO: Add js injection
/*
func JsInject(body io.ReadCloser) io.ReadCloser {
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
