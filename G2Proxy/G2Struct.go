package G2Proxy

type G2ReqData struct {
	Id   string
	Addr string
	Host string
	Port string
	Envs string
	Path string
	Uris string

	Form map[string]string //请求参数
	//File []G2IFile         //上传文件
	Body []byte //其它数据

	Source string //访问来源
	Region string //数据区域
	Agents string //应用代码
	Secure int    //检查项目

	UsrID string   //用户代码
	Roles []string //隶属角色
	Safes int      //访问权限
	Login int      //是否登录
}

type G2ResData struct {
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
}
