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
	"errors"
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
		valStr = re2.ReplaceAllString(valStr, "path=" + global.Prefix + base64.URLEncoding.EncodeToString([]byte(global.ProxyURL.String())) + "/" + ";")
	}

	val = strings.Split(valStr, "; ")

	return val
}

func internalHTML(key string, val string) (string, error) {
	if key == "href" || key == "src" || key == "poster" || key == "data" || key == "action" || key == "srcset" || key == "data-src" || key == "data-href" {
		attrURL, err := url.Parse(val)
		if err != nil || attrURL.Scheme == "" || attrURL.Host == "" {
			if val != "" {
				val = global.URL.Scheme + "//" + global.URL.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(global.ProxyURL.String() + val[1:]))
			}
		} else {
			val = global.URL.Scheme + "//" + global.URL.Host + global.Prefix + base64.URLEncoding.EncodeToString([]byte(val))
		}
	}
	if key == "style" {
		val, err := CSS(val)
		if err != nil {
			return "", err
		}
	}
	attr := " " + key + "=" + "\"" + val + "\""
	return attr, nil
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
			tagnameBytes, _ := tokenizer.TagName()
			tagname := string(tagnameBytes)
			if tagname == "style" {
				valInterface, err := CSS(token.Data) 
				val := valInterface.(string)
				if err == nil {
					token.Data = val
				} else {
					return nil, err
				}
			}
			out += token.Data
		case html.StartTagToken:
			attr := ""
			for _, elm := range token.Attr {
				attrTemp, err := internalHTML(elm.Key, elm.Val)
				attr += attrTemp
				if err != nil {
					return nil, err
				}
			}

			out += "<" + token.Data + attr + ">"

			if token.Data == "head" {
				out += "<script src=\"../js/inject.js\" data-config=\"" + base64.URLEncoding.EncodeToString([]byte("{\"url\":\"" + global.ProxyURL.String() + "\"}")) + "\"></script>"
			}
		case html.EndTagToken:
			out += "</" + token.Data + ">"
		case html.SelfClosingTagToken:
			attr := ""
			for _, elm := range token.Attr {
				attrTemp, err := internalHTML(elm.Key, elm.Val)
				if err == nil {
					attr += attrTemp
				} else {
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

func CSS(bodyInterface interface{}) (interface{}, error) {
	switch bodyInterface.(type) {
	case io.ReadCloser:
		body := bodyInterface.(io.ReadCloser)
		tokenizer := css.NewLexer(parse.NewInput(body))
	case string:
		body := bodyInterface.(string)
		// There might be a more efficient method
		tokenizer := css.NewLexer(parse.NewInput(strings.NewReader(body)))
	default:
		err := errors.New("Invalid argument type passed to CSS function")
		return nil, err
	}

	out := ""

	for {
		tokenType, token := tokenizer.Next()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		} else {
			switch bodyInterface.(type) {
			case io.ReadCloser:
				return nil, err
			case string:
				return "", err
			}
		}

		switch tokenType {
		case css.AtKeywordToken:
			val := strings.Replace(string(token), "'", "", 1)
			val = strings.Replace(string(val), "'", "", 1)
			val = internalCSS(val)

			out += val
		case css.URLToken:
			val := strings.Replace(string(token), "url(", "", 4)
			val = strings.Replace(string(val), ")", "", 1)
			val = internalCSS(val)

			out += "url(" + val + ")"
		default:
			out += string(token)
		}
	}

	switch bodyInterface.(type) {
	case io.ReadCloser:
		body = bodyInterface.(io.ReadCloser)
		body = ioutil.NopCloser(strings.NewReader(out))
		body.Close()
		return body, nil
	case string:
		return out, nil
	}
}

// Low Priority

// TODO: Add xml rewrites for external entities
// Use https://golang.org/pkg/encoding/xml/

// TODO: Add svg rewrites
// Use https://github.com/rustyoz/svg/
