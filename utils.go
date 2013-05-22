/**
 * Author:        Tony.Shao
 * Email:         xiocode@gmail.com
 * Github:        github.com/xiocode
 * File:          utils.go
 * Description:   Utils.go
 */

package bingo

import (
	"compress/gzip"
	"github.com/xiocode/toolkit/to"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func encodeParams(params map[string]interface{}) (result string, err error) {
	if len(params) > 0 {
		values := url.Values{}
		for key, value := range params {
			values.Add(key, to.String(value))
		}
		result = values.Encode()
	}
	return
}

func read_body(response *http.Response) (body string, err error) {
	var reader io.ReadCloser
	var contents []byte
	using_gzip := response.Header.Get("Content-Encoding")
	switch using_gzip {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return
		}
		defer reader.Close()
	default:
		reader = response.Body
	}

	contents, err = ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	body = to.String(contents)
	return body, nil
}
