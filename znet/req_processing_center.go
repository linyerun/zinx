package znet

import (
	"fmt"
	. "github.com/linyerun/zinx/global_properties"
	"github.com/linyerun/zinx/ziface"
	"math/rand"
	"time"
)

// ReqProcessingCenter 消息处理模块的实现
type ReqProcessingCenter struct {
	HandlersMap  map[uint32]ziface.IHandler //msgId:IHandler
	WorkPoolSize uint32                     //工作池大小
	Workers      []chan ziface.IRequest     //每个任务协程都有一个与之对应的待处理的任务channel
}

func newReqProcessingCenter() ziface.IReqProcessingCenter {
	return &ReqProcessingCenter{
		HandlersMap:  make(map[uint32]ziface.IHandler),
		WorkPoolSize: GlobalObject.WorkPoolSize,
		Workers:      make([]chan ziface.IRequest, GlobalObject.MaxWorkerTask),
	}
}

func (m *ReqProcessingCenter) GetWorkPoolSize() uint32 {
	return m.WorkPoolSize
}

func (m *ReqProcessingCenter) AddHandler(msgId uint32, handler ziface.IHandler) {
	_, ok := m.HandlersMap[msgId]
	if ok {
		panic("msgId has existed in HandlersMap")
	}
	m.HandlersMap[msgId] = handler
	fmt.Printf("add handler(msgId=%d) successfully\n", msgId)
}

func (m *ReqProcessingCenter) ExecReqAppointedHandler(req ziface.IRequest) {
	router, ok := m.HandlersMap[req.GetMsgId()]
	if !ok {
		fmt.Println("Api msgId =", req.GetMsgId(), "is Not found!", "you need Register")
		return
	}
	//模板方法模式
	router.PreHandle(req)
	router.Handle(req)
	router.PostHandle(req)
}

func (m *ReqProcessingCenter) SendReqToWorkers(req ziface.IRequest) {
	//随机分配给一个任务协程处理
	rand.Seed(time.Now().Unix())
	workerID := rand.Uint32() % m.WorkPoolSize
	fmt.Println("add ConnID =", req.GetConnection().GetConnID(), "to WorkerID =", workerID, "to execute")
	//发送到对于的worker去处理
	m.Workers[workerID] <- req
}

// StartWorkers 开启协程池
func (m *ReqProcessingCenter) StartWorkers() {
	for i := uint32(0); i < m.WorkPoolSize; i++ {
		//开启一个工作协程
		m.Workers[i] = make(chan ziface.IRequest, GlobalObject.MaxWorkerTask)
		go m.startOneWorker(i, m.Workers[i])
	}
}

// 开启一个协程
func (m *ReqProcessingCenter) startOneWorker(workID uint32, taskQueue chan ziface.IRequest) {
	fmt.Println("workID =", workID, "start successfully....")
	for {
		//取到一个之后，放到DoMsgHandle里面去执行
		//取不到就阻塞
		m.ExecReqAppointedHandler(<-taskQueue)
	}
}
