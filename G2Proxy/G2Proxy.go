package G2Proxy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func HandleToRPC(c *gin.Context) {

	clientConfig := &RpcClient{
		protocol: "tcp",
		address:  "localhost:9890",
	}
	client, err := clientConfig.NewClient()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer client.Close()

	// 获取请求路径
	path := c.Param("path")
	fmt.Println(path)
	// 解析表单数据
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}
	// 将表单数据转换为 map[string]string
	reqData := make(map[string]string)
	err = c.BindJSON(&reqData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}
	// 准备请求参数
	args := &G2ReqData{
		Id:     "",
		Addr:   "",
		Host:   "",
		Port:   "",
		Envs:   "",
		Path:   path,
		Uris:   "",
		Form:   reqData,
		Body:   nil,
		Source: "mgrs",
		Region: "localhost",
		Agents: "isme",
		Secure: 0,
		UsrID:  "",
		Roles:  []string{"Admin", "Guest"},
		Safes:  1,
		Login:  1,
	}
	var result = &G2ResData{}

	// 调用远程方法
	err = client.Call("Provide.Service", args, result)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, err)
	}
	c.JSON(200, result)
}

type RpcClient struct {
	protocol string
	address  string
}

func (c *RpcClient) NewClient() (*rpc.Client, error) {
	// 连接到 RPC 服务
	client, err := jsonrpc.Dial("tcp", "localhost:9890")
	if err != nil {
		log.Println("Dialing:", err)
		return nil, err
	}
	return client, nil
}
