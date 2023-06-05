package znet

import (
	"github.com/linyerun/zinx/ziface"
)

// BaseHandler 实现IHandler时，先嵌入这个BaseHandler基类，然后根据需要堆这个基类的方法进行重写就🆗了
type BaseHandler struct {
}

//我们继承这个BaseHandler之后，可以只重写自己需要的方法就可以了

func (r *BaseHandler) PreHandle(_ ziface.IRequest) {

}

func (r *BaseHandler) Handle(_ ziface.IRequest) {

}

func (r *BaseHandler) PostHandle(_ ziface.IRequest) {

}
