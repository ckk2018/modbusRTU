package mbrtu

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/tarm/serial"

	"ckklearn.com/testmodbus/global"
	"ckklearn.com/testmodbus/mbrtu/mbcrc"
	"ckklearn.com/testmodbus/mbrtu/request"
)

const (
	// rtu下返回报文最小长度（异常返回报文）
	minResLen int = 5
)

// RtuMaster modbus主站结构
type RtuMaster struct {
	s           *serial.Port
	reqCrcOrder binary.ByteOrder // 最后一次请求的crc16校验码字节序
	reqFunCode  global.FunCode   // 最后一次请求的功能码
	reqExpLen   int              // 最后一次请求的期望返回报文字节长度
}

// NewRtuMaster 构造函数
// 串口在这里被初始化
func NewRtuMaster(c *serial.Config) (*RtuMaster, error) {
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	return &RtuMaster{s: s}, nil
}

// Close 关闭主站
func (m *RtuMaster) Close() error {
	return m.s.Close()
}

// Write 向串口写数据
func (m *RtuMaster) Write(r request.RtuRequest, crcOrder binary.ByteOrder) (int, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 8))
	err := r.Serialize(buf, crcOrder)
	if err != nil {
		return 0, err
	}
	n, err := m.s.Write(buf.Bytes())
	if err != nil {
		return 0, err
	}
	m.reqCrcOrder = crcOrder
	m.reqFunCode = r.FunCode()
	m.reqExpLen = r.ExpectedLen()
	return n, nil
}

// 将超时也作为异常抛出
func (m *RtuMaster) _read(p []byte) (int, error) {
	n, err := m.s.Read(p)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, fmt.Errorf("read timeout")
	}
	return n, nil
}

// Read 读取串口数据
func (m *RtuMaster) Read(p []byte) (int, error) {
	read := 0
	raw := make([]byte, 1024)

	// 先保证读到最小长度的报文
	for read < minResLen {
		n, err := m._read(raw[read:])
		if err != nil {
			return 0, err
		}
		read += n
	}

	// 如果不是异常，则继续读取剩余报文
	if raw[1] != m.reqFunCode+0x80 {
		for read < m.reqExpLen {
			n, err := m._read(raw[read:])
			if err != nil {
				return 0, err
			}
			read += n
		}
	}

	// 利用crc16校验接收包
	var readCrc uint16
	binary.Read(bytes.NewReader(raw[read-2:]), m.reqCrcOrder, &readCrc)
	calCrc := mbcrc.Crc16(raw[:read-2])
	if readCrc != calCrc {
		return 0, fmt.Errorf("validate failed: readcrc %x, calcrc %x", readCrc, calCrc)
	}

	// 解析从站返回数据
	n, err := RtuParseResponse(p, raw, m.reqFunCode)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// ---- 组合 `Read` 和 `Write` 的高级方法 ----
// ---- 标准mbrtu ----

// ReadCoils 读取线圈
func (m *RtuMaster) ReadCoils(p []byte, addr byte, offset, num uint16, crcOrder binary.ByteOrder) (int, error) {
	_, err := m.Write(request.NewRtuReadRequest(addr, global.ReadCoils, offset, num), crcOrder)
	if err != nil {
		return 0, err
	}
	return m.Read(p)
}

// ReadInputs 读取输出
func (m *RtuMaster) ReadInputs(p []byte, addr byte, offset, num uint16, crcOrder binary.ByteOrder) (int, error) {
	_, err := m.Write(request.NewRtuReadRequest(addr, global.ReadInputs, offset, num), crcOrder)
	if err != nil {
		return 0, err
	}
	return m.Read(p)
}

// ReadHoldingRegisters 读取保持寄存器
func (m *RtuMaster) ReadHoldingRegisters(p []byte, addr byte, offset, num uint16, crcOrder binary.ByteOrder) (int, error) {
	_, err := m.Write(request.NewRtuReadRequest(addr, global.ReadHoldingRegisters, offset, num), crcOrder)
	if err != nil {
		return 0, err
	}
	return m.Read(p)
}

// ReadInputRegisters 读取输入寄存器
func (m *RtuMaster) ReadInputRegisters(p []byte, addr byte, offset, num uint16, crcOrder binary.ByteOrder) (int, error) {
	_, err := m.Write(request.NewRtuReadRequest(addr, global.ReadInputRegisters, offset, num), crcOrder)
	if err != nil {
		return 0, err
	}
	return m.Read(p)
}

// WriteSingleCoil 写单个线圈
func (m *RtuMaster) WriteSingleCoil(addr byte, offset uint16, on bool, crcOrder binary.ByteOrder) error {
	var data uint16
	if on {
		data = 0xff00
	}
	_, err := m.Write(request.NewRtuWriteSingleRequest(addr, global.WriteSingleCoil, offset, data), crcOrder)
	if err != nil {
		return err
	}
	_, err = m.Read(make([]byte, 0))
	return err
}

// WriteSingleRegister 写单个保持寄存器
func (m *RtuMaster) WriteSingleRegister(addr byte, offset uint16, data uint16, crcOrder binary.ByteOrder) error {
	_, err := m.Write(request.NewRtuWriteSingleRequest(addr, global.WriteSingleRegister, offset, data), crcOrder)
	if err != nil {
		return err
	}
	_, err = m.Read(make([]byte, 0))
	return err
}

// WriteMultiCoils 写多个线圈
func (m *RtuMaster) WriteMultiCoils(addr byte, offset uint16, on []bool, crcOrder binary.ByteOrder) error {
	coilNum := len(on)
	// 线圈数除以8再向上取整，得到字节数
	data := make([]byte, (coilNum-1)/8+1)
	for i, b := range on {
		if b {
			idx := i / 8
			data[idx] += 1 << (i - 8*idx)
		}
	}
	_, err := m.Write(request.NewRtuWriteMultiCoilsRequest(addr, offset, uint16(coilNum), data), crcOrder)
	if err != nil {
		return err
	}
	_, err = m.Read(make([]byte, 0))
	return err
}

// WriteMultiRegisters 写多个保持寄存器
func (m *RtuMaster) WriteMultiRegisters(addr byte, offset uint16, data []uint16, crcOrder binary.ByteOrder) error {
	_, err := m.Write(request.NewRtuWriteMultiRegsRequest(addr, offset, data), crcOrder)
	if err != nil {
		return err
	}
	_, err = m.Read(make([]byte, 0))
	return err
}

// ---- 南瑞项目变体 ----

// NRWriteMultiRegisters 写多个保持寄存器
func (m *RtuMaster) NRWriteMultiRegisters(addr byte, offset uint16, data []uint16, crcOrder binary.ByteOrder) error {
	_, err := m.Write(request.NewNRWriteMultiRegsRequest(addr, offset, data), crcOrder)
	if err != nil {
		return err
	}
	_, err = m.Read(make([]byte, 0))
	return err
}
