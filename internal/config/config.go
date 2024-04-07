package config

import (
	"errors"
	"regexp"
)

import (
	"github.com/go-ini/ini"
)

// Config is a struct to describe the information of configure file
type Config struct {
	UrlListFile     string
	OutputDirectory string
	MaxDepth        int
	CrawlInterval   int
	CrawlTimeout    int
	TargetUrl       string
	ThreadCount     int

	TargetUrlRE *regexp.Regexp
}

func NewConfig(urlListFile string, outputDirectory string, maxDepth int, crawlInterval int, crawlTimeout int, targetUrl string, threadCount int) *Config {
	return &Config{UrlListFile: urlListFile, OutputDirectory: outputDirectory, MaxDepth: maxDepth, CrawlInterval: crawlInterval, CrawlTimeout: crawlTimeout, TargetUrl: targetUrl, ThreadCount: threadCount}
}

// NewConfigFromFile reads the configure items from the configInputPath then build a Config struct
func NewConfigFromFile(configInputPath string) (*Config, error) {
	conf, err := ini.Load(configInputPath)
	if err != nil {
		return nil, errors.New("please give right configure file path")
	}
	spiderSection := conf.Section("spider")

	cfg := NewConfig(
		spiderSection.Key("urlListFile").String(),
		spiderSection.Key("outputDirectory").String(),
		spiderSection.Key("maxDepth").MustInt(),
		spiderSection.Key("crawlInterval").MustInt(),
		spiderSection.Key("crawlTimeout").MustInt(),
		spiderSection.Key("targetUrl").String(),
		spiderSection.Key("threadCount").MustInt(),
	)

	cfg.TargetUrlRE, err = regexp.Compile(cfg.TargetUrl)
	if err != nil {
		return nil, errors.New("fail to compile regexp in configure file")
	}
	return cfg, nil
}

// ValidateConfig checks all member variables in cfg is validated
func ValidateConfig(cfg *Config) error {
	if cfg.UrlListFile == "" {
		return errors.New("UrlListFile is empty")
	}
	if cfg.OutputDirectory == "" {
		return errors.New("OutputDirectory is empty")
	}
	if cfg.MaxDepth < 1 {
		return errors.New("MaxDepth < 1")
	}
	if cfg.CrawlInterval < 0 {
		return errors.New("CrawlInterval < 0")
	}
	if cfg.CrawlTimeout < 1 {
		return errors.New("CrawlTimeout < 1")
	}
	if cfg.TargetUrl == "" {
		return errors.New("TargetUrl is empty")
	}
	if cfg.ThreadCount < 1 {
		return errors.New("TreadCount < 1")
	}
	return nil
}
