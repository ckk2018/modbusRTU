package request

import (
	"bytes"
	"encoding/binary"

	"ckklearn.com/testmodbus/global"
	"ckklearn.com/testmodbus/mbrtu/mbcrc"
)

// RtuReadRequest 读数据请求
// 读取连续的寄存器或线圈
type RtuReadRequest struct {
	base
	num uint16
}

// NewRtuReadRequest 构造函数
func NewRtuReadRequest(addr byte, fun global.FunCode, offset, num uint16) *RtuReadRequest {
	return &RtuReadRequest{
		base: base{
			addr:   addr,
			fun:    fun,
			offset: offset,
		},
		num: num,
	}
}

// Serialize 将结构序列化为rtu请求报文
func (r *RtuReadRequest) Serialize(buf *bytes.Buffer, crcOrder binary.ByteOrder) error {
	err := binary.Write(buf, binary.BigEndian, r)
	if err != nil {
		return err
	}
	crc16 := mbcrc.Crc16(buf.Bytes())
	return binary.Write(buf, crcOrder, crc16)
}

// ExpectedLen 期望的返回报文字节长度
func (r *RtuReadRequest) ExpectedLen() int {
	var dataLen int

	switch r.fun {
	case global.ReadCoils, global.ReadInputs:
		dataLen = (int(r.num)-1)/8 + 1
	case global.ReadHoldingRegisters, global.ReadInputRegisters:
		dataLen = int(r.num) * 2
	}

	return dataLen + 5
}
