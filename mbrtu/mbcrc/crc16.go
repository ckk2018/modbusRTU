package mbcrc

type crc16Table = [256]uint16

var table crc16Table

// 在初始化函数中初始化crc16查找表
func init() {
	getCrcTable(&table)
}

// 计算crc16查找表
func getCrcTable(t *crc16Table) {
	var lsb uint16
	var crc uint16

	for i := uint16(0); i < 256; i++ {
		crc = i
		for j := 0; j < 8; j++ {
			lsb = crc & 0x0001
			crc >>= 1
			if lsb == 0 {
				continue
			}
			crc ^= 0xa001
		}
		t[i] = crc
	}
}

// Crc16 计算数据的crc16校验码
func Crc16(data []byte) uint16 {
	var crc16 uint16 = 0xffff
	for _, d := range data {
		idx := byte(crc16) ^ d
		crc16 = (crc16 >> 8) ^ table[idx]
	}
	return crc16
}
