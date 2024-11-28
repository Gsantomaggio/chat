from wire_formatting import (
    read_header_components,
    read_uint32,
    read_string,
)


users = set()
def read_message(buffer: bytes, offset: int = 0):
    msg_len, offset = read_uint32(buffer, offset)
    msg_len_rcv = len(buffer[offset:])
    if msg_len < msg_len_rcv:
        return f"Message not correct, declared len {msg_len}, but received len {msg_len_rcv}"
    return "Message length ok!"


def read_header(buffer: bytes, offset: int):
    _, command, offset = read_header_components(buffer, offset)
    if command == 1:
        return "command: login"
    elif command == 2:
        return "command: message"
    else:
        return f"Error command in the header: {command}"


def read_correlationId(buffer: bytes, offset: int):
    return read_uint32(buffer, offset)


def login(buffer: bytes, offset: int):
    username, offset = read_string(buffer, offset)
    if username not in users:
        users.add(username)
        return f"user: {username} login confirmed"
    else:
        return f"user: {username} already logged"

print(login(b"\x00\x05\x75\x73\x65\x72\x31", 0))
