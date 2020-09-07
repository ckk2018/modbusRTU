package global

// FunCode modbus功能码类型别名
type FunCode = byte

// 支持的modbus功能码
const (
	ReadCoils            FunCode = 0x01
	ReadInputs           FunCode = 0x02
	ReadHoldingRegisters FunCode = 0x03
	ReadInputRegisters   FunCode = 0x04
	WriteSingleCoil      FunCode = 0x05
	WriteSingleRegister  FunCode = 0x06
	WriteMultiCoils      FunCode = 0x0f
	WriteMultiRegisters  FunCode = 0x10
)

// SlaveError modbus从站异常类型
type SlaveError struct {
	Code byte
	msg  string
}

func (e SlaveError) Error() string {
	return e.msg
}

// SlaveErrorMap 备用的modbus异常
var SlaveErrorMap = map[byte]SlaveError{
	0x01: {0x01, "illegal function"},
	0x02: {0x02, "illegal data address"},
	0x03: {0x03, "illegal data value"},
	0x04: {0x04, "slave device failure"},
	0x05: {0x05, "acknowledge"},
	0x06: {0x06, "slave device busy"},
	0x0a: {0x0a, "gateway path unavailable"},
	0x0b: {0x0b, "gateway target device failed to respond"},
}
