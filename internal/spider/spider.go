package spider

import (
	"minispider/internal/config"
	"net/http"
)

type Spider struct {
	Cfg        *config.Config
	TaskQ      *TaskQueue
	HttpClient *http.Client
	UniMap     *UniMap
}

// NewSpider builds a new spider container
func NewSpider(Cfg *config.Config, client *http.Client, sq *TaskQueue, UniMap *UniMap) *Spider {
	return &Spider{
		Cfg:        Cfg,
		TaskQ:      sq,
		HttpClient: client,
		UniMap:     UniMap,
	}
}
