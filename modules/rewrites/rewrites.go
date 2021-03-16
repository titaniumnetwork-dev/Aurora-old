package rewrites

import (
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
	"github.com/titaniumnetwork-dev/Aurora/modules/config"
	"golang.org/x/net/html"

	//	"encoding/xml"

	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
)

// TODO: Start switching to using fmt.Sprintf()

var err error

// TODO: Continue adding more header rewrites until it's done
func Header(key string, vals []string) []string {
	for i, val1 := range vals {
		split1 := strings.Split(val1, "; ")
		for j, val2 := range split1 {
			switch key {
			// Request headers
			case "Host":
				split2 := strings.Split(val2, ":")
				split2[0] = config.ProxyURL.Host
				val2 = strings.Join(split2, ":")
			// Response headers
			case "Set-Cookie":
				split2 := strings.Split(val2, "=")
				switch split2[0] {
				case "domain":
					split2[1] = config.URL.Hostname()
				case "path":
					split2[1] = config.YAML.HTTPPrefix + base64.URLEncoding.EncodeToString([]byte(config.ProxyURL.String()))
				}
				val2 = strings.Join(split2, "=")
			}
			split1[j] = val2
		}
		val1 = strings.Join(split1, "; ")
		vals[i] = val1
	}

	return vals
}

func internalHTML(key string, val string) string {
	if key == "href" || key == "src" || key == "poster" || key == "action" || key == "srcset" {
		url, err := url.Parse(val)
		if err != nil || url.Scheme == "" || url.Host == "" {
			split := strings.Split(val, "../")
			if len(split) >= 2 {
				val = split[len(split)-1]

				pSplit := strings.Split(config.ProxyURL.Path, "/")
				for i := 1; i <= len(pSplit); i++ {
					pSplit[len(pSplit)-i] = ""
				}

				val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s://%s%s/%s", config.ProxyURL.Scheme, config.ProxyURL.Host, config.ProxyURL.Path[:len(strings.Join(pSplit, ""))], val))))
			} else if strings.HasPrefix(val, "/") {
				val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s://%s%s", config.ProxyURL.Scheme, config.ProxyURL.Host, val))))
			} else {
				val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(config.ProxyURL.String()+val)))
			}
		} else {
			val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(val)))
		}
	}
	if key == "style" {
		valInterface := CSS(val)
		val = valInterface.(string)
	}
	attr := fmt.Sprintf(" %s=\"%s\"", key, val)
	return attr
}

func internalCSS(val string) string {
	url, err := url.Parse(val)
	if err != nil || url.Scheme == "" || url.Host == "" {
		split := strings.Split(val, "../")
		if len(split) >= 2 {
			val = split[len(split)-1]

			pSplit := strings.Split(config.ProxyURL.Path, "/")
			for i := 1; i <= len(pSplit); i++ {
				pSplit[len(pSplit)-i] = ""
			}

			val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s://%s%s/%s", config.ProxyURL.Scheme, config.ProxyURL.Host, config.ProxyURL.Path[:len(strings.Join(pSplit, ""))], val))))
		} else if strings.HasPrefix(val, "/") {
			val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s://%s%s", config.ProxyURL.Scheme, config.ProxyURL.Host, val))))
		} else {
			val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(config.ProxyURL.String()+val)))
		}
	} else {
		val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(val)))
	}

	return val
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
			if string(tokenizer.Text()) == "style" {
				valInterface := CSS(token.Data)
				val := valInterface.(string)
				token.Data = val
			}
			out += token.Data
		case html.StartTagToken:
			attr := ""
			for _, elm := range token.Attr {
				attrSel := internalHTML(elm.Key, elm.Val)
				attr += attrSel
			}

			out += fmt.Sprintf("<%s%s>", token.Data, attr)

			if token.Data == "head" {
				out += fmt.Sprintf("<script src=\"../js/inject.js\" data-config=\"%s\"></script>", base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("{\"url\":\"%s\",\"proxyurl\":\"%s\",\"httpprefix\":\"%s\",\"wsprefix\":\"%s\"}", config.URL.String(), config.ProxyURL.String(), config.YAML.HTTPPrefix, config.YAML.WSPrefix))))
			}
			if token.Data == "html" {
				// Temporary solution
				// TODO: Doesn't work; fix
				// token.Attr = append("id", "domsel")
			}
		case html.SelfClosingTagToken:
			attr := ""
			for _, elm := range token.Attr {
				attrSel := internalHTML(elm.Key, elm.Val)
				attr += attrSel
			}

			out += fmt.Sprintf("<%s%s/>", token.Data, attr)
		default:
			out += token.String()
		}
	}

	body = ioutil.NopCloser(strings.NewReader(out))
	body.Close()
	return body
}

func CSS(bodyInterface interface{}) interface{} {
	var tokenizer *css.Lexer
	switch bodyInterface.(type) {
	case string:
		body := bodyInterface.(string)
		tokenizer = css.NewLexer(parse.NewInput(strings.NewReader(body)))
	default:
		body := bodyInterface.(io.ReadCloser)
		tokenizer = css.NewLexer(parse.NewInput(body))
	}

	out := ""

	for {
		tokenType, token := tokenizer.Next()

		err = tokenizer.Err()
		if err == io.EOF {
			break
		}

		tokenStr := string(token)
		switch tokenType {
		case css.URLToken:
			val := strings.Replace(tokenStr, "url(", "", 4)
			val = strings.Replace(val, ")", "", 1)
			val = internalCSS(val)

			out += fmt.Sprintf("url(%s)", val)
		default:
			out += tokenStr
		}
	}

	switch bodyInterface.(type) {
	case string:
		return out
	default:
		body := ioutil.NopCloser(strings.NewReader(out))
		body.Close()
		return body
	}
}

// Low Priority

// TODO: Add xml rewrites for external entities
// Use https://golang.org/pkg/encoding/xml/

// TODO: Add svg rewrites
// Use https://github.com/rustyoz/svg/
