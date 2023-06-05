package global_properties

import (
	"encoding/json"
	"fmt"
	"github.com/linyerun/zinx/ziface"
	"io/fs"
	"io/ioutil"
)

type GlobalObj struct {
	/**
	Server
	*/
	TcpServer ziface.IServer //当前Zinx全局Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口
	Name      string         //当前服务器的名称
	/**
	Zinx
	*/
	Version        string //当前Zinx的版本号
	MaxConn        int    //当前服务器主机允许的最大连接数
	MaxPacketSize  uint32 //当前Zinx框架数据包的最大值
	WorkPoolSize   uint32 //工作池大小
	MaxWorkerTask  uint32 //任务队列长度
	MaxMsgChanSize uint32 //连接的Writer最大处理信息数
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("config/zinx.json")
	if _, ok := err.(*fs.PathError); ok {
		fmt.Println("您没用在`/config/zinx`配置全局变量,全局变量均使用默认值")
		return
	}
	if err != nil { //出现错误我就直接使用默认配置
		fmt.Println("your config file has error,it will use default config!")
		fmt.Println(err)
		return
	}
	if err = json.Unmarshal(data, GlobalObject); err != nil {
		fmt.Println("your config json may have error,it will use default config!")
		fmt.Println(err)
		return
	}
	if GlobalObject.MaxPacketSize < 1 {
		GlobalObject.MaxPacketSize = 1024
	}
}

// Init 初始化当前GlobalObject
func Init() {
	//首先，设置默认值
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V1.0", //这个是自定义的，不能被用户更改
		Host:           "0.0.0.0",
		TcpPort:        8999,
		MaxConn:        1000,
		MaxPacketSize:  1024 * 4,
		WorkPoolSize:   10,
		MaxWorkerTask:  1024,
		MaxMsgChanSize: 10,
	}
	//然后尝试从config/zinx.json中去加载一些用户自定义参数
	GlobalObject.Reload()
	GlobalObject.Version = "V1.0"
}
