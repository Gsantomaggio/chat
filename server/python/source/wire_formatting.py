from datetime import datetime, timedelta, timezone


def read_header(buffer: bytes, offset: int):
    version, offset = read_uint8(buffer, offset)
    command, offset = read_uint16(buffer, offset)
    return version, command, offset


def read_string(buffer: bytes, offset: int) -> tuple:
    length, offset = read_uint16(buffer, offset)
    data_string = bytes(buffer[offset : offset + length]).decode(errors="ignore")
    offset += length

    return data_string, offset


def read_uint64(buffer: bytes, offset: int) -> tuple:
    total_move_offset = offset + 8
    data = int.from_bytes(buffer[offset:total_move_offset], "big")

    return data, total_move_offset


def read_uint32(buffer: bytes, offset: int) -> tuple:
    total_move_offset = offset + 4
    data = int.from_bytes(buffer[offset:total_move_offset], "big")

    return data, total_move_offset


def read_uint16(buffer: bytes, offset: int) -> tuple:
    total_move_offset = offset + 2
    data = int.from_bytes(buffer[offset:total_move_offset], "big")

    return data, total_move_offset


def read_uint8(buffer: bytes, offset: int) -> tuple:
    total_move_offset = offset + 1
    data = int.from_bytes(buffer[offset:total_move_offset], "big")

    return data, total_move_offset


def read_timestamp(buffer: bytes, offset: int) -> tuple:
    value, offset = read_uint64(buffer, offset)
    return value, offset


def write_uint8(num: int) -> bytes:
    return num.to_bytes()


def write_uint16(num: int) -> bytes:
    return num.to_bytes(2)


def write_uint32(num: int) -> bytes:
    return int(num).to_bytes(4)
