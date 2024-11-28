from wire_formatting import (
    read_header_components,
    read_uint32,
    read_string,
    read_timestamp,
)

from login import login


class Message:
    def __init__(self, buffer: bytes, user: str = "") -> None:
        self.buffer = None
        if self.buffer is None:
            self.buffer = buffer
        self.offset = 0
        self.length = None
        self.version = None
        self.key = None
        self.command = None
        self.correlationId = None
        self.username = user if user != "" else None
        self.message_field = None
        self.from_field = None
        self.to_field = None
        self.timestamp = None

        self.read_message_length()
        self.read_header()
        self.read_correlationId()

    def read_message_length(self):
        self.length, self.offset = read_uint32(self.buffer, self.offset)
        msg_len_rcv = len(self.buffer[self.offset :])
        if self.length < msg_len_rcv:
            raise ValueError(
                f"Message not correct, declared len {self.length}, but received len {msg_len_rcv}"
            )

    def read_header(self):
        self.version, self.key, self.offset = read_header_components(
            self.buffer, self.offset
        )
        match self.key:
            case 1:
                self.command = "CommandLogin"
                self.read_command_login()
            case 2:
                self.command = "CommandMessage"
                if self.username:
                    self.read_command_message()
                else:
                    raise Exception(
                        "Message sent without a login. Please send a CommandLogin message"
                    )
            case _:
                raise ValueError(f"Error command in the header. Key: {self.key}")

    def read_correlationId(self):
        self.correlationId, self.offset = read_uint32(self.buffer, self.offset)

    def read_command_login(self):
        is_logged_already, self.username, self.offset = login(self.buffer, self.offset)
        if is_logged_already:
            raise ValueError(f"user: {self.username} already logged!")

    def read_command_message(self):
        self.message_field, self.offset = read_string(self.buffer, self.offset)
        self.from_field, self.offset = read_string(self.buffer, self.offset)
        self.to_field, self.offset = read_string(self.buffer, self.offset)
        timestamp, self.offset = read_timestamp(self.buffer, self.offset)
        self.timestamp = timestamp.strftime("%d-%m-%Y %H:%M:%S")


if __name__ == "__main__":
    buf = b"\x00\x05\x75\x03\x01\x00\x01\x00\x01\x75\x73\x00\x05\x75\x73\x65\x72\x31\x65\x72\x31\x73\x65\x72\x31"
    mex = Message(buf)
    print(mex.message_field)
