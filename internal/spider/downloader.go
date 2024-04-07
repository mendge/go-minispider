package spider

import (
	"errors"
	"io"
	"net/http"
)

type Downloader interface {
	// DownloadHTML downloads the page of specified url.
	DownloadHTML(httpClient *http.Client, urlStr string) ([]byte, error)
}

func (s *Spider) DownloadHTML(httpClient *http.Client, urlStr string) ([]byte, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, errors.New("fail to create a request")
	}
	// TODO 通用化配置请求
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := httpClient.Do(req)
	// Close的正确使用姿势
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.New("fail to request url")
	}

	pageBody, err := io.ReadAll(resp.Body)
	return pageBody, err
}
