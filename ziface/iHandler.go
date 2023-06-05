package ziface

// IHandler 开发者自定义请求处理器的抽象层
type IHandler interface {
	PreHandle(request IRequest)  //在处理conn业务之前的钩子
	Handle(request IRequest)     //在处理conn业务的主方法hook
	PostHandle(request IRequest) //在处理conn业务之后的钩子
}
