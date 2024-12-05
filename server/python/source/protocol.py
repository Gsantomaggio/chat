mex_length = b"\x00\x00\x00\x0e"
version = b"\x01"
key_login = b"\x00\x01"
key_message = b"\x00\x02"
correlationId = b"\x00\x00\x00\x01"
length = b"\x00\x05"
user1 = bytes("user1", "utf-8")
user2 = b"\x00\x05\x75\x73\x65\x72\x32"
message = "testo di prova".encode().hex()
fromfield = "user1".encode().hex()
tofield = "user2".encode().hex()
timestamp = (1733236118).to_bytes(8)

login_message = mex_length + version + key_login + correlationId + length + user1

print(login_message)
