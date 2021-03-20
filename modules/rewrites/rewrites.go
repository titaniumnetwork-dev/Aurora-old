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

// Here contains some debug logging which would be removed after bugs are fixed
func URL(val string) string {
	url, err := url.Parse(val)
	if err != nil || url.Scheme == "" || url.Host == "" {
		pURL := config.ProxyURL.String()

		split1_0 := strings.Split(val, "../")
		split1_1 := strings.Split(config.ProxyURL.String(), "/")
		split1_2 := strings.Split(split1_1[len(split1_1)-1], ".")

		split2 := strings.Split(val, ":")

		if len(split1_2) >= 2 {
			split1_1[len(split1_1)-1] = ""
			if len(split1_0) >= 2 {
				split1_1[len(split1_1)-2] = ""
			}
			pURL = strings.Join(split1_1, "/")
		}

		//log.Println("URL Invalid")
		switch true {
		case len(split2) == 2:
			//log.Println("Protocol url: " + val)
		case len(split1_0) >= 2:
			val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s%s", pURL, val))))
			//log.Println("../ url: " + val)
		case strings.HasPrefix(val, "//"):
			val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(pURL+val)))
			//log.Println("// url: " + val)
		case strings.HasPrefix(val, "/"):
			val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s://%s%s", config.ProxyURL.Scheme, config.ProxyURL.Host, val))))
			//log.Println("/ url: " + val)
		default:
			val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(pURL+val)))
			//log.Println("url: " + val)
		}
	} else {
		//log.Println("URL Valid")
		val = fmt.Sprintf("%s://%s%s%s", config.Scheme, config.URL.Host, config.YAML.HTTPPrefix, base64.URLEncoding.EncodeToString([]byte(val)))
		//log.Println("url: " + val)
	}

	return val
}

// TODO: Send multiple key and values in at the same time because for example how would it know about this?
// <meta http-equiv="refresh" content="3;url=https://www.mozilla.org">
// See https://www.w3schools.com/Tags/att_object_usemap.asp
// Add usemap
func internalHTML(key string, val string, tag string) string {
	switch true {
	// See https://stackoverflow.com/questions/28652648/how-to-use-external-svg-in-html and http://bl.ocks.org/clhenrick/0b73208409a14144e1f5
	case key == "href" || key == "src" || key == "poster" || key == "action" || key == "formaction" || key == "data":
		val = URL(val)
	case key == "srcset":
		// As seen from google logo when proxified the file extension is split too!
		split := strings.Split(val, " ")
		// TODO: Switch to using range
		for i := 0; i <= len(split)-1; i++ {
			if i^1 == i+1 {
				split[i] = URL(split[i])
			}
		}
		val = strings.Join(split, " ")
	case strings.HasPrefix(key, "on"):
		val = fmt.Sprintf("{let document=audocument;%s}", val)
	case key == "srcdoc":
		// TODO: Rewrite html again... why does this have to exist :(
		// I will have to make html return and take in an interface
	case key == "style":
		valInterface := CSS(val)
		val = valInterface.(string)
	case tag == "meta":
		// See https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name and https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta
		// http-equiv
		// name
	}

	attr := fmt.Sprintf(" %s=\"%s\"", key, val)
	return attr
}

// There is a bug that can be seen here with the html rewrites http://localhost/go/aHR0cDovL21hZGVieWV2YW4uY29tL29ic2N1cmUtY3BwLWZlYXR1cmVzLw==
func HTML(body io.ReadCloser) io.ReadCloser {
	tokenizer := html.NewTokenizer(body)
	out := ""

	isScript := false
	isStyle := false

	for {
		tokenType := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tokenType {
		case html.TextToken:
			switch true {
			case isScript:
				token.Data = fmt.Sprintf("{let document=audocument;%s}</script>", token.Data)
				out += token.Data
				isScript = false
			case isStyle:
				dataInterface := CSS(token.Data)
				token.Data = dataInterface.(string)
				out += token.Data
				isStyle = false
			default:
				out += token.Data
			}
		case html.StartTagToken:
			attr := ""
			for _, elm := range token.Attr {
				// See https://www.w3.org/TR/2016/REC-SRI-20160623/
				if elm.Key != "integrity" {
					// TODO: Delete directly instad
					attrSel := internalHTML(elm.Key, elm.Val, token.Data)
					attr += attrSel
				}
			}

			out += fmt.Sprintf("<%s%s>", token.Data, attr)

			switch token.Data {
			case "script":
				isScript = true
			case "style":
				isStyle = true
			case "head":
				out += "<script src=\"/inject\"></script>"
			}
		case html.SelfClosingTagToken:
			attr := ""
			for _, elm := range token.Attr {
				// TODO: Delete directly instad
				if elm.Key != "integrity" {
					attrSel := internalHTML(elm.Key, elm.Val, token.Data)
					attr += attrSel
				}
			}

			out += fmt.Sprintf("<%s%s/>", token.Data, attr)
		case html.EndTagToken:
			if token.String() == "</script>" {
				break
			}
			out += token.String()
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
			val = strings.Replace(val, "'", "", 1)
			val = strings.Replace(val, "'", "", 1)
			val = strings.Replace(val, "\"", "", 1)
			val = strings.Replace(val, "\"", "", 1)
			val = URL(val)

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

// TODO: Parse js server side and rewrite es6 imports
func JS(body io.ReadCloser) io.ReadCloser {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return body
	}
	bodyStr := fmt.Sprintf("{let document=audocument;%s}", string(bytes))
	newBody := ioutil.NopCloser(strings.NewReader(bodyStr))
	newBody.Close()
	return newBody
}

// Low Priority

// TODO: Add svg rewrites
// Use https://github.com/rustyoz/svg/
