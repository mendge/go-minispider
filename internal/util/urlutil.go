package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strings"
)

import (
	"golang.org/x/net/html"
)

// ParseSeeds parses url seeds from file of seedFilePath. error will return if anything wrong.
func ParseSeeds(seedFilePath string) ([]string, error) {
	file, err := os.Open(seedFilePath)
	if err != nil {
		return nil, errors.New("please give a right seeds file")
	}
	defer file.Close()

	var urls []string
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&urls); err != nil {
		return nil, errors.New("fail to parse seeds file in json format")
	}
	return urls, nil
}

// RemoveProtocol removes the protocol in urlStr.
func RemoveProtocol(urlStr string) string {
	index := strings.Index(urlStr, "://")
	if index == -1 {
		return urlStr
	}
	return urlStr[index+3:]
}

// ExtractLinks extracts all links in page
func ExtractLinks(faURL string, pageContent []byte) ([]string, error) {
	doc, err := html.Parse(bytes.NewReader(pageContent))
	if err != nil {
		return nil, errors.New("fail to parse html")
	}

	baseURL, err := url.Parse(faURL)
	if err != nil {
		return nil, err
	}

	var links []string
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					// 相对路径和绝对路径转换
					linkURL, err := baseURL.Parse(attr.Val)
					if err != nil {
						continue
					}
					links = append(links, linkURL.String())
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(doc)
	return links, err
}
