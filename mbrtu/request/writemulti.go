package request

import (
	"bytes"
	"encoding/binary"

	"ckklearn.com/testmodbus/global"
	"ckklearn.com/testmodbus/mbrtu/mbcrc"
)

// 所有写多个请求的共用信息
// 仅用于嵌套，方便序列化
type writeMultiBase struct {
	base
	num      uint16
	dataSize byte
}

// ---- 写线圈 ----

// RtuWriteMultiCoilsRequest 写多个线圈请求
// 写入连续的线圈
type RtuWriteMultiCoilsRequest struct {
	writeMultiBase
	data []byte
}

// NewRtuWriteMultiCoilsRequest 构造函数
func NewRtuWriteMultiCoilsRequest(addr byte, offset, num uint16, data []byte) *RtuWriteMultiCoilsRequest {
	return &RtuWriteMultiCoilsRequest{
		writeMultiBase: writeMultiBase{
			base: base{
				addr:   addr,
				fun:    global.WriteMultiCoils,
				offset: offset,
			},
			num:      num,
			dataSize: byte(len(data)),
		},
		data: data,
	}
}

// Serialize 将结构序列化为rtu请求报文
func (r *RtuWriteMultiCoilsRequest) Serialize(buf *bytes.Buffer, crcOrder binary.ByteOrder) error {
	err := binary.Write(buf, binary.BigEndian, r.writeMultiBase)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.BigEndian, r.data)
	if err != nil {
		return err
	}
	crc16 := mbcrc.Crc16(buf.Bytes())
	return binary.Write(buf, crcOrder, crc16)
}

// FunCode 功能码
func (r *RtuWriteMultiCoilsRequest) FunCode() global.FunCode {
	return r.fun
}

// ExpectedLen 期望的返回报文字节长度
func (r *RtuWriteMultiCoilsRequest) ExpectedLen() int {
	return 8
}

// ---- 写寄存器 ----

// RtuWriteMultiRegsRequest 写多个保持寄存器请求
// 写入连续的保持寄存器
type RtuWriteMultiRegsRequest struct {
	writeMultiBase
	data []uint16
}

// NewRtuWriteMultiRegsRequest 构造函数
func NewRtuWriteMultiRegsRequest(addr byte, offset uint16, data []uint16) *RtuWriteMultiRegsRequest {
	regNum := uint16(len(data))
	return &RtuWriteMultiRegsRequest{
		writeMultiBase: writeMultiBase{
			base: base{
				addr:   addr,
				fun:    global.WriteMultiRegisters,
				offset: offset,
			},
			num:      regNum,
			dataSize: byte(regNum * 2),
		},
		data: data,
	}
}

// Serialize 将结构序列化为rtu请求报文
func (r *RtuWriteMultiRegsRequest) Serialize(buf *bytes.Buffer, crcOrder binary.ByteOrder) error {
	err := binary.Write(buf, binary.BigEndian, r.writeMultiBase)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.BigEndian, r.data)
	if err != nil {
		return err
	}
	crc16 := mbcrc.Crc16(buf.Bytes())
	return binary.Write(buf, crcOrder, crc16)
}

// FunCode 功能码
func (r *RtuWriteMultiRegsRequest) FunCode() global.FunCode {
	return r.fun
}

// ExpectedLen 期望的返回报文字节长度
func (r *RtuWriteMultiRegsRequest) ExpectedLen() int {
	return 8
}
