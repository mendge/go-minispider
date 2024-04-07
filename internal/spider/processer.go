package spider

import (
	"regexp"
)

import (
	"minispider/internal/util"
)

type Processor interface {
	// ParseSubURLs parses the subURLs in pageBody with specified regular expression.
	ParseSubURLs(faURL string, pageBody []byte, re *regexp.Regexp) ([]string, error)
}

func (s *Spider) ParseSubURLs(faURL string, pageBody []byte, re *regexp.Regexp) ([]string, error) {
	var subURLs []string
	links, err := util.ExtractLinks(faURL, pageBody)
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		if re.MatchString(link) {
			subURLs = append(subURLs, link)
		}
	}
	return subURLs, nil
}
