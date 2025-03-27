package main

import (
	"awesomeProject/G2Proxy"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/*path", func(c *gin.Context) {
		G2Proxy.HandleToRPC(c)
	})
	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080
}
