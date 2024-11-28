from wire_formatting import read_string

users = {}


def login(buffer: bytes, offset: int) -> tuple:
    username, offset = read_string(buffer, offset)
    status = users.setdefault(username, "online")
    if status == "online":
        return True, username, offset
    return False, username, offset


def logout(username: str) -> None:
    users[username] = "offline"
