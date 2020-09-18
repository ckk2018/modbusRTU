// 针对南瑞项目的变体rtu请求结构

package request

import (
	"bytes"
	"encoding/binary"

	"ckklearn.com/testmodbus/global"
	"ckklearn.com/testmodbus/mbrtu/mbcrc"
)

type nrWriteMultiBase struct {
	base
	num      uint16
	dataSize uint16 // 2个字节，区别于标准mbrtu
}

// NRWriteMultiRegsRequest 写多个保持寄存器
type NRWriteMultiRegsRequest struct {
	nrWriteMultiBase
	data []uint16
}

// NewNRWriteMultiRegsRequest 构造函数
func NewNRWriteMultiRegsRequest(addr byte, offset uint16, data []uint16) *NRWriteMultiRegsRequest {
	regNum := uint16(len(data))
	return &NRWriteMultiRegsRequest{
		nrWriteMultiBase: nrWriteMultiBase{
			base: base{
				addr:   addr,
				fun:    global.WriteMultiRegisters,
				offset: offset,
			},
			num:      regNum,
			dataSize: regNum * 2,
		},
		data: data,
	}
}

// Serialize 将结构序列化为请求报文
func (r *NRWriteMultiRegsRequest) Serialize(buf *bytes.Buffer, crcOrder binary.ByteOrder) error {
	err := binary.Write(buf, binary.BigEndian, r.nrWriteMultiBase)
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

// ExpectedLen 期望的返回报文字节长度
func (r *NRWriteMultiRegsRequest) ExpectedLen() int {
	return 8
}
