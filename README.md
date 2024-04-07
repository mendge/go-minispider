
## [题目]
使用go开发一个迷你定向抓取器 mini_spider，对一些网站进行定向抓取。实现对种子链接的抓取，并把URL长相符合特定pattern的网页保存到磁盘上。
### [程序运行]
./mini_spider -c ../conf -l ../log
### [配置文件]
```azure
[spider]
# 种子文件路径
urlListFile = ../url.data
# 抓取结果存储目录
outputDirectory = ../output
# 最大抓取深度(种子为0级)
maxDepth = 3
# 抓取间隔. 单位: 秒
crawlInterval =  1
# 抓取超时. 单位: 秒
crawlTimeout = 3
# 需要存储的目标网页URL pattern(正则表达式)
targetUrl = .*.(htm|html)$
# 抓取routine数
threadCount = 22
```
### [种子文件]
json 格式
```["http://www.baidu.com","http://www.sina.com.cn",...]```
### [要求]
1. 支持命令行参数处理。具体包含:
   *  -h(帮助)
   * -v(版本)
   * -c(配置文件路径）
   * -l(配置日志文件路径)
2. 抓取网页的顺序没有限制
3. 单个网页抓取或解析失败，不导致整个程序退出。在日志中记录下错误原因并继续
4. 当程序完成所有抓取任务后，程序优雅退出
5. 从HTML提取链接时处理相对路径和绝对路径
6. 需要能够处理不同字符编码的网页(例如 utf-8 或 gbk )
7. 网页存储时每个网页单独存为一个文件，以URL为文件名对，URL中的特殊字符做转义
8. 支持多routine并行抓取（并不是指简单设置GOMAXPROCS>1)
9. 编写单元测试，单元测试有效而且通过
10. 实现“去重”抓取功能，对于已经抓取过的网页不要重复抓取
11. 对于配置文件的进行异常检查和处理
12. 保存程序两个运行日志：mini_spider.log(信息日志) 和 mini_spider.wf.log(告警日志)
 
## [安装运行]
```azure
$ git clone https://github.com/mendge/go-minispider
$ cd go-minispider
$ go build -o ./bin/mini_spider main.go
...按需要修改配置文件
$ cd bin
$ mini_spider -c ../conf -l ../log
```
