from datetime import datetime, timedelta, timezone


def read_header_components(buffer: bytes, offset: int):
    version, offset = read_uint8(buffer, offset)
    command, offset = read_uint16(buffer, offset)
    return version, command, offset


def read_string(buffer: bytes, offset: int) -> tuple:
    length, offset = read_uint16(buffer, offset)
    data_string = bytes(buffer[offset:offset+length]).decode(errors="ignore")
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
    date_time_offset = datetime_from_unix_milliseconds(value)
    return date_time_offset, offset


def datetime_from_unix_milliseconds(ms: int) -> datetime:
    delta = timedelta(milliseconds=ms)
    utc_epoch = datetime(1970, 1, 1, tzinfo=timezone.utc)
    dt_with_offset = utc_epoch + delta

    return dt_with_offset