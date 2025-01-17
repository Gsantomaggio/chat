from socket import socket
from source.exceptions import AlreadyLoggedException
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


"""
This module handles user messages in a TCP server environment.

Classes:
    MessageHandler: A class to handle reading, processing, and sending messages between users.

Functions:
    read_message(buffer: bytes, user: User | None) -> User: Reads and processes a message from the buffer.
    _read_correlationId(buffer: bytes, offset: int) -> int: Reads the correlation ID from the buffer.
    _read_command_message(buffer: bytes, offset: int, correlationId: int) -> Message: Reads a command message from the buffer.
    _create_command_message(m: Message) -> bytes: Creates a command message in bytes format.
    _send_message(m: Message) -> int: Sends a message to the intended recipient.
    _send_user_messages(user: User) -> list: Sends all queued messages to the user.
    _send_response(correlationId: int, code: int, conn: socket) -> bytes: Sends a response back to the client.

Usage:
    This module is used to handle user login, message sending, and response handling in a TCP server. It ensures that messages are correctly formatted, sent, and logged.

Example:
    handler = MessageHandler(conn, users)
    user = handler.read_message(buffer, user)
"""


class MessageHandler:
    def __init__(self, conn: socket, users: dict):
        self.conn = conn
        self.users = users

    def read_message(self, buffer: bytes, user: User | None) -> User:
        _, key, offset = read_header(buffer, 0)
        correlationId, offset = self._read_correlationId(buffer, offset)

        if key == 1:
            username, _ = read_string(buffer, offset)
            usr: User = self.users.setdefault(username, User(username))
            response_code = usr.login(self.conn)
            if response_code == 1:
                logger.info(f"User {usr.username} logged in")
                self._send_response(correlationId, response_code, usr.conn)
                self._send_user_messages(usr)
            else:
                self._send_response(correlationId, response_code, self.conn)
                raise AlreadyLoggedException(f"User {usr.username} already logged")

            return usr

        elif key == 2:
            if user:
                message = self._read_command_message(buffer, offset, correlationId)
                response_code = self._send_message(message)
                self._send_response(correlationId, response_code, user.conn)
            else:
                self._send_response(correlationId, 3, self.conn)

            return user

        else:
            raise ValueError(f"Received wrong COMMAND in the header. KEY: {key}")

    def _read_correlationId(self, buffer: bytes, offset: int) -> int:
        return read_uint32(buffer, offset)

    def _read_command_message(
        self, buffer: bytes, offset: int, correlationId: int
    ) -> Message:
        message_field, offset = read_string(buffer, offset)
        from_field, offset = read_string(buffer, offset)
        to_field, offset = read_string(buffer, offset)
        timestamp, offset = read_timestamp(buffer, offset)

        return Message(correlationId, message_field, from_field, to_field, timestamp)

    def _create_command_message(self, m: Message) -> bytes:
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

    def _send_message(self, m: Message) -> int:
        receiver = m.to_field
        try:
            user = self.users[receiver]
            user.messages.append(m)
            if user.isonline:
                self._send_user_messages(user)
            else:
                logger.warning(
                    f"User {user.username} is offline and received a message from {m.from_field}"
                )
            return 1
        except KeyError:
            logger.error(f"User {receiver} not found")
            return 3

    def _send_user_messages(self, user: User) -> list:
        while user.messages:
            m: Message = user.messages.pop(0)
            mex = self._create_command_message(m)
            user.conn.send(mex)

            logger.info(
                f"Sent message from {m.from_field} to {user.username}: {m.message}"
            )

        return user.messages

    def _send_response(self, correlationId: int, code: int, conn: socket) -> bytes:
        version = write_uint8(1)
        key = write_uint16(3)
        corrId = write_uint32(correlationId)
        response_code = write_uint16(code)
        resp = version + key + corrId + response_code
        resp_length = write_uint32(9)
        response = resp_length + resp
        conn.send(response)

        logger.debug(f"Response sent with correlationId {correlationId}")

        return response
