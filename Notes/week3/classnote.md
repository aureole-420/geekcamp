# 1. goroutine

[main.go](./main.go)

### References
* [浅谈goroutine](https://www.jianshu.com/p/7ebf732b6e1f)
* [go channel, buffered与unbuffered的区别](https://www.dazhuanlan.com/maplelalala/topics/1026868)
* [go routine内存泄漏](https://zhuanlan.zhihu.com/p/352589023)


# 2. memory model


# 3. Package sync

### (1) Share Memory by communicating: 并发调用resource poller的例子，sync.mutex实现 vs channel实现

### (2) Race conditions: 
*   `go build -race` 检测; 
*   底层数据的赋值mutex或者atomic更合适，channel也可以做但太重了.;
*   interface data race的例子:IcecreamMaker, two implementation, Ben and Jerry。不要对golang数据结构做任何假设，比如slice，interface都不是data race safe(non-atomic), map是data race safe(因为返回的是map指针)但仍然有data race. 建议用atomic value

### (3) sync.atomic 
* example: Config struct, 四个人读，一个人写
* 同步语义
    - 互斥锁 mutex  
    - 读写锁 rwmutex； 相对更重，涉及到更多的goroutine 上下文切换
        ```golang
        var l sync.RWMutex

        l.Lock()
        write
        l.Unlock()


        l.RLock()
        read
        l.RUnlock()
        ```
    - atomic value; 如果读特别多，那么atomic value性能极好
        ```golang
        var v atomic.Value

        // write 
        cfg := &Config{...} // 每次构建新对象，copy on write
        v.Store(cfg)


        // read
        cfg := v.Load().(*Config)  // 需要assert
        ```
### (3.1) Copy-On-Write 
* 思路在microservice degredation 或者local cache中普遍使用。copy on write + atomic values 既满足 原子性，也满足 可见性（数据更新后，内存地址也变化，别的goroutine也拿到的是相同的；如果复用老数据，内存地址不变，某些goroutine可能从store cache里面拿老数据）
* example： 共享map： 用互斥锁mutex保证map的写，更新时放入 atomic value中，读操作只用load atomic value即可。

### (4) Mutex 
* 获取锁的goroutine会被加入FIFO waiting queue中,锁被释放时会唤醒队首goroutine。唤醒需要时间
* 某些情况下会有starvation，某goroutine不停的拿锁，等待的goroutine被唤醒后发现还是拿不到锁。
    *  Bargin （无脑让出锁） vs Handoff （降低吞吐，释放锁的goroutine hands off lock to goroutine in FIFO waiting queue.） vs Spinning (空waiting queue时多尝试一下).
    * Go 1.9后，starvation mode，保证handoff会在出现starvation(>1ms有界等待)kick in

### (5) errgroup
* 利用 sync.Waitgroup 管理并执行goroutine
    * 并行工作流
    * 错误处理，优雅降级
    * context传播和取消
    * 局部变量+闭包(closure)
    ```golang
    g, ctx := errgroup.WithContext(context.Background())

    g.Go(func() error {
        // dosomething 1
    })

    g.Go(func() error {
        // dosomething 2
    })

    g.Go(func() error {
        // dosomething 3
    })

    err := g.Wait()
    fmt.Println(err)
    fmt.Println(ctx.Err())
    ```
* source code: [golang/sync/errgroup/errgroup.go](https://github.com/golang/sync/blob/master/errgroup/errgroup.go) 总共就几十行代码.
    - sync.Once 单例 只能被执行一次

### (6) sync.Pool
保存和复用临时对象(in stack), 减少GC压力.
不能放连接池或带装状态的东西，只能放不确定什么时候会被回收的东西.
全局对象,
> When there is an expensive object you have to create it frequently, it can be very beneficial to use sync.Pool

> Pools often contain things like *bytes.Buffer, which are temporary and re-usable.


### References
* [Golang并发：channel vs sync.Mutex](https://segmentfault.com/a/1190000017890174)
    > channel的核心是数据流动，关注到并发问题中的数据流动，把流动的数据放到channel中，就能使用channel解决这个并发问题。

    > mutex的能力是数据不动，某段时间只给一个协程访问数据的权限擅长数据位置固定的场景
* [Go: mutex and starvation](https://medium.com/a-journey-with-go/go-mutex-and-starvation-3f4f4e75ad50) 讲解了golang mutex中Bargin vs Handoff vs Spinning.

# 4. Package context
### (1) Channel
* 常见最佳实践: 任务分发
* unbuffered channel
    * 任何叫交换前需要两端goroutine同时准备好。 -- 同步通信
    * 本质是保证同步;接受端尚未接受时，发送端goroutine会被block;发送端伤未发送消息前，接受端会被阻塞；
* buffered channel
    * 只有当buffer满时，发送端才会block;只有当buffer未空时，接受端才会block
    * 缓冲区大小会影响性能，buffer太小会频繁的出现阻塞; buffer远大于consumer数量时性能影响就不大了。

### (2) Go concurrency patterns
> Doc [Go concurrency patterns](https://blog.golang.org/concurrency-timeouts)
* Timing out: 用channel发送超时消息. In practice use `time.After`返回的buffered channel即可.
    ```golang
    timeout := make(channel bool, 1) // size=1, timeout goroutine can  send and then exit.
    go func() {
        time.Sleep(1 * time.Second)
        timeout <- true
    }

    select {
        case <-ch:
            // read from a ch occurred
        case <- timeout: // 直接替换成time.After()
            // the read from ch has timed out.
    }
    ```
* Moving on: [nonblocking channel operations](https://gobyexample.com/non-blocking-channel-operations), 用select default语句，如果channel operation block住了，那就自动跳到default. 
    ```golang
    func Query(conns []Conn, query string) Result {
        ch := make(chan Result) //这里是个bug, 如果send先于receive，那么所有的send都会被block，全部goroutine都fall to default；解决方法： 换成 make(chan Result, 1)至少让一个goroutine能够把send message.
        for _, conn := range conns {
            go func(c Conn) {
                select {
                case ch <- c.DoQuery(query): 
                default:
                }
            }(conn)
        }
        return <-ch // receive
    }
    ```
* Pipeline
* Fanout, Fanin
* Cancellation:用channel表示cancellation
    - close限于receive发生；一定要在没有写操作后再关闭通道。
* Context


### (3) Design philosophy
* buffered channel: 考虑满了后的行为，需要等待还是丢弃 
* 实际直接用channel比较少，一般用errgroup比较多

### (4) Package context:
> Doc: [Go concurrency patterns: Context](https://blog.golang.org/context)
 
go 1.7    引入context解决了concurrency中最典型两大问题：超时 和 取消

#### Request-scoped context
Fail fast: 下游goroutine出现高耗时的情况时，只有快速地让goroutine失败回收资源才能保证系统稳定.
```golang
// A Context carries a deadline, cancelation signal, and request-scoped values
// across API boundaries. Its methods are safe for simultaneous use by multiple
// goroutines.
type Context interface {
    // Done returns a channel that is closed when this Context is canceled
    // or times out.
    Done() <-chan struct{} // 只读channel

    // Err indicates why this context was canceled, after the Done channel
    // is closed.
    Err() error

    // Deadline returns the time when this Context will be canceled, if any.
    Deadline() (deadline time.Time, ok bool)

    // Value returns the value associated with key or nil if none.
    Value(key interface{}) interface{} // 挂载元数据
}conte
```
* context 是goroutine safe,不要修改context里面的元数据;use withValue来copy on write.
* `context.WithCancel`, `context.WithDeadline` 获取cancel方法从而获得cancel权利 `defer cancel()`
* 简单讲解go context的文章[在Go中用Context取消操作](https://zhuanlan.zhihu.com/p/163061156)
    - 监听取消事件: Select case <- context.Done()
    - 触发取消事件: cancel()
    - 过期: context.WithTimeout(); client端直接把ctx给req，httpClient.Do(req) 会自动超时过期。