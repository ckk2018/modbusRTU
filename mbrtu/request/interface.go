package request

import (
	"bytes"
	"encoding/binary"

	"ckklearn.com/testmodbus/global"
)

// RtuRequest 所有rtu请求结构都要实现这个接口
type RtuRequest interface {
	// Serialize 将结构数据序列化为报文
	Serialize(buf *bytes.Buffer, crcOrder binary.ByteOrder) error

	// FunCode 请求的功能码
	FunCode() global.FunCode

	// ExpectedLen 当前请求所期望的返回报文字节长度
	// 期望的长度为正常返回的长度，异常返回的长度需要自行判断
	ExpectedLen() int
}
