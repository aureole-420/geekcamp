## Assignment
按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

## Solution
Kratos直接提供了解决方案
* project layout将project分为biz, data, server, service层。
* 使用proto buff定义API，并提供注册方法
* 使用wire启动kratos application 
* Krater.New()进行服务启动时，可选参数包括信号处理。

## References
1. Kratos
    - [Kratos 项目结构](https://go-kratos.dev/docs/intro/layout/)
    - [Kratos 注册接口](https://go-kratos.dev/docs/component/api)
    - [Kratos 依赖注入](https://go-kratos.dev/docs/guide/wire)
    - [Kratos crud示例 Blog](https://github.com/go-kratos/kratos/tree/main/examples/blog)  比hello-world有更多细节，也足够简单

1. Wire 依赖注入 
    - [wire core concepts](https://github.com/google/wire/blob/main/docs/guide.md) providers and injectors
    - [GoLang的wire框架使用例子](https://www.cnblogs.com/llh4cnblogs/p/13636195.html)
    - [undefined: InitializeEvent](https://github.com/google/wire/issues/224) wire_gen.go is also part of the main package; either `go build` or `go build main.go wire_gen.go`