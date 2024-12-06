from socket import socket
from source.wire_formatting import (
    read_header,
    read_uint32,
    read_string,
    read_timestamp,
    write_uint16
)
from source.message import Message
from source.users import login, User


def read_message(buffer: bytes, conn: socket, user: User | None) -> tuple:
    _, offset = read_message_length(buffer)
    _, key, offset = read_header(buffer, offset)
    correlationId, offset = read_correlationId(buffer, offset)
    if key == 1:
        response_code, user = login(buffer, offset, conn)
        send_response(correlationId, response_code, user)
        send_user_messages(user)

        return user, "CommandLogin"

    elif key == 2:
        if user:
            response_code = 1
            send_response(correlationId, response_code, user)
            message = read_command_message(buffer, offset, correlationId)
        else:
            send_response(correlationId, 3, user)
            message = None
        
        return message, "CommandMessage"

    else:
        raise ValueError(f"Error command in the header. Key: {key}")


def read_message_length(buffer: bytes, offset: int = 0) -> tuple:
    length, offset = read_uint32(buffer, offset)
    msg_len_rcv = len(buffer[offset:])
    if length != msg_len_rcv:
        raise ValueError(
            f"Message not correct, declared len {length}, but received len {msg_len_rcv}"
        )
    return length, offset


def read_correlationId(buffer: bytes, offset: int):
    return read_uint32(buffer, offset)


def read_command_message(buffer: bytes, offset: int, correlationId: int) -> Message:
    message_field, offset = read_string(buffer, offset)
    from_field, offset = read_string(buffer, offset)
    to_field, offset = read_string(buffer, offset)
    timestamp, offset = read_timestamp(buffer, offset)
    timestamp = timestamp.strftime("%d-%m-%Y %H:%M:%S")

    return Message(correlationId, message_field, from_field, to_field, timestamp)


def send_user_messages(user: User) -> None:
    while user.messages:
        mex = user.messages.pop()
        user.conn.send(str(mex).encode())


def send_response(correlationId: int, code: int, user: User) -> None:
    version = (1).to_bytes()
    key = write_uint16(3)
    response_code = write_uint16(code)
    response = version+key+correlationId+response_code
    user.conn.send(response)
