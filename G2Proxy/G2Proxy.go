package G2Proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"reflect"
)

type ResData struct {
	RetCode int         //输出代码
	Message string      //输出消息
	OutData interface{} //输出结果
	Succeed int         //运行状态
	OutType string      //输出类型

	RunSafe int     //运行权限
	RunFlag int     //登录状态
	RunTime float64 //运行耗时
	SrvTime int64   //服务时间
	TokenID string  //运行身份
	state   int     //发送状态
}

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

	grpcClient := NewG2GrpcClient(client)

	// 调用远程方法
	cs, err := grpcClient.Service(context.Background(), args)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oPars := &ResData{}

	for {
		mRes, e4 := cs.Recv()

		if e4 == io.EOF {
			break
		}

		if e4 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": e4.Error()})
			return
		}
		//..........................................................................

		// 输出转换
		if e5 := MergeStructs(mRes, oPars); e5 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": e5.Error()})
			return
		}
		//..........................................................................

	}
	c.JSON(200, oPars)
}

type RpcClient struct {
	protocol string
	address  string
}

func (c *RpcClient) NewClient() (*grpc.ClientConn, error) {
	// 连接到 GRPC 服务
	client, err := grpc.Dial("localhost:9890", grpc.WithInsecure())

	if err != nil {
		log.Println("Dialing:", err)
		return nil, err
	}
	return client, nil
}

func MergeStructs(iFr, iTo interface{}) error {

	// 对象转换
	Change := func() error {
		v, e1 := json.Marshal(iFr) //
		if e1 != nil {
			return e1
		}
		//..............................................................................

		e := json.Unmarshal(v, iTo)
		if e != nil {
			return e1
		}
		//..............................................................................
		return nil
	}
	//..................................................................................

	// 输出转换
	mFr := reflect.ValueOf(iFr).Elem() //实例
	f, b := mFr.Type().FieldByName("OutData")

	if b && f.Type.String() == "string" {
		if e1 := Change(); e1 != nil {
			return e1
		}

		v := make([]interface{}, 0)
		Out := iFr.(*G2ResData).OutData
		e2 := json.Unmarshal([]byte(Out), &v)
		if e2 != nil {
			return e2
		}
		iTo.(*ResData).OutData = v[0]
		return nil
	}
	//..................................................................................

	// 输出转换
	if b {
		Out := []interface{}{iFr.(*G2ResData).OutData}
		v, e1 := json.Marshal(Out)
		if e1 != nil {
			return e1
		}
		iFr.(*G2ResData).OutData = string(v)
		if e2 := Change(); e2 != nil {
			return e2
		} else {
			return nil
		}
	}
	//..................................................................................
	return Change()
}
