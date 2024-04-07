package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

import (
	"minispider/internal/config"
	"minispider/internal/spider"
	"minispider/internal/util"
)

import (
	"github.com/rs/zerolog"
)

const (
	version                = "v1.0"
	defaultConfigInputPath = "./conf"
	defaultLogOutputDir    = "./log/"
	helpText               = `This is minispider, a distributed program for crawling HTML in specified url format. You can use the following parameters in command:
	-h			Get help
	-v			Show version of program
	-c			Specify configuration file input path, or default to ` + defaultConfigInputPath + `
	-l			Specify log output dir, or default to ` + defaultLogOutputDir
)

const timeFormat = "2006-01-02 15:04:05"

var (
	looger   *zerolog.Logger
	rfLogger *zerolog.Logger
)

func main() {
	allTaskDone := make(chan struct{})
	closeLogFiles := make(chan struct{})
	cfgPath, logDir := processCMD()
	looger, rfLogger = initLoggers(logDir, closeLogFiles)
	conf := initConfig(cfgPath)
	uniMap := &spider.UniMap{}
	taskQue := initTaskQue(conf)
	httpClient := &http.Client{
		Timeout: time.Duration(conf.CrawlTimeout) * time.Second,
	}
	for i := 0; i < conf.ThreadCount; i++ {
		s := spider.NewSpider(conf, httpClient, taskQue, uniMap)
		go crawling(allTaskDone, s)
	}
	select {
	case <-allTaskDone:
		// this notices goroutine that it's time to close the log files.
		closeLogFiles <- struct{}{}
		fmt.Println("minispider Done.")
	}
}

// processCMD parses the parameters in CMD and processes it.
func processCMD() (string, string) {
	helpFlag := flag.Bool("h", false, "help")
	versionFlag := flag.Bool("v", false, "version")
	configInputPath := flag.String("c", defaultConfigInputPath, "config input path")
	logOutputDir := flag.String("l", defaultLogOutputDir, "log output dir")

	flag.Parse()

	if *helpFlag {
		fmt.Println(helpText)
		os.Exit(0)
	}
	if *versionFlag {
		fmt.Println("minispider", version, "implemented by mendge")
		os.Exit(0)
	}
	return *configInputPath, *logOutputDir
}

// initConfig inits a Config based on file on cfgPath.
func initConfig(cfgPath string) *config.Config {
	conf, err := config.NewConfigFromFile(cfgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	err = config.ValidateConfig(conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return conf
}

// initTaskQue inits a task queue based on configure information.
func initTaskQue(cfg *config.Config) *spider.TaskQueue {
	sdq := spider.NewTaskDeque(cfg.ThreadCount)
	seeds, err := util.ParseSeeds(cfg.UrlListFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	for _, urlStr := range seeds {
		sdq.PushTask(spider.NewTask(0, urlStr))
	}
	return sdq
}

// initLoggers inits two logger redirected to log files in logDer.
// if logDir is not existed, make logDir build.
func initLoggers(logDir string, closeLogFiles chan struct{}) (*zerolog.Logger, *zerolog.Logger) {
	err := util.BuildDir(logDir)
	if err != nil {
		fmt.Println("can not build a new dir as root dir of logs")
		os.Exit(-1)
	}

	// Info level logger
	logFile, err := os.OpenFile(logDir+"logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("can not build log.txt in specified log output dir")
		os.Exit(-1)
	}
	logger := zerolog.New(logFile).With().Logger()

	// Warn level logger
	rfLogFile, err := os.OpenFile(logDir+"logs.wf.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("can not build log.rf.txt in specified log output dir")
		os.Exit(-1)
	}

	rfLogger := zerolog.New(rfLogFile).With().Logger()

	// 另起一个协程用于在程序退出时关闭日志文件
	go func(file *os.File, rfFile *os.File) {
		select {
		case <-closeLogFiles:
			file.Close()
			rfFile.Close()
		}
	}(logFile, rfLogFile)

	return &logger, &rfLogger
}

// crawling makes a spider container keep working until receives allTaskDone signal.
func crawling(allTaskDone chan struct{}, s *spider.Spider) {
	for {
		// gets task
		var scheduler spider.Scheduler = s
		task, noMoreTask := scheduler.FetchTask(s.TaskQ)
		if noMoreTask {
			allTaskDone <- struct{}{}
			return
		}
		getTaskTime := time.Now().Format(timeFormat)

		// downloads page
		var downloader spider.Downloader = s
		pageBody, err := downloader.DownloadHTML(s.HttpClient, task.DstURL)
		if err != nil {
			rfLogger.Warn().Str("URL", task.DstURL).Str("time", time.Now().Format(timeFormat)).Msg(fmt.Sprintf("fail to download page: %s", err))
			continue
		}
		downloadPageTime := time.Now().Format(timeFormat)

		// parses subURLs and add new tasks
		var processor spider.Processor = s
		links, err := processor.ParseSubURLs(s.Cfg.TargetUrl, pageBody, s.Cfg.TargetUrlRE)
		if err != nil {
			rfLogger.Warn().Str("URL", task.DstURL).Str("time", time.Now().Format(timeFormat)).Msg(fmt.Sprintf("fail to parse subURLs in page: %s", err))
		}
		newTaskCount := 0
		for _, urlStr := range links {
			visited := s.UniMap.IsKeyVisitedOrVisit(urlStr)
			if !visited && task.NowDepth < s.Cfg.MaxDepth {
				task := spider.NewTask(task.NowDepth+1, urlStr)
				scheduler.AddTask(task, s.TaskQ)
				newTaskCount++
			}
		}
		parseSubURLsTime := time.Now().Format(timeFormat)

		// saves page content
		var persistencer spider.Persistencer = s
		err = persistencer.SaveToFile(task.DstURL, s.Cfg.OutputDirectory, pageBody)
		if err != nil {
			rfLogger.Warn().Str("URL", task.DstURL).Str("time", time.Now().Format(timeFormat)).Msg(fmt.Sprintf("fail to save page to file: %s", err))
			continue
		}
		savePageTime := time.Now().Format(timeFormat)

		// if task done normally, record a log.
		looger.Info().Str("URL", task.DstURL).Str("get task", getTaskTime).
			Str("download Page", downloadPageTime).Str("parse sub URLs", parseSubURLsTime).
			Str("save to file", savePageTime).Str("new unprocessed task count", strconv.Itoa(newTaskCount)).Send()
		// sets a interval
		time.Sleep(time.Duration(s.Cfg.CrawlInterval) * time.Second)
	}
}
