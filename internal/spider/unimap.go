package spider

import "sync"

type UniMap struct {
	urls sync.Map
}

// IsKeyVisitedOrVisit check whether urlStr has been visited, mark it as visited for next time if not.
func (um *UniMap) IsKeyVisitedOrVisit(urlStr string) bool {
	_, ok := um.urls.LoadOrStore(urlStr, struct{}{})
	return ok
}
