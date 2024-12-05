class Message:
    def __init__(self, correlationId, message, from_field, to_field, timestamp):
        self.correlationId = correlationId
        self.message = message
        self.from_field = from_field
        self.to_field = to_field
        self.timestamp = timestamp
