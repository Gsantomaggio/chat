# import struct
# from datetime import datetime, timedelta, timezone


def read_header(buffer: bytes, offset: int):
    version = buffer[offset : offset + 1]
    offset += 1
    command = buffer[offset : offset + 2]
    offset += 2
    return version, command, offset


def read_string(buffer: bytes, offset: int) -> tuple:
    total_move_offset = offset + 2
    length = int.from_bytes(buffer[offset:total_move_offset], "big")
    offset += total_move_offset + length
    data_string = bytes(buffer[total_move_offset:offset]).decode(errors="ignore")

    return data_string, offset


def read_uint32(buffer: bytes, offset: int) -> tuple:
    total_move_offset = offset + 4
    data = int.from_bytes(buffer[offset:total_move_offset], "big")

    return data, total_move_offset


def read_uint16(buffer: bytes, offset: int) -> tuple:
    total_move_offset = offset + 2
    data = int.from_bytes(buffer[offset:total_move_offset], "big")

    return data, total_move_offset


"""
def read_any(buffer: bytes, offset: int) -> tuple:
    signed = False
    data_type = buffer[offset]
    offset += 1
    if data_type is FormatCode.Sym8:
        move_offset = 1
        return read_string(buffer, offset, move_offset)
    elif data_type is FormatCode.Str8:
        move_offset = 1
        return read_string(buffer, offset, move_offset)
    elif data_type is FormatCode.Sym32:
        move_offset = 4
        return read_string(buffer, offset, move_offset)
    elif data_type is FormatCode.Str32:
        move_offset = 4
        return read_string(buffer, offset, move_offset)
    elif data_type is FormatCode.Ulong0:
        return 0, offset
    elif data_type is FormatCode.Uint0:
        return 0, offset
    elif data_type is FormatCode.SmallUlong or data_type is FormatCode.SmallUint:
        move_offset = 1
    elif data_type is FormatCode.Ulong:
        move_offset = 8
    elif data_type is FormatCode.Uint:
        move_offset = 4
    elif data_type is FormatCode.Ushort:
        move_offset = 2
    elif data_type is FormatCode.Ubyte:
        move_offset = 1
    elif data_type is FormatCode.Byte:
        signed = True
        move_offset = 1
    elif data_type is FormatCode.Vbin8:
        move_offset = 1
        return read_bytes(buffer, offset, move_offset)
    elif data_type is FormatCode.Vbin32:
        move_offset = 4
        return read_bytes(buffer, offset, move_offset)
    elif data_type is FormatCode.Smallint or data_type is FormatCode.Smalllong:
        signed = True
        move_offset = 1
    elif data_type is FormatCode.Int:
        signed = True
        move_offset = 4
    elif data_type is FormatCode.Long:
        signed = True
        move_offset = 8
    elif data_type is FormatCode.Short:
        signed = True
        move_offset = 2
    elif data_type is FormatCode.Bool:
        move_offset = 1
        return read_bool(buffer, offset, move_offset)
    elif data_type is FormatCode.BoolTrue:
        return True, offset
    elif data_type is FormatCode.BoolFalse:
        return False, offset
    elif data_type is FormatCode.Timestamp:
        return read_timestamp(buffer, 0)
    elif data_type is FormatCode.Float:
        data_float_double = ">f"
        move_offset = 4
        return read_float_double(buffer, offset, move_offset, data_float_double)
    elif data_type is FormatCode.Double:
        data_float_double = ">d"
        move_offset = 8
        return read_float_double(buffer, offset, move_offset, data_float_double)
    elif data_type is FormatCode.Null:
        return None, offset
    else:
        raise ValueError("Error in read_any, value {} not found".format(data_type))

    return read_generic(buffer, offset, move_offset, signed)


def read_string(buffer: bytes, offset: int, move_offset: int) -> tuple:
    total_move_offset = offset + move_offset
    length = int.from_bytes(buffer[offset:total_move_offset], "big")
    length += total_move_offset
    data_string = bytes(buffer[total_move_offset:length]).decode(errors="ignore")

    return data_string, length


def read_bool(buffer: bytes, offset: int, move_offset: int) -> tuple:
    total_move_offset = offset + move_offset
    value = int.from_bytes(buffer[offset:total_move_offset], "big")

    return value != 0, total_move_offset


def read_bytes(buffer, offset, move_offset):
    total_move_offset = offset + move_offset
    length = int.from_bytes(buffer[offset:total_move_offset], "big")
    length += total_move_offset
    binary_data = buffer[total_move_offset:length]

    return binary_data, length


def read_float_double(
    buffer: bytes, offset: int, move_offset: int, data_type: str
) -> tuple:
    total_move_offset = offset + move_offset
    data_float = struct.unpack(data_type, buffer[offset:total_move_offset])[0]

    return data_float, total_move_offset


def read_generic(
    buffer: bytes, offset: int, move_offset: int, signed: bool = False
) -> tuple:
    total_move_offset = offset + move_offset
    data_generic = int.from_bytes(
        buffer[offset:total_move_offset], "big", signed=signed
    )

    return data_generic, total_move_offset


def read_composite_header(buffer: bytes, offset: int) -> tuple:
    data_type = buffer[offset]

    if data_type != 0:
        raise ValueError("Invalid composite header %#02x {type}")
    offset += 1
    _, offset = read_any(buffer, offset)
    fields, offset = read_list_header(buffer, offset)
    return offset, fields


def read_list_header(buffer: bytes, offset: int) -> tuple:
    data_type = buffer[offset]

    if data_type is FormatCode.List0:
        length = 0
        offset += 1
        return length, offset
    elif data_type is FormatCode.List8:
        # move the offset by two to discard unused size value
        offset += 2
        length = int.from_bytes(buffer[offset : offset + 1], "big")
        offset += 1
        return length, offset
    elif data_type is FormatCode.List32:
        # move the offset by five to discard unused size value
        offset += 5
        length = int.from_bytes(buffer[offset : offset + 4], "big")
        offset += 4
        return length, offset
    else:
        raise ValueError("Invalid list header {}".format(data_type))


def try_read_null(byte):
    offset = value = 1 if byte == FormatCode.Null else 0
    return offset, value


def read_timestamp(buffer: bytes, offset: int) -> tuple:
    data_type = buffer[offset]
    if data_type is FormatCode.Timestamp:
        offset += 1
        value = int.from_bytes(buffer[offset : offset + 8], "big")
        offset += 8
        date_time_offset = datetime_from_unix_milliseconds(value)
        return date_time_offset, offset

    raise ValueError("Read_timestamp invalid type {type}")


def datetime_from_unix_milliseconds(ms):
    delta = timedelta(milliseconds=ms)
    utc_epoch = datetime(1970, 1, 1, tzinfo=timezone.utc)
    dt_with_offset = utc_epoch + delta

    return dt_with_offset
"""
