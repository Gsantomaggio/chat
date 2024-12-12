from socket import socket
from source.wire_formatting import (
    read_header,
    read_uint32,
    read_string,
    read_timestamp,
    write_uint16,
    write_uint32,
    write_uint8,
)
from source.message import Message
from source.users import User

from source import Logger

logger = Logger(__name__)


def read_message(buffer: bytes, conn: socket, user: User | None, users: dict) -> User:
    _, offset = read_message_length(buffer)
    _, key, offset = read_header(buffer, offset)
    correlationId, offset = read_correlationId(buffer, offset)

    if key == 1:
        username, _ = read_string(buffer, offset)
        usr: User = users.setdefault(username, User(username))
        response_code = usr.login(conn)
        logger.info(f"User {usr.username} logged in")

        send_response(correlationId, response_code, usr)
        send_user_messages(usr)

        return usr

    elif key == 2:
        if user:
            response_code = 1
            send_response(correlationId, response_code, user)
            message = read_command_message(buffer, offset, correlationId)
            send_message(message, users)
        else:
            send_response(correlationId, 3, user)
            message = None

        return user

    else:
        raise ValueError(f"Received wrong COMMAND in the header. KEY: {key}")


def read_message_length(buffer: bytes, offset: int = 0) -> tuple:
    length, offset = read_uint32(buffer, offset)
    msg_len_rcv = len(buffer[offset:])
    if length != msg_len_rcv:
        raise ValueError(
            f"Message not correct, declared len {length}, but received len {msg_len_rcv}"
        )

    return length, offset


def read_correlationId(buffer: bytes, offset: int) -> int:
    return read_uint32(buffer, offset)


def read_command_message(buffer: bytes, offset: int, correlationId: int) -> Message:
    message_field, offset = read_string(buffer, offset)
    from_field, offset = read_string(buffer, offset)
    to_field, offset = read_string(buffer, offset)
    timestamp, offset = read_timestamp(buffer, offset)

    return Message(correlationId, message_field, from_field, to_field, timestamp)


def create_command_message(m: Message) -> bytes:
    version = write_uint8(1)
    key = write_uint16(2)
    correlationId = write_uint32(m.correlationId)
    prefix = version + key + correlationId
    message = bytes(m.message, "utf-8")
    message_length = write_uint16(len(message))
    from_field = bytes(m.from_field, "utf-8")
    from_length = write_uint16(len(from_field))
    to_field = bytes(m.to_field, "utf-8")
    to_length = write_uint16(len(to_field))
    timestamp = int(m.timestamp).to_bytes(8)
    mex = (
        prefix
        + message_length
        + message
        + from_length
        + from_field
        + to_length
        + to_field
        + timestamp
    )
    mex_length = write_uint32(len(mex))

    return mex_length + mex


def send_message(m: Message, users: dict) -> None:
    receiver = m.to_field
    user: User = users.setdefault(receiver, User(receiver))
    user.messages.append(m)
    if user.isonline:
        send_user_messages(user)
    else:
        logger.warning(f"User {user.username} is offline and received a message from {m.from_field}")


def send_user_messages(user: User) -> None:
    while user.messages:
        m: Message = user.messages.pop(0)
        mex = create_command_message(m)
        user.conn.send(mex)
        
        logger.info(f"Sent message from {m.from_field} to {user.username}: {m.message}")


def send_response(correlationId: int, code: int, user: User) -> None:
    version = write_uint8(1)
    key = write_uint16(3)
    corrId = write_uint32(correlationId)
    response_code = write_uint16(code)
    resp = version + key + corrId + response_code
    resp_length = len(resp).to_bytes(4)
    response = resp_length + resp
    user.conn.send(response)

    logger.debug(f"Response sent to user {user.username} with correlationId {correlationId}")