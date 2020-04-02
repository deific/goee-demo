package main

import (
	"github.com/deific/goee"
	"github.com/deific/goee/core"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// 启动入口
func main() {
	// 实例化
	g := goee.Default()

	// 加载配置文件
	g.LoadConfig("conf/goee-dev.yaml")

	// 注册静态文件处理
	g.Static("/static", "E:/temp/static")
	// 加载模板
	g.LoadHtmlTemplate("templates/*")

	// 注册路由分组
	v1 := g.Group("/v1")
	v1.GET("/hello", func(c *core.Context) {
		c.HTML(http.StatusOK, "<h1>hello "+c.GetParam("name")+"</h1>")
	}).GET("/hello/:name", func(c *core.Context) {
		c.HTML(http.StatusOK, "<h1>hello "+c.Param("name")+"</h1>")
	})

	// 注册路由，想到根分组
	g.GET("/", func(c *core.Context) {
		c.HTML(http.StatusOK, "<h1>Welcome to Goee!</h1> you request path = "+c.Path)
	})
	g.GET("/hello/:name", func(c *core.Context) {
		log.Println("执行实际处理函数")
		c.HTML(http.StatusOK, "<h1>hello "+c.Param("name")+"</h1>")
	})
	g.GET("/json", func(c *core.Context) {
		c.JSON(http.StatusOK, core.HMap{"code": 200, "msg": "成功", "success": true})
	})
	// 使用html模板
	g.GET("/html", func(c *core.Context) {
		c.HtmlTemplate(http.StatusOK, "index.tpl", core.HMap{"title": "张三", "now": time.Now()})
	})
	// 代理
	g.GET("/proxy", func(c *core.Context) {
		trueServer := "http://www.baidu.com/"
		url, err := url.Parse(trueServer)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("proxy to:", url)
		proxy := httputil.NewSingleHostReverseProxy(url)
		// Update the headers to allow for SSL redirection
		c.Req.URL.Host = url.Host
		c.Req.URL.Scheme = url.Scheme
		c.Req.Header.Set("X-Forwarded-Host", c.Req.Header.Get("Host"))
		c.Req.Host = url.Host
		c.Req.URL.Path = url.Path
		proxy.ServeHTTP(c.Writer, c.Req)
	})

	// 异常路由
	g.GET("/panic", func(c *core.Context) {
		name := []string{"goee"}
		c.JSON(http.StatusOK, name[10])
	})

	// 启动服务
	//g.Run(":9000")
	g.Start()
}
