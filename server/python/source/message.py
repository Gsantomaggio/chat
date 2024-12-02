class Message:
    def __init__(self, correlationId, message, from_field, to_field, timestamp):
        self.correlationId = correlationId
        self.message = message
        self.from_field = from_field
        self.to_field = to_field
        self.timestamp = timestamp


if __name__ == "__main__":
    buf = b"\x00\x05\x75\x03\x01\x00\x01\x00\x01\x75\x73\x00\x05\x75\x73\x65\x72\x31\x65\x72\x31\x73\x65\x72\x31"
    mex = Message(buf)
    print(mex.message_field)
