
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
