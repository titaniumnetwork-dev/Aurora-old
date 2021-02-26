package rewrites

import (
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/global"
	"golang.org/x/net/html"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
	//	"encoding/xml"
	"io"
	"io/ioutil"
	"encoding/base64"
	"net/url"
	"bytes"
	"regexp"
	"log"
	"strings"
)

// TODO: Write a proper header parser
func Header(key string, val []string) []string {
	valStr := strings.Join(val, "; ")

	// TODO: Continue adding more header rewrites
	switch key {
	case "Location":
		// TODO: Change the global config of the status code once global config is added
	case "Set-Cookie":
		// TODO: Fix broken regex
		re1 := regexp.MustCompile(`Domain=(.*?);`)
		valStr = re1.ReplaceAllString(valStr, "Domain=" + global.URI + ";")
		re2 := regexp.MustCompile(`Path=(.*?);`)
		// TODO: This won't work when base64 is fully encoded maybe I can use a cookiejar and split the path in the future
		valStr = re2.ReplaceAllString(valStr, "Path=" + global.Path + ";")
	}

	val = strings.Split(valStr, "; ")

	return val
}

func elmAttr(key string, val string) string {
	if key == "href" || key == "src" || key == "poster" || key == "data" || key == "action" || key == "srcset" || key == "data-src" || key == "data-href" {
		attrURI, err := url.Parse(val)
		if err != nil || attrURI.Scheme == "" || attrURI.Host == "" {
			val = global.Scheme + global.Host + global.Prefix + base64.StdEncoding.EncodeToString([]byte(global.ProxyURI + val))
		} else {
			val = global.Scheme + global.Host + global.Prefix + base64.StdEncoding.EncodeToString([]byte(val))
		}
	}
	attr := " " + key + "=" + "\"" + val + "\""
	return attr
}

func HTML(body io.ReadCloser) io.ReadCloser {
	tokenizer := html.NewTokenizer(body)
	out := ""

	for {
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
				attr += elmAttr(elm.Key, elm.Val)
			}
			out += "<" + token.Data + attr + ">"
		case html.EndTagToken:
			out += "</" + token.Data + ">"
		case html.SelfClosingTagToken:
			attr := ""
			for _, elm := range token.Attr {
				attr += elmAttr(elm.Key, elm.Val)
			}
			out += "<" + token.Data + attr + "/>"
		case html.CommentToken:
			out += "<!--" + token.Data + "-->"
		case html.DoctypeToken:
			out += "<!DOCTYPE " + token.Data + ">"
		}
	}

	body = ioutil.NopCloser(strings.NewReader(out))
	body.Close()
	return body
}

func CSS(body io.ReadCloser) io.ReadCloser {
	tokenizer := css.NewLexer(parse.NewInput(body))
	out := ""

	for {
		tokenType, tokenBytes := tokenizer.Next()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tokenType {
		case css.URLToken:
			data := strings.Replace(string(tokenBytes), "url(", "", 4)
			data = strings.Replace(string(data), ")", "", 1)
		
			uri, err := url.Parse(data)
			if err != nil || uri.Scheme == "" || uri.Host == "" {
				data = global.Scheme + global.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(global.ProxyURI + data))
			} else {
				data = global.Scheme + global.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(data))
			}
	
			out += "url(" + data + ")"
		default:
			out += string(tokenBytes)
		}
	}

	body = ioutil.NopCloser(strings.NewReader(out))
	body.Close()
	return body
}

// TODO: Add xml rewrites for external entities (low priority)
// Use https://golang.org/pkg/encoding/xml/

// TODO: Add svg rewrites
// Use https://github.com/rustyoz/svg/

// TODO: Add js injection
func JS(body io.ReadCloser) io.ReadCloser {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	bodyString := buf.String()

	// TODO: Have newlines
	fileBytes, err := ioutil.ReadFile("././static/inject.js")
	if err != nil {
		log.Println(err)
	}
	file := string(fileBytes)

	out := file + bodyString

	body = ioutil.NopCloser(strings.NewReader(out))
	body.Close()
	return body
}
