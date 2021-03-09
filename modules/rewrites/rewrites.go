package rewrites

import (
	"github.com/titaniumnetwork-dev/Aurora/modules/config"
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
	"log"
)

// TODO: Write a proper header parser
func Header(key string, valArr []string) []string {
	val := strings.Join(valArr, "; ")

	// TODO: Continue adding more header rewrites
	switch key {
	case "Set-Cookie":
		re1 := regexp.MustCompile(`domain=(.*?);`)
		val = re1.ReplaceAllString(val, "domain=" + config.URL.Hostname() + ";")
		re2 := regexp.MustCompile(`path=(.*?);`)
		val = re2.ReplaceAllString(val, "path=" + config.Prefix + base64.URLEncoding.EncodeToString([]byte(config.ProxyURL.String())) + "/" + ";")
	}

	valArr = strings.Split(val, "; ")

	return valArr
}

func internalHTML(key string, val string) (string, error) {
	if key == "href" || key == "src" || key == "poster" || key == "data" || key == "action" || key == "srcset" || key == "data-src" || key == "data-href" {
		url, err := url.Parse(val)
		if err != nil || url.Scheme == "" || url.Host == "" {
			if val != "" {
				val = config.URL.Scheme + "://" + config.URL.Host + config.Prefix + base64.URLEncoding.EncodeToString([]byte(config.ProxyURL.String() + val[1:]))
			} else {
				err = errors.New("No value in attribute" + key + "set")
				return "", nil
			}
		} else {
			val = config.URL.Scheme + "://" + config.URL.Host + config.Prefix + base64.URLEncoding.EncodeToString([]byte(val))
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
		val = config.URL.Scheme + "://" + config.URL.Host + config.Prefix + base64.URLEncoding.EncodeToString([]byte(config.URL.String() + val))
	} else {
		val = config.URL.Scheme + "://" + config.URL.Host + config.Prefix + base64.URLEncoding.EncodeToString([]byte(val))
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
			return nil, nil
		}

		switch tokenType {
		case html.TextToken:
			if string(tokenizer.Text()) == "style" {
				valInterface, err := CSS(token.Data)
				if err != nil {
					return nil, err
				}
				val := valInterface.(string)
				token.Data = val
			}
			out += token.Data
		case html.StartTagToken:
			attr := ""
			for _, elm := range token.Attr {
				attrSel, err := internalHTML(elm.Key, elm.Val)
				if err == nil {
					attr += attrSel
				} else {
					log.Println(err)
				}
			}

			out += "<" + token.Data + attr + ">"

			if token.Data == "head" {
				out += "<script src=\"../js/inject.js\" data-config=\"" + base64.URLEncoding.EncodeToString([]byte("{\"url\":\"" + config.ProxyURL.String() + "\"}")) + "\"></script>"
			}
		case html.EndTagToken:
			out += token.String()
		case html.SelfClosingTagToken:
			attr := ""
			for _, elm := range token.Attr {
				attrSel, err := internalHTML(elm.Key, elm.Val)
				if err == nil {
					attr += attrSel
				} else {
					return nil, err
				}
			}

			out += "<" + token.Data + attr + "/>"
		case html.CommentToken:
			out += token.String()
		case html.DoctypeToken:
			out += token.String()
		}
	}

	body = ioutil.NopCloser(strings.NewReader(out))
	body.Close()
	return body, nil
}

func CSS(bodyInterface interface{}) (interface{}, error) {
	var tokenizer css.Lexer
	switch bodyInterface.(type) {
	case io.ReadCloser:
		body := bodyInterface.(io.ReadCloser)
		tokenizer := css.NewLexer(parse.NewInput(body))
	case string:
		body := bodyInterface.(string)
		tokenizer := css.NewLexer(parse.NewInput(strings.NewReader(body)))
	default:
		err := errors.New("Invalid argument type passed")
		return nil, err
	}

	out := ""

	for {
		tokenType, token := tokenizer.Next()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		} else if err != nil {
			switch bodyInterface.(type) {
			case io.ReadCloser:
				return nil, err
			case string:
				return "", err
			}
		}

		tokenStr := string(token)
		switch tokenType {
		case css.AtKeywordToken:
			val := strings.Replace(tokenStr, "'", "", 1)
			val = strings.Replace(val, "'", "", 1)
			val = internalCSS(val)

			out += val
		case css.URLToken:
			val := strings.Replace(tokenStr, "url(", "", 4)
			val = strings.Replace(val, ")", "", 1)
			val = internalCSS(val)

			out += "url(" + val + ")"
		default:
			out += tokenStr
		}
	}

	switch bodyInterface.(type) {
	case string:
		return out, nil
	default:
		body := ioutil.NopCloser(strings.NewReader(out))
		body.Close()
		return body, nil
	}
}

// Low Priority

// TODO: Add xml rewrites for external entities
// Use https://golang.org/pkg/encoding/xml/

// TODO: Add svg rewrites
// Use https://github.com/rustyoz/svg/
