package mbrtu

import (
	"fmt"

	"ckklearn.com/testmodbus/global"
)

// RtuParseResponse 解析从站返回报文
func RtuParseResponse(dst, src []byte, reqFunCode global.FunCode) (int, error) {
	// `src[1]` 是读取的报文的功能码
	switch src[1] {
	case reqFunCode:
		// 如果是读取，则写入读取到的数据，否则不写入
		// `src[2]` 是读取的数据的字节长度
		switch reqFunCode {
		case global.ReadCoils, global.ReadInputs, global.ReadInputRegisters, global.ReadHoldingRegisters:
			return copy(dst, src[3:3+src[2]]), nil
		default:
			return 0, nil
		}

	case reqFunCode + 0x80:
		// 判断从站返回的异常类型
		errCode := src[2]
		err, ok := global.SlaveErrorMap[errCode]
		if !ok {
			return 0, fmt.Errorf("unknow error code `%x`", errCode)
		}
		return 0, err

	default:
		return 0, fmt.Errorf("internal error")
	}
}
