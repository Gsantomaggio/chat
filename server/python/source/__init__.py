import logging


RESET = "\033[0m"
COLORS = {
    'DEBUG': "\033[36m",  # Cyan
    'INFO': "\033[32m",   # Green
    'WARNING': "\033[33m",  # Yellow
    'ERROR': "\033[31m",  # Red
    'CRITICAL': "\033[1;31m",  # Bold Red
}

class CustomFormatter(logging.Formatter):
    
    def format(self, record):
        log_color = COLORS.get(record.levelname, RESET)
        record.msg = f"{log_color}{record.msg}{RESET}"
        return super().format(record)


class Logger:
    def __new__(cls, module):
        logger = logging.getLogger(module)
        logger.setLevel(logging.DEBUG)
        cls._configure_handler(logger)
        return logger

    @staticmethod
    def _configure_handler(logger):
        formatter = CustomFormatter('%(asctime)s - %(message)s')
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