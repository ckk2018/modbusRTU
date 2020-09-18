package request

import (
	"ckklearn.com/testmodbus/global"
)

// 所有读写请求的共用信息
// 仅用于嵌套，方便序列化
type base struct {
	addr   byte           // 从站号
	fun    global.FunCode // 功能码
	offset uint16         // 偏移量
}

// FunCode 请求的功能码
func (b base) FunCode() global.FunCode {
	return b.fun
}
