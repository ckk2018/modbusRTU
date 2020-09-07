# -*- coding: utf-8 -*-

import logging
from typing import List
from ctypes import CDLL, c_int, c_char, c_char_p, c_ubyte, c_uint, \
    c_ushort, string_at

from IPython import embed

BIG_ENDIAN = b'>'
LITTLE_ENDIAN = b'<'

ErrString = c_ubyte * 1024
ReadData = c_ubyte * 1024
WriteData = c_ushort * 512


class ModbusDll:
    def __init__(self, dll_path: str):
        self._dll = CDLL(dll_path)

        self._err_string = ErrString()
        self._read_data = ReadData()

        # 初始化所有要测试的函数
        self._init_func()

    def _init_func(self):
        self._mb_open = self._dll.Open
        self._mb_open.restype = c_int
        self._mb_open.argtypes = (c_char_p, c_uint, c_ubyte, c_char, c_ubyte,
                                  c_uint, ErrString)

        self._mb_close = self._dll.Close
        self._mb_close.restype = c_int
        self._mb_close.argtypes = (ErrString, )

        self._mb_read_holding_registers = self._dll.ReadHoldingRegisters
        self._mb_read_holding_registers.restype = c_int
        self._mb_read_holding_registers.argtypes = (ReadData, c_ubyte,
                                                    c_ushort, c_ushort, c_char,
                                                    ErrString)

        self._mb_write_single_register = self._dll.WriteSingleRegister
        self._mb_write_single_register.restype = c_int
        self._mb_write_single_register.argtypes = (c_ushort, c_ubyte, c_ushort,
                                                   c_char, ErrString)

        self._mb_write_multi_registers = self._dll.WriteMultiRegisters
        self._mb_write_multi_registers.restype = c_int
        self._mb_write_multi_registers.argtypes = (WriteData, c_int, c_ubyte,
                                                   c_ushort, c_char, ErrString)

        self._mb_nr_write_multi_registers = self._dll.NRWriteMultiRegisters
        self._mb_nr_write_multi_registers.restype = c_int
        self._mb_nr_write_multi_registers.argtypes = (WriteData, c_int,
                                                      c_ubyte, c_ushort,
                                                      c_char, ErrString)

    def _last_err_msg(self) -> str:
        return string_at(self._err_string).decode('utf-8')

    def _last_read_data(self, length: int) -> bytes:
        return string_at(self._read_data, length)

    def open(self, name: str, timeout: int) -> bool:
        res = self._mb_open(name.encode('utf-8'), 9600, 8, b'n', 1, timeout,
                            self._err_string)
        if res < 0:
            logging.error(self._last_err_msg())
            return False
        return True

    def close(self) -> bool:
        res = self._mb_close(self._err_string)
        if res < 0:
            logging.error(self._last_err_msg())
            return False
        return True

    def read_holding_registers(self, addr: int, offset: int, num: int,
                               crc_order: bytes) -> bytes:
        res = self._mb_read_holding_registers(self._read_data, addr, offset,
                                              num, crc_order, self._err_string)
        if res < 0:
            logging.error(self._last_err_msg())
            return b''
        return self._last_read_data(length=res)

    def write_single_register(self, data: int, addr: int, offset: int,
                              crc_order: bytes) -> bool:
        res = self._mb_write_single_register(data, addr, offset, crc_order,
                                             self._err_string)
        if res < 0:
            logging.error(self._last_err_msg())
            return False
        return True

    def write_multi_registers(self, data: List[int], addr: int, offset: int,
                              crc_order: bytes) -> bool:
        res = self._mb_write_multi_registers(WriteData(*data), len(data), addr,
                                             offset, crc_order,
                                             self._err_string)
        if res < 0:
            logging.error(self._last_err_msg())
            return False
        return True

    def nr_write_multi_registers(self, data: List[int], addr: int, offset: int,
                                 crc_order: bytes) -> bool:
        res = self._mb_nr_write_multi_registers(WriteData(*data), len(data),
                                                addr, offset, crc_order,
                                                self._err_string)
        if res < 0:
            logging.error(self._last_err_msg())
            return False
        return True


if __name__ == '__main__':
    logging.basicConfig(
        format='%(asctime)s %(levelname)s:%(message)s',
        level=logging.INFO,
    )

    try:
        dll = ModbusDll('./mbrtu.dll')
    except Exception as e:
        logging.error(e)
        input('输入任意键退出...')
    else:
        embed()
