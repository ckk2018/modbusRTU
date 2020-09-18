package request

import (
	"bytes"
	"encoding/binary"

	"ckklearn.com/testmodbus/global"
	"ckklearn.com/testmodbus/mbrtu/mbcrc"
)

// RtuWriteSingleRequest 写单个数据请求
// 写入单个的寄存器或线圈
type RtuWriteSingleRequest struct {
	base
	data uint16
}

// NewRtuWriteSingleRequest 构造函数
func NewRtuWriteSingleRequest(addr byte, fun global.FunCode, offset, data uint16) *RtuWriteSingleRequest {
	return &RtuWriteSingleRequest{
		base: base{
			addr:   addr,
			fun:    fun,
			offset: offset,
		},
		data: data,
	}
}

// Serialize 将结构序列化为rtu请求报文
func (r *RtuWriteSingleRequest) Serialize(buf *bytes.Buffer, crcOrder binary.ByteOrder) error {
	err := binary.Write(buf, binary.BigEndian, r)
	if err != nil {
		return err
	}
	crc16 := mbcrc.Crc16(buf.Bytes())
	return binary.Write(buf, crcOrder, crc16)
}

// ExpectedLen 期望的返回报文字节长度
func (r *RtuWriteSingleRequest) ExpectedLen() int {
	return 8
}
