package rewrites

import (
	//	"github.com/tdewolff/parse/v2/css"
	"golang.org/x/net/html"
	//	"encoding/xml"
	"io/ioutil"
	"io"
	"fmt"
	//	"os"
	"regexp"
	"strings"
)

// This would have to be modified in the future when path support is added
func ProxyUri(proxyUri string) string {
	// TODO: Support paths
	re := regexp.MustCompile(`(\:\/)([^\/])`)
	proxyUri = re.ReplaceAllString(proxyUri, "$1/$2")

	return proxyUri
}

// TODO: Write a proper header parser
func Header(key string, val []string) []string {
	// TODO: Continue adding more header rewrites
	valStr := strings.Join(val, "; ")
	switch key {
	case "Location":
		// TODO: Change the global config of the status code once global config is added
	case "Set-Cookie":
		// TODO: Fix broken regex
		re1 := regexp.MustCompile(`Domain=(.*?);`)
		// TODO: Figure out how to put domain variable in the middle of the string
		valStr = re1.ReplaceAllString(valStr, "Domain=(insert hostname);")
		re2 := regexp.MustCompile(`Path=(.*?);`)
		// TODO: Figure out how to put domain variable in the middle of the string
		valStr = re2.ReplaceAllString(valStr, "Path=(insert proxy path);")
	}

	val = strings.Split(valStr, "; ")

	return val
}

func Html(body io.ReadCloser) io.ReadCloser {
	tokenizer := html.NewTokenizer(body)
	out := ""

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()

		// Order based on docs
		switch (tokenType) {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
		case html.TextToken:
			out += token.Data
		case html.StartTagToken:
			attr := ""
			for _, elm := range token.Attr {
				if elm.Key == "href" || elm.Key == "src" || elm.Key == "poster" || elm.Key == "data" || elm.Key == "action" || elm.Key == "srcset" || elm.Key == "data-src" || elm.Key == "data-href" {
					// TODO: Support custom paths
					if strings.HasPrefix(elm.Val, "/") {
						elm.Val = "https://e93c3221b7e9.ngrok.io" + "/" + "aHR0cHM6Ly9nb29nbGUuY29tLw==" + elm.Val
					} else {
						elm.Val = "https://e93c3221b7e9.ngrok.io" + "/" + elm.Val
					}
				}
				attr += " " + elm.Key + "=" + "\"" + elm.Val + "\""
			}
			out += "<" + token.Data + attr + ">"
		case html.EndTagToken:
			out += "</" + token.Data + ">"
		case html.SelfClosingTagToken:
			out += "<" + token.Data + "/>"
		case html.CommentToken:
			//	out += token.Data
		case html.DoctypeToken:
			//	out += token.Data
		}
	}
	fmt.Println(out)
	body = ioutil.NopCloser(strings.NewReader(out))
	body.Close()
	return body
}

// TODO: Add css rewrites
// See https://github.com/tdewolff/parse/tree/master/css
/*
func Css(body io.ReadCloser) io.ReadCloser {
	// I'm unsure if this will work with io.ReadCloser
	tokenizer := css.NewLexer(parse.NewInput(body))

	for {
		tokenType, token := tokenizer.Next()
		// TODO: Check eof error and break if so
		switch tokenType {
		case css.ErrorToken:
			// TODO: Do error logging
		case css.AtKeywordToken:
		case css.URLToken:
		}
	}

	// TODO: Return io.ReadCloser body
	return body
}
*/

// TODO: Add xml rewrites for external entities (low priority)
// Use https://golang.org/pkg/encoding/xml/

// TODO: Add js injection
/*
func Js(body io.ReadCloser) io.ReadCloser {
	// Needs to read bytes instead
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	file, err := ioutil.ReadFile("././static/inject.js")

	// Don't know if this formatting will work
	bodyBytes := append(file, buf)
	// TODO: Convert bodyBytes to io.ReadCloser

	return body
}
*/
