package rewrites

import (
	//	"github.com/tdewolff/parse/v2/css"
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/global"
	"golang.org/x/net/html"
	//	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	//	"os"
	"encoding/base64"
	"net/url"
	"bytes"
	"regexp"
	"strings"
)

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
		valStr = re1.ReplaceAllString(valStr, "Domain="+global.URI+";")
		re2 := regexp.MustCompile(`Path=(.*?);`)
		valStr = re2.ReplaceAllString(valStr, "Path="+global.Path+";")
	}

	val = strings.Split(valStr, "; ")

	return val
}

func elmAttrRewrite(key string, val string) string {
	if key == "href" || key == "src" || key == "poster" || key == "data" || key == "action" || key == "srcset" || key == "data-src" || key == "data-href" {
		attrURI, err := url.Parse(val)
		if err != nil || attrURI.Scheme == "" || attrURI.Host == "" {
			val = global.Proto + global.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(global.ProxyURI + val))
		} else {
			val = global.Proto + global.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(val))
		}
	}
	attr := " " + key + "=" + "\"" + val + "\""
	return attr
}

// TODO: include womginx in inline and element scripts
func Html(body io.ReadCloser) io.ReadCloser {
	tokenizer := html.NewTokenizer(body)
	out := ""

	for {
		// Maybe they can be combined like in CSS
		tokenType := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tokenType {
		case html.TextToken:
			out += token.Data
		case html.StartTagToken:
			attr := ""
			for _, elm := range token.Attr {
				attr += elmAttrRewrite(elm.Key, elm.Val)
			}
			out += "<" + token.Data + attr + ">"
			// fmt.Println("<" + token.Data + attr + ">")
		case html.EndTagToken:
			out += "</" + token.Data + ">"
		case html.SelfClosingTagToken:
			attr := ""
			for _, elm := range token.Attr {
				attr += elmAttrRewrite(elm.Key, elm.Val)
			}
			out += "<" + token.Data + attr + "/>"
		case html.CommentToken:
			out += "<!--" + token.Data + "-->"
		case html.DoctypeToken:
			out += "<!DOCTYPE " + token.Data + ">"
		}
	}

	// fmt.Println(out)

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

	fmt.Println("Debug: HTML rewrites")
	for {
		tokenType, token := tokenizer.Next()
		// TODO: Check eof error and break if so

		switch tokenType {
		case css.URLToken:
		default:
			out += token.Data
		}
	}

	body = ioutil.NopCloser(strings.NewReader(out))
	body.Close()
	return body
}
*/

// TODO: Add xml rewrites for external entities (low priority)
// Use https://golang.org/pkg/encoding/xml/

// TODO: Add svg rewrites
// Use https://github.com/rustyoz/svg/

// TODO: Add js injection
func Js(body io.ReadCloser) io.ReadCloser {
	// Needs to read bytes instead
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	bodyString := buf.String()

	//	file, err := ioutil.ReadFile("././static/inject.js")

	out := bodyString
	fmt.Println(out)

	body = ioutil.NopCloser(strings.NewReader(out))
	body.Close()
	return body
}
