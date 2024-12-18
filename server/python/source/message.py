class Message:
    """
    Message class represents a message with its metadata.

    Attributes:
        correlationId (str): The unique identifier for the message.
        message (str): The content of the message.
        from_field (str): The sender of the message.
        to_field (str): The recipient of the message.
        timestamp (float): The time the message was created.

    Methods:
        __str__(): Returns a string representation of the message.
    """

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
