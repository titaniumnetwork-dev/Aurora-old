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
	"regexp"
	"strings"
)

// TODO: Write a proper header parser
func Header(key string, val []string) []string {
	valStr := strings.Join(val, "; ")

	// TODO: Continue adding more header rewrites
	switch key {
	case "Location":
	case "Set-Cookie":
		re1 := regexp.MustCompile(`domain=(.*?);`)
		valStr = re1.ReplaceAllString(valStr, "domain=" + global.URL + ";")
		re2 := regexp.MustCompile(`path=(.*?);`)
		valStr = re2.ReplaceAllString(valStr, "path=" + global.Prefix + base64.URLEncoding.EncodeToString([]byte(global.ProxyURL)) + "/" + ";")
	}

	val = strings.Split(valStr, "; ")

	return val
}

func elmAttr(key string, val string) string {
	if key == "href" || key == "src" || key == "poster" || key == "data" || key == "action" || key == "srcset" || key == "data-src" || key == "data-href" {
		attrURL, err := url.Parse(val)
		if err != nil || attrURL.Scheme == "" || attrURL.Host == "" {
			if val != "" {
				val = global.Scheme + "//" + global.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(global.ProxyURL + val[1:]))
 			}
		} else {
			val = global.Scheme + "//" + global.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(val))
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
		
			if token.Data == "head" {
				out += "<script src=\"../js/inject.js\" data-config=\"" + base64.URLEncoding.EncodeToString([]byte("{\"url\":\"" + global.ProxyURL + "\"}")) + "\"></script>"
			}
			if token.Data == "style" {
				// TODO: Send this to CSS rewrite function (Do what eli told me to do)
			}
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
		tokenType, token := tokenizer.Next()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tokenType {
		case css.StringToken:
			data := strings.Replace(string(token), "'", "", 1)
			data = strings.Replace(string(data), "'", "", 1)

			url, err := url.Parse(data)
			if err != nil || url.Scheme == "" || url.Host == "" {
				data = global.Scheme + "//" + global.Host + global.Prefix + base64.StdEncoding.EncodeToString([]byte(global.URL + data))
			} else if strings.HasPrefix(data, "http://") || strings.HasPrefix(data, "https://") {
				data = global.Scheme + "//" + global.Host + global.Prefix + base64.StdEncoding.EncodeToString([]byte(data))
			}

			out += data
		case css.URLToken:
			data := strings.Replace(string(token), "url(", "", 4)
			data = strings.Replace(string(data), ")", "", 1)
		
			url, err := url.Parse(data)
			if err != nil || url.Scheme == "" || url.Host == "" {
				data = global.Scheme + global.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(global.URL + data))
			} else {
				data = global.Scheme + global.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(data))
			}
	
			out += "url(" + data + ")"
		default:
			out += string(token)
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