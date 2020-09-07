package main

import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
	"unsafe"

	"github.com/tarm/serial"

	"ckklearn.com/testmodbus/mbrtu"
)

const (
	errNotOpen         string = "should call `Open` first"
	errNotClose        string = "should call `Close` first"
	errInvalidParity   string = "invalid parity"
	errInvalidSBits    string = "invalid sBits"
	errInvalidCrcOrder string = "invalid crcOrder"
	retErr             C.int  = -1
)

// 全局的modbus主站对象
var master *mbrtu.RtuMaster

// Open 打开串口
//export Open
func Open(name *C.char, baud C.uint, dBits C.uchar, parity C.char, sBits C.uchar, readTimeout C.uint, errMsg *C.uchar) C.int {
	if master != nil {
		writeString(errMsg, errNotClose)
		return retErr
	}

	var (
		_parity   serial.Parity
		_stopBits serial.StopBits
	)

	switch parity {
	case 'n':
		_parity = serial.ParityNone
	case 'o':
		_parity = serial.ParityOdd
	case 'e':
		_parity = serial.ParityEven
	default:
		writeString(errMsg, errInvalidParity)
		return retErr
	}

	switch sBits {
	case 1:
		_stopBits = serial.Stop1
	case 2:
		_stopBits = serial.Stop2
	case 15:
		_stopBits = serial.Stop1Half
	default:
		writeString(errMsg, errInvalidSBits)
		return retErr
	}

	rtuMaster, err := mbrtu.NewRtuMaster(&serial.Config{
		Name:        C.GoString(name),
		Baud:        int(baud),
		Size:        byte(dBits),
		Parity:      _parity,
		StopBits:    _stopBits,
		ReadTimeout: time.Second * time.Duration(readTimeout),
	})
	if err != nil {
		writeString(errMsg, err.Error())
		return retErr
	}

	// 如果没有任何异常，则给全局主站变量赋值
	master = rtuMaster

	return 0
}

// Close 关闭串口
//export Close
func Close(errMsg *C.uchar) C.int {
	if master != nil {
		err := master.Close()
		if err != nil {
			writeString(errMsg, err.Error())
			return retErr
		}
		master = nil
	}
	return 0
}

// ReadHoldingRegisters 读取保持寄存器
//export ReadHoldingRegisters
func ReadHoldingRegisters(data *C.uchar, addr C.uchar, offset, num C.ushort, crcOrder C.char, errMsg *C.uchar) C.int {
	if master == nil {
		writeString(errMsg, errNotOpen)
		return retErr
	}

	_crcOrder, err := getCrcOrder(crcOrder)
	if err != nil {
		writeString(errMsg, err.Error())
		return retErr
	}

	p := make([]byte, 1024)
	n, err := master.ReadHoldingRegisters(p, byte(addr), uint16(offset), uint16(num), _crcOrder)
	if err != nil {
		writeString(errMsg, err.Error())
		return retErr
	}

	writeBytes(data, p[:n])

	return C.int(n)
}

// WriteSingleRegister 写单个保持寄存器
//export WriteSingleRegister
func WriteSingleRegister(data C.ushort, addr C.uchar, offset C.ushort, crcOrder C.char, errMsg *C.uchar) C.int {
	if master == nil {
		writeString(errMsg, errNotOpen)
		return retErr
	}

	_crcOrder, err := getCrcOrder(crcOrder)
	if err != nil {
		writeString(errMsg, err.Error())
		return retErr
	}

	err = master.WriteSingleRegister(byte(addr), uint16(offset), uint16(data), _crcOrder)
	if err != nil {
		writeString(errMsg, err.Error())
		return retErr
	}

	return 0
}

// WriteMultiRegisters 写多个连续保持寄存器
//export WriteMultiRegisters
func WriteMultiRegisters(data *C.ushort, dataLen C.int, addr C.uchar, offset C.ushort, crcOrder C.char, errMsg *C.uchar) C.int {
	if master == nil {
		writeString(errMsg, errNotOpen)
		return retErr
	}

	_crcOrder, err := getCrcOrder(crcOrder)
	if err != nil {
		writeString(errMsg, err.Error())
		return retErr
	}

	// 将数据转换成同类型的切片
	reader := bytes.NewReader(C.GoBytes(unsafe.Pointer(data), dataLen*C.sizeof_ushort))
	_data := make([]uint16, dataLen)
	binary.Read(reader, binary.LittleEndian, _data)

	err = master.WriteMultiRegisters(byte(addr), uint16(offset), _data, _crcOrder)
	if err != nil {
		writeString(errMsg, err.Error())
		return retErr
	}

	return 0
}

// ---- 变体 ----

// NRWriteMultiRegisters 写多个连续保持寄存器
//export NRWriteMultiRegisters
func NRWriteMultiRegisters(data *C.ushort, dataLen C.int, addr C.uchar, offset C.ushort, crcOrder C.char, errMsg *C.uchar) C.int {
	if master == nil {
		writeString(errMsg, errNotOpen)
		return retErr
	}

	_crcOrder, err := getCrcOrder(crcOrder)
	if err != nil {
		writeString(errMsg, err.Error())
		return retErr
	}

	// 将数据转换成同类型的切片
	reader := bytes.NewReader(C.GoBytes(unsafe.Pointer(data), dataLen*C.sizeof_ushort))
	_data := make([]uint16, dataLen)
	binary.Read(reader, binary.LittleEndian, _data)

	err = master.NRWriteMultiRegisters(byte(addr), uint16(offset), _data, _crcOrder)
	if err != nil {
		writeString(errMsg, err.Error())
		return retErr
	}

	return 0
}

// ---- 功能函数 ----

func write(dst *C.uchar, src []byte, asString bool) {
	var i int
	hptr := unsafe.Pointer(dst)
	for i = 0; i < len(src); i++ {
		*(*byte)(unsafe.Pointer(uintptr(hptr) + uintptr(i))) = src[i]
	}

	// 如果写入的是字符串信息，则加上结束符
	if asString {
		*(*byte)(unsafe.Pointer(uintptr(hptr) + uintptr(i))) = '\x00'
	}
}

func writeBytes(dst *C.uchar, src []byte) {
	write(dst, src, false)
}

func writeString(dst *C.uchar, src string) {
	write(dst, []byte(src), true)
}

func getCrcOrder(o C.char) (binary.ByteOrder, error) {
	switch o {
	case '>':
		return binary.BigEndian, nil
	case '<':
		return binary.LittleEndian, nil
	default:
		return nil, fmt.Errorf(errInvalidCrcOrder)
	}
}

func main() {}
