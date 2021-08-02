## Requirement: 
基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

## Hints given by TA
> 本周的作业要求实现一个http server， http server是什么大家应该都知道的， 如果对此没有概念可以去学一下go语言基础小课 ，或看一下这篇官方教程https://golang.google.cn/doc/tutorial/web-service-gin； 作业中还提到国两点要求需要注意， 1使用errgroup管理协程； 2.能正确处理signal； errgroup是本周课毛老师重点讲过的；signal是编程中的一个重要概念， 以前不了解signal的可以去先自行学习一下相关知识，讲linux(UNIX)编程在书都会有专门章节来讲这个

> signal(信号)是进程间通信的一种方式， 它的成本很低，你只需要知道进程的ID就可以向该进程发信号；进程需要处理收到的信号， 如果你的程序没有设置信号处理，那就会使用默认方式去处理； 比如我们在console上执行一条命令， 这条命令执行到一半我们想终止它， 这个时候你按下ctrl+c， 命令可能就终止了，这是为什么呢？ 当你的终端上按下ctrl+c时，控制程序会向正在执行的命令发送一个SIGINT信号， 而SIGINT信号的默认处理程序就是直接退出执行；想象一下如果你的http server 正在处理客户端发来的请求， 这时收到了SIGINT信号， http server 直接退出会发生什么？   客户端发来的请求没有收到相应就连接中断了，这是很不“优雅“的；应该怎样做呢？  收到SIGINT信号后， 程序稍微多等一会， 等正在处理的请求都处理完成， 再退出，就比较"优雅"了

## References:
0. [gin-gonic/examples] graceful-shutdown/notify-with-context
1. [Tutorial: Developing a RESTful API with Go and Gin](https://golang.google.cn/doc/tutorial/web-service-gin)
2. [Golang之信号处理(signal)](https://zhuanlan.zhihu.com/p/128953024.vs)
3. [Golang开发笔记 9.10 Go Signal信号处理](https://www.bookstack.cn/read/golang_development_notes/zh-9.10.md)
4. [Graceful Restart in Golang](https://grisha.org/blog/2014/06/03/graceful-restart-in-golang/) In-progress requests completion/timeout