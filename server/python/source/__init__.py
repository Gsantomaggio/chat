import logging

"""
This module sets up a custom logging system with color-coded log levels.

Classes:
    CustomFormatter: A custom logging formatter that adds color to log messages based on their severity level.
    Logger: A logger class that configures and returns a logger with the custom formatter.

Usage:
    The module can be run as a standalone script to demonstrate the logging functionality. When run, it will log messages of various severity levels (DEBUG, INFO, WARNING, ERROR, CRITICAL) to the console with appropriate colors.

Example:
    logger = Logger(__name__)
    logger.debug("This is a debug message")
    logger.info("This is an info message")
    logger.warning("This is a warning message")
    logger.error("This is an error message")
    logger.critical("This is a critical message")
"""


RESET = "\033[0m"
COLORS = {
    "DEBUG": "\033[36m",  # Cyan
    "INFO": "\033[32m",  # Green
    "WARNING": "\033[33m",  # Yellow
    "ERROR": "\033[31m",  # Red
    "CRITICAL": "\033[1;31m",  # Bold Red
}


class CustomFormatter(logging.Formatter):
    def format(self, record):
        log_color = COLORS.get(record.levelname, RESET)
        message = super().format(record)
        return f"{log_color}{message}{RESET}"


class Logger:
    def __new__(cls, module):
        logger = logging.getLogger(module)
        logger.setLevel(logging.DEBUG)
        cls._configure_handler(logger)
        return logger

    @staticmethod
    def _configure_handler(logger):
        formatter = CustomFormatter(
            fmt="%(asctime)s - %(message)s", datefmt="%Y-%m-%d %H:%M:%S"
        )
        console_handler = logging.StreamHandler()
        console_handler.setLevel(logging.DEBUG)
        console_handler.setFormatter(formatter)
        logger.addHandler(console_handler)


if __name__ == "__main__":
    logger = Logger(__name__)
    logger.debug("This is a debug message")
    logger.info("This is an info message")
    logger.warning("This is a warning message")
    logger.error("This is an error message")
    logger.critical("This is a critical message")
