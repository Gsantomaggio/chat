class AlreadyLoggedException(Exception):
    """
    AlreadyLoggedException is a custom exception class to handle cases where a user is already logged in.

    Attributes:
        message (str): The error message to be displayed.
        errors (list, optional): Additional errors related to the exception.

    Methods:
        __init__(message, errors=None): Initializes the exception with a message and optional errors.
    """

    def __init__(self, message, errors=None):
        super().__init__(message)
        self.errors = errors
