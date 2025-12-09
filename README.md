# 这是木犀审核的后端仓库

## 如何使用？

你可以通过MakeFIle的相关指令便捷操作：

```bash
# 直接构建应用 -- 您需要预先完善好config.yaml和.env文件
make build 
# 这会在dist目录生成可执行文件 （windows和linux支持），请保证env文件和可执行文件在同级目录

#如果您是mac用户，或者有云部署的需求
make docker
make deploy
```



项目架构：

```go
main :---APP -{
			config
			server // 核心，所有相关服务的入口，内置http server
				}

server : {
    router // 具体的http.Server实例，如gin.router
    close()  // 服务关闭函数，负责调度需要在shutdown前处理的事情
}

router : {
    rgs : //具体服务的路由
    {
        middlewares  //  统一处理请求前，请求中，请求后需要执行的逻辑的中间件
    	controllers // 具体服务的控制器，根据不同路径处理请求体并执行相关核心逻辑
    } 
}

controller {
    service  // 服务的核心逻辑
    api_code  // 统一返回的具体的响应码
    api_error // 若有错误，统一封装并返回的错误
}

service {
    middlewares // 需要使用的各类中间件，如数据库，日志，消息队列
}


// 最后的所有都通过wire依赖注入串联
```

不足之处：

```go
sdk (方便外部服务使用的包) 与 langchain 部分都根据领域原则独立成包，但审核的具体服务如user、project等都不独立，model、service等写在了一起，亟待优化。
```