package netutils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpRequest struct {
	url         string
	method      string
	header      map[string]string
	body        io.Reader
	contentType string
}

func NewHttpRequest(url string) (h *HttpRequest) {
	h = &HttpRequest{
		url:         url,
		method:      "GET",
		header:      make(map[string]string),
		body:        nil,
		contentType: "",
	}
	return h
}

func (h *HttpRequest) SetHeader(k, v string) {
	h.header[k] = v
}

func (h *HttpRequest) SetMethod(m string) error {
	m = strings.ToUpper(m)
	if m != "GET" && m != "POST" {
		return fmt.Errorf("only support GET or POST")
	}
	h.method = m
	return nil
}

func (h *HttpRequest) SetBodyStr(b string, contentType string) {
	h.body = strings.NewReader(b)
	h.contentType = contentType
}

func (h *HttpRequest) SetBodyBytes(b []byte, contentType string) {
	h.body = bytes.NewReader(b)
	h.contentType = contentType
}

func (h *HttpRequest) SetBodyFields(fields map[string]string) {
	h.contentType = "application/x-www-form-urlencoded"
	buf := bytes.NewBufferString("")
	for k, v := range fields {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(v))
	}
	h.body = buf

}

func (h *HttpRequest) Resp(timeout time.Duration) (resp *http.Response, err error) {
	client := http.Client{}
	client.Timeout = timeout
	req, _ := http.NewRequest(h.method, h.url, h.body)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	if len(h.header) > 0 {
		for k, v := range h.header {
			req.Header.Set(k, v)
		}
	}

	if len(h.contentType) > 0 {
		req.Header.Set("Content-Type", h.contentType)
	}

	resp, err = client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *HttpRequest) ExecReader(timeout time.Duration) (reader io.Reader, resp *http.Response, err error) {

	resp, err = h.Resp(timeout)

	if err != nil {
		return nil, resp, err
	}

	if strings.Compare("gzip", resp.Header.Get("Content-Encoding")) == 0 {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				panic(err)
			}
		}()

		r, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, resp, err
		}

		return r, resp, nil
	}

	return resp.Body, resp, nil
}

func (h *HttpRequest) Exec(timeout time.Duration) (result []byte, resp *http.Response, err error) {

	resp, err = h.Resp(timeout)

	if err != nil {
		return nil, resp, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if strings.Compare("gzip", resp.Header.Get("Content-Encoding")) == 0 {
		r, err := gzip.NewReader(bytes.NewReader(content))
		if err != nil {
			return nil, resp, err
		}

		rs, err := io.ReadAll(r)
		if err != nil {
			return nil, resp, err
		}
		result = rs
	} else {
		result = content
	}

	return result, resp, nil
}
