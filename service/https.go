package service

import (
	"errors"
	"io"
	"net/http"
	"utilities/uerror"
)

type httpReq struct {
}

var HttpReq httpReq

type UrlUploadImage struct {
	Url string `json:"url"`
}

func (r httpReq) CrawlByURL(method, url string) (body io.Reader, status int, err error) {
	client := http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		err = uerror.StackTrace(err)
		return
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		err = uerror.StackTrace(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		status = resp.StatusCode
		err = uerror.StackTrace(errors.New("Status crawl fail " + resp.Status))
		return
	}
	body = resp.Body
	status = resp.StatusCode
	return
}
