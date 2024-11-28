from wire_formatting import read_string

users = set()

def login(buffer: bytes, offset: int) -> tuple:
    username, offset = read_string(buffer, offset)
    if username not in users:
        users.add(username)
        return False, username, offset
    else:
        return True, username, offset