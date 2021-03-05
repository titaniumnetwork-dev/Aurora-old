package rewrites

import (
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/config"
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
	case "Set-Cookie":
		re1 := regexp.MustCompile(`domain=(.*?);`)
		valStr = re1.ReplaceAllString(valStr, "domain=" + global.URL.Hostname() + ";")
		re2 := regexp.MustCompile(`path=(.*?);`)
		valStr = re2.ReplaceAllString(valStr, "path=" + global.Prefix + base64.URLEncoding.EncodeToString([]byte(global.ProxyURL)) + "/" + ";")
	}

	val = strings.Split(valStr, "; ")

	return val
}

func internalHTML(key string, val string) (string, error) {
	if key == "href" || key == "src" || key == "poster" || key == "data" || key == "action" || key == "srcset" || key == "data-src" || key == "data-href" {
		attrURL, err := url.Parse(val)
		if err != nil || attrURL.Scheme == "" || attrURL.Host == "" {
			if val != "" {
				val = global.URL.Scheme + "//" + global.URL.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(global.ProxyURL + val[1:]))
			}
		} else {
			val = global.URL.Scheme + "//" + global.URL.Host + global.URL.Prefix + base64.URLEncoding.EncodeToString([]byte(val))
		}
	} else if key == "style" {
		val, err := CSS(val)
		if err != nil {
			return nil, error
		}
	}
	attr := " " + key + "=" + "\"" + val + "\""
	return attr
}

func internalCSS(val string) string {
	url, err := url.Parse(val)
	if err != nil || url.Scheme == "" || url.Host == "" {
		val = global.URL.Scheme + "//" + global.URL.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(global.URL.String() + val))
	} else if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
		val = global.URL.Scheme + "//" + global.URL.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(val))
	}

	return val
}

func HTML(body io.ReadCloser) (io.ReadCloser, error) {
	tokenizer := html.NewTokenizer(body)
	out := ""

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		} else {
			return nil, err
		}

		switch tokenType {
		case html.TextToken:
			out += token.Data
		case html.StartTagToken:
			attr := ""
			for _, elm := range token.Attr {
				attr, err += internalHTML(elm.Key, elm.Val)
				if err != nil {
					return nil, err
				}
			}

			out += "<" + token.Data + attr + ">"

			if token.Data == "head" {
				out += "<script src=\"../js/inject.js\" data-config=\"" + base64.URLEncoding.EncodeToString([]byte("{\"url\":\"" + global.ProxyURL + "\"}")) + "\"></script>"
			}
			if token.Data == "style" {
				val, err := CSS(token.Text())
				if err != nil {
					return nil, err
				}
				out += "<style>" + val + "</script>"
			}
		case html.EndTagToken:
			out += "</" + token.Data + ">"
		case html.SelfClosingTagToken:
			attr := ""
			for _, elm := range token.Attr {
				attr, err += internalHTML(elm.Key, elm.Val)
				if err != nil {
					return nil, err
				}
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
	return body, nil
}

func CSS(body interface{}) (interface{}, error) {
	if body.(type) == io.ReadCloser {
		tokenizer := css.NewLexer(parse.NewInput(body))
	} else if body.(type) == string {
		tokenizer := css.NewLexer(strings.NewReader(body))
	} else {
		return nil, errors.New("Invalid argument type passed to CSS function " + body)
	}

	out := ""

	for {
		tokenType, token := tokenizer.Next()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		} else if {
			return err
		}

		switch tokenType {
		case css.AtKeywordToken:
			val := strings.Replace(string(token), "'", "", 1)
			val = strings.Replace(string(val), "'", "", 1)
			val = internalCSS(val)

			out += val
		case css.URLToken:
			val := strings.Replace(string(token), "url(", "", 4)
			val = strings.Replace(string(data), ")", "", 1)
			val = internalCSS(val)

			out += "url(" + val + ")"
		default:
			out += string(token)
		}
	}

	if body.(type) == io.ReadCloser {
		body = ioutil.NopCloser(strings.NewReader(out))
		body.Close()
		return body, nil
	} else if body.(type) == string {
		return out, nil
	}
}

// Low Priority

// TODO: Add xml rewrites for external entities
// Use https://golang.org/pkg/encoding/xml/

// TODO: Add svg rewrites
// Use https://github.com/rustyoz/svg/
