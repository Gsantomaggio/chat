using Microsoft.Extensions.Logging;

namespace server.src
{
    /// <summary>
    /// Provides logging functionality with different log levels and colored output.
    /// </summary>
    internal static class Logger
    {
        private static readonly ILogger _logger = LoggerFactory.Create(builder =>
        {
            builder
            .AddConsole()
            .SetMinimumLevel(LogLevel.Debug);
        }).CreateLogger(nameof(ServerTCP));

        public static void LogInformation(string message) => Log(LogLevel.Information, message);
        public static void LogInformation(string message, params object[] args) => Log(LogLevel.Information, message, args);

        public static void LogDebug(string message) => Log(LogLevel.Debug, message);
        public static void LogDebug(string message, params object[] args) => Log(LogLevel.Debug, message, args);

        public static void LogWarning(string message) => Log(LogLevel.Warning, message);
        public static void LogWarning(string message, params object[] args) => Log(LogLevel.Warning, message, args);

        public static void LogError(string message) => Log(LogLevel.Error, message);
        public static void LogError(string message, params object[] args) => Log(LogLevel.Error, message, args);

        public static void LogCritical(string message) => Log(LogLevel.Critical, message);
        public static void LogCritical(string message, params object[] args) => Log(LogLevel.Critical, message, args);

        private static void Log(LogLevel level, string message)
        {
            string messageToPrint = $"{Colors.ColorMap[level]}{message}{Colors.Reset}";
            _logger.Log(level, messageToPrint);
        }
        private static void Log(LogLevel level, string message, params object[] args)
        {
            string messageToPrint = $"{Colors.ColorMap[level]}{message}{Colors.Reset}";
            _logger.Log(level, messageToPrint, args);
        }
    }

    /// <summary>
    /// Contains color codes for different log levels to be used in console output.
    /// </summary>
    internal static class Colors
    {
        public const string Reset = "\u001b[0m";
        public static readonly Dictionary<LogLevel, string> ColorMap = new()
        {
            { LogLevel.Debug, "\u001b[36m" },
            { LogLevel.Information, "\u001b[32m" },
            { LogLevel.Warning, "\u001b[33m" },
            { LogLevel.Error, "\u001b[31m" },
            { LogLevel.Critical, "\u001b[1;31m" }
        };
    }
}
