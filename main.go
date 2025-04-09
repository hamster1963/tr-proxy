package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tr-proxy/G2Proxy"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 如果是 OPTIONS 请求，直接返回状态码 204
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	// 使用跨域中间件
	r.Use(Cors())

	r.POST("/*path", func(c *gin.Context) {
		G2Proxy.HandleToRPC(c)
	})
	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080
}
