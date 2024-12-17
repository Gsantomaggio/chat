class Message:
    def __init__(self, correlationId, message, from_field, to_field, timestamp):
        self.correlationId = correlationId
        self.message = message
        self.from_field = from_field
        self.to_field = to_field
        self.timestamp = timestamp

    def __str__(self) -> str:
        return (
            f"\nmessage: {self.message}\nfrom: {self.from_field}\nto: {self.to_field}"
        )
