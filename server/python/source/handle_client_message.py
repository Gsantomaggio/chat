from socket import socket
from wire_formatting import (
    read_header,
    read_uint32,
    read_string,
    read_timestamp,
)
from message import Message
from users import login


def read_message(buffer: bytes, conn: socket, is_logged_correctly: bool):
    _, offset = read_message_length(buffer)
    _, key, offset = read_header(buffer, offset)
    correlationId, offset = read_correlationId(buffer, offset)
    if key == 1:
        username = login(buffer, offset, conn)
        return username, "CommandLogin"
    elif key == 2:
        if is_logged_correctly:
            message = read_command_message(buffer, offset, correlationId)
            return message, "CommandMessage"
        else:
            raise ValueError(
                "Message sent without a login. Please send a CommandLogin message"
            )
    else:
        raise ValueError(f"Error command in the header. Key: {key}")


def read_message_length(buffer: bytes, offset: int = 0) -> tuple:
    length, offset = read_uint32(buffer, offset)
    msg_len_rcv = len(buffer[offset:])
    if length < msg_len_rcv:
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
