from wire_formatting import (
    read_uint16,
    read_uint32,
)


users = {}
def read_message(buffer: bytes, offset: int = 0):
    msg_len, offset = read_uint32(buffer, offset)
    msg_len_rcv = len(buffer[offset:])
    if msg_len < msg_len_rcv:
        return f"Message not correct, declared len {msg_len}, but received len {msg_len_rcv}"


print(read_message(b"\x00\x00\x00\x01\x00\x01", 0))
