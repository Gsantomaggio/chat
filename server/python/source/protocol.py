import time

version = b"\x01"
key_login = b"\x00\x01"
key_message = b"\x00\x02"
correlationId = b"\x00\x00\x00\x01"


def create_login_message(username):
    user = bytes(username, "utf-8")
    length = len(user).to_bytes(2)
    mex = version + key_login + correlationId + length + user
    mex_length = len(mex).to_bytes(4)

    return mex_length + mex


def create_msg_message(msg, to, frm):
    prefix = version + key_message + correlationId
    message = bytes(msg, "utf-8")
    message_length = len(message).to_bytes(2)
    from_field = bytes(frm, "utf-8")
    from_length = len(from_field).to_bytes(2)
    to_field = bytes(to, "utf-8")
    to_length = len(to_field).to_bytes(2)
    timestamp = int(time.time()).to_bytes(8)
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
    mex_length = len(mex).to_bytes(4)

    return mex_length + mex
