# Cera 
#### Cera是一个极简的http小框架，在阅读fasthttp的时候顺手编写的一个http框架雏形, 路由在fashttprouter基础上修改的，该框架建议作为学习娱乐使用，暂不建议用于生产环境，除非您熟悉cera源码

支持功能如下：
- [x] 支持基本的mVC架构，没有对orm进行整合，个人觉得没有这个必要，如果需要使用orm框架请自行引用gorm等框架。
- [x] 支持基础http路由
- [x] 支持头部中间件和尾部处理（中间件）
- [x] View支持api view（json返回值）和web view
- [x] Render支持golang template渲染，通过渲染输出页面
- [x] 支持jwt基础功能，但暂未将token回调解析出的信息放入user结构中（暂未想到合理的安放方式）
- [ ] xrsf token计划支持
- [ ] session 计划支持
- [ ] 计划支持jinjia2模板

#### 最简单的使用方式
##### 1. 创建基础的目录结构
您的应该创建如下目录结构：
```shell
---project #新建一个项目目录
  |---- template #您的模板目录，此处存放您的html模板，文件名应以.htm结尾
  |---- static   #您的静态文件目录，此处存放所需要的静态文件资源
  |---- views    #您的视图函数目录，您的所有视图函数可以放在此处。当然，您也可以放在项目根目录下
  |---- routes  #存放您的http路由文件
      |---- routes.go #存放您的http路由
  |---- {{project}}.go #项目起始文件,您也可以单独创建cmd目录来存放main方法，但要注意，所有的静态文件、模板目录路径目前使用的都是./xxxx，如有需要，可以自行修改
```
cera的http server使用fasthttp封装，采用复用式请求处理方法，性能方面应该无需担心
##### 2. 创建一个视图
```go
// project/views/loginView.go

type Login struct {
	view.ApiView  // 如果是只处理返回json值，此处如要引用ApiView，如果是需要渲染页面，此处则需要引用view.View
}

func (l *Login) Get() {  // 处理Get方法
	l.Data["a"] = "test" // 返回值，注意Data中所有数据都将被渲染输出
}
func (l *Login) Post() { // 处理Post方法
  
	l.Data["a"] = "test"
}
// 如果路由中包含此处未处理的方法，默认返回404 not found

// project/views/paasView.go

type Paas struct {
	view.ApiView
}

func (p *Paas) Get() {
	p.Ctx.Response.Header.SetStatusCode(200)
	p.Data["hello"] = "world"
}

```
##### 3. 创建一个路由并引用视图
```go
// project/routes/routes.go
func Router() *router.Router {
    r := router.New()
    r.GET("/auth/login", &views.Login{}) // 创建一个路径为/auth/login的路由，接受get方法，并指定由login view处理
    r.POST("/auth/login", &views.Login{})
    r.ANY("/", &views.Paas{})  // 接受所有方法
    return r
}

```

##### 4. 创建一个http server

```go
//{{project}}.go
func main(){
    h := http.NewHttpServe("127.0.0.1", "9999")
    logger := logrus.New() //传入的logger需要实现cera.SimpleLogger接口
    h.SetLogger(logger) // 为http server设置logger
    h.SetRouter(routes.Router())   // 设置路由
    h.Start() // 启动服务
}
```

