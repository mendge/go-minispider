package spider

import (
	"fmt"
	"minispider/internal/util"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"regexp"
	"sync"
	"testing"
	"time"
)

// same as TestTaskQueue_PushTask
func TestSpider_AddTask(t *testing.T) {
	// the difference is just use an interface Scheduler implemented by &Spider{} to call function
}

func TestSpider_DownloadHTML(t *testing.T) {
	const defaultTimeout = 10
	httpClient := &http.Client{
		Timeout: time.Duration(defaultTimeout) * time.Second,
	}

	tests := []struct {
		name    string
		urlStr  string
		wantErr bool
	}{
		{"200", "http://www.baidu.com", false},
		{"404", "http://www.baidubaidu.com", true},
		{"timeout", "https://www.google.com", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var downloader Downloader = NewSpider(nil, httpClient, nil, nil)
			body, err := downloader.DownloadHTML(httpClient, test.urlStr)
			if err != nil {
				t.Errorf("%s\n", err)
				return
			}
			if len(body) == 0 {
				t.Errorf("DownloadHTML() got body len %v, want len != 0", len(body))
			}
		})
	}
}

// same as TestTaskQueue_PopTask
func TestSpider_FetchTask(t *testing.T) {
	// the difference is just use an interface Scheduler implemented by &Spider{} to call function
}

func TestSpider_ParseSubURLs(t *testing.T) {
	re, _ := regexp.Compile(".*.(htm|html)$")
	tests := []struct {
		name      string
		faURL     string
		wantLinks []string
		html      string
		wantErr   bool
	}{
		{
			"sibling links",
			"https://www.example.com",
			[]string{
				"https://www.example.com/aaa/bbb/page1.html",
				"https://www.example.com/page2.html",
				"https://www.example.com/page3.html",
			},
			`<!DOCTYPE html>
<html lang="en">
<body>
    <h1>Test Page</h1>
    <ul>
        <li><a href="https://www.example.com/aaa/bbb/page1.html">Page 1</a></li>
        <li><a href="../../page2.html">Page 2</a></li>
        <li><a href="/page3.html">Page 3</a></li>
        <li><a href="https://www.example.com/page3">Page 4</a></li>
    </ul>
</body>
</html>`,

			false,
		},
		{
			"Nested links",
			"https://mendge.com",
			[]string{
				"https://mendge.com/xxx/yyy/page1.html",
				"https://mendge.com/page2.html",
				"https://mendge.com/page3.html",
			},
			`<!DOCTYPE html>
<html lang="en">
<head>
</head>
<body>
    <h1>Nested Links Page</h1>
    <ul>
        <li><a href="https://mendge.com/xxx/yyy/page1.html">Page 1</a></li>
        <li>
            <a href="../../page2.html">Page 2</a>
            <ul>
                <li><a href="/page3.html">page 3</a></li>
				<ul>
                	<li><a href="../../image.png"> image </a></li>
            	</ul>
            </ul>
        </li>
    </ul>
</body>
</html>`,
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var processor Processor = &Spider{}
			links, err := processor.ParseSubURLs(test.faURL, []byte(test.html), re)
			if err != nil {
				t.Errorf("parseSubURLs() got error:%s\n", err)
				return
			}
			if !reflect.DeepEqual(test.wantLinks, links) {
				t.Errorf("parseSubURLs() got links %v, want %v\n", links, test.wantLinks)
			}
		})
	}
}

func TestSpider_SaveToFile(t *testing.T) {
	html := "<!doctype html>\n<html>\n<head>\n    <title>Example Domain</title>\n\n    <meta charset=\"utf-8\" />\n    <meta http-equiv=\"Content-type\" content=\"text/html; charset=utf-8\" />\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n    <style type=\"text/css\">\n    body {\n        background-color: #f0f0f2;\n        margin: 0;\n        padding: 0;\n        font-family: -apple-system, system-ui, BlinkMacSystemFont, \"Segoe UI\", \"Open Sans\", \"Helvetica Neue\", Helvetica, Arial, sans-serif;\n        \n    }\n    div {\n        width: 600px;\n        margin: 5em auto;\n        padding: 2em;\n        background-color: #fdfdff;\n        border-radius: 0.5em;\n        box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);\n    }\n    a:link, a:visited {\n        color: #38488f;\n        text-decoration: none;\n    }\n    @media (max-width: 700px) {\n        div {\n            margin: 0 auto;\n            width: auto;\n        }\n    }\n    </style>    \n</head>\n\n<body>\n<div>\n    <h1>Example Domain</h1>\n    <p>This domain is for use in illustrative examples in documents. You may use this\n    domain in literature without prior coordination or asking for permission.</p>\n    <p><a href=\"https://www.iana.org/domains/example\">More information...</a></p>\n</div>\n</body>\n</html>"
	htmlData := []byte(html)
	type args struct {
		urlStr    string
		outputDir string
		data      []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"right occasion", args{"www.example.com", "../testdata/", htmlData}, false},
		{"nonexistent outDir", args{"www.example.com", "../testdata/testDir/", htmlData}, false},
		{"unescaped url", args{"https://www.example.com/", "../testdata/", htmlData}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outputDirNotExist := false
			_, err := os.Stat(test.args.outputDir)
			if os.IsNotExist(err) {
				outputDirNotExist = true
			}
			var persistencer Persistencer = &Spider{}
			err = persistencer.SaveToFile(test.args.urlStr, test.args.outputDir, htmlData)
			if err != nil {
				t.Errorf("%s", err)
				return
			}
			// remove file or directory been created
			fileName := url.QueryEscape(util.RemoveProtocol(test.args.urlStr))
			_ = os.Remove(path.Join(test.args.outputDir, fileName))
			if outputDirNotExist {
				_ = os.Remove(test.args.outputDir)
			}
		})
	}
}

// TODO
func TestTaskQueue_PopTask(t *testing.T) {
	//const initialTaskCount = 2
	//const multiples = 2 // 倍数
	//const concurrencyCount = multiples * initialTaskCount
	//
	//var allTasksDone = make(chan struct{})
	//taskQ := NewTaskDeque(concurrencyCount)
	//i := 0
	//
	//// 相当于初始化任务队列的任务
	//for ; i < initialTaskCount; i++ {
	//	urlStr := fmt.Sprintf("Task[ %d ]", i)
	//	taskQ.PushTask(NewTask(1, urlStr))
	//	t.Logf("Push %s\n", urlStr)
	//}
	//
	//var produceAndConsume func(id int, onlyConsume bool)
	//produceAndConsume = func(id int, onlyConsume bool) {
	//	task, noMoreTask := taskQ.PopTask()
	//	if noMoreTask {
	//		allTasksDone <- struct{}{}
	//		return
	//	}
	//	if onlyConsume {
	//		// 概率只消费不生产，使得最后刚好消费完
	//		t.Logf("Pop %s only", task.DstURL)
	//	} else {
	//		// 消费完正常生产，填补消费的那个任务
	//		newTask := NewTask(1, fmt.Sprintf("Task[ %d ]", id))
	//		taskQ.PushTask(newTask)
	//		t.Logf("Pop %s \t\t Push Task[ %d ]", task.DstURL, id)
	//	}
	//	t.Logf("\t\tLen of deque %d", taskQ.dq.Len())
	//}
	//
	////
	//for ; i <= initialTaskCount+concurrencyCount; i++ {
	//	go produceAndConsume(i, i%multiples == 0)
	//}
	//go produceAndConsume(i, true)
	//
	//select {
	//case <-allTasksDone:
	//	if taskQ.dq.Len() != concurrencyCount {
	//		t.Errorf("In the case of %d coroutines popping and pushing task concurrently,"+
	//			"when detected no more task, the expected queue length is %d,"+
	//			" but the actual length is %d\n", concurrencyCount, initialTaskCount, taskQ.dq.Len())
	//		return
	//	}
	//}
}

func TestTaskQueue_PushTask(t *testing.T) {
	const concurrencyCount = 1000000
	wg := sync.WaitGroup{}
	wg.Add(concurrencyCount)
	taskQ := NewTaskDeque(concurrencyCount)
	for i := 0; i < concurrencyCount; i++ {
		go func(id int) {
			task := NewTask(1, fmt.Sprintf("Task[ %d ]", id))
			taskQ.PushTask(task)
			//t.Logf("Push Task[ %d ]\n", id)
			wg.Done()
		}(i)
	}
	wg.Wait()
	if taskQ.dq.Len() != concurrencyCount {
		t.Errorf("In the case of %d coroutines pushing task concurrently, the expected queue length is %d,"+
			" but the actual length is %d\n", concurrencyCount, concurrencyCount, taskQ.dq.Len())
		return
	}
}

func TestUniMap_IsKeyVisitedOrVisit(t *testing.T) {
	const count = 10003
	const circleSpan = 5
	const basedURL = "http://www.example.com/image"
	uniMap := UniMap{}
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(id int) {
			urlStr := fmt.Sprintf("%s_%d", basedURL, id/circleSpan)
			uniMap.IsKeyVisitedOrVisit(urlStr)
			wg.Done()
		}(i)
	}
	wg.Wait()
	mapSize := 0
	uniMap.urls.Range(func(key, value interface{}) bool {
		mapSize++
		return true
	})
	// the mathematical expression is to calculate the different value of i/circleSpan (i from 0 to count-1)
	wantSize := (count-1)/circleSpan + 1
	if mapSize != wantSize {
		t.Errorf("In the case of coroutines working concurrently, the expected uniMap size is %d,"+
			" but the actual size  is %d\n", wantSize, mapSize)
	}
}
