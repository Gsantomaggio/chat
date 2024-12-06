mex_length = b"\x00\x00\x00\x0e"
version = b"\x01"
key_login = b"\x00\x01"
key_message = b"\x00\x02"
correlationId = b"\x00\x00\x00\x01"
length = b"\x00\x05"
user1 = bytes("user1", "utf-8")
user2 = bytes("user2", "utf-8")
message = bytes("testo di prova", "utf-8")
fromfield = bytes("user1", "utf-8")
tofield = bytes("user2", "utf-8")
timestamp = (1733236118).to_bytes(8)

login_message = mex_length + version + key_login + correlationId + length + user1
login_message2 = mex_length + version + key_login + correlationId + length + user2
mex = b"\x00\x0e" + message + length + fromfield + length + tofield + timestamp
message_message = b"\x00\x00\x00\x2d" + version + key_message + correlationId + mex
