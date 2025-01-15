using Microsoft.Extensions.Logging;

namespace server.src
{
    /// <summary>
    /// Provides logging functionality with different log levels and colored output.
    /// </summary>
    public sealed class ColoredLogger
    {
        private readonly ILogger _logger;
        private static readonly Lazy<ColoredLogger> _instance = new(() => new ColoredLogger());

        public static ColoredLogger Instance => _instance.Value;

        private ColoredLogger()
        {
            _logger = LoggerFactory.Create(builder =>
                builder.AddConsole()
                      .SetMinimumLevel(LogLevel.Debug))
                .CreateLogger(nameof(ServerTCP));
        }

        public void Log(LogLevel level, string message, params object[] args)
        {
            var color = ConsoleColors.GetColor(level);
            var coloredMessageTemplate = $"{color}{message}{ConsoleColors.Reset}";

            if (args.Length > 0)
            {
                _logger.Log(level, coloredMessageTemplate, args);
            }
            else
            {
                _logger.Log(level, coloredMessageTemplate);
            }
        }

        public void LogInformation(string message, params object[] args) =>
            Log(LogLevel.Information, message, args);
        public void LogDebug(string message, params object[] args) =>
            Log(LogLevel.Debug, message, args);
        public void LogWarning(string message, params object[] args) =>
            Log(LogLevel.Warning, message, args);
        public void LogError(string message, params object[] args) =>
            Log(LogLevel.Error, message, args);
        public void LogCritical(string message, params object[] args) =>
            Log(LogLevel.Critical, message, args);
    }

    /// <summary>
    /// Contains color codes for different log levels to be used in console output.
    /// </summary>
    internal static class ConsoleColors
    {
        public const string Reset = "\u001b[0m";

        private static readonly Dictionary<LogLevel, string> _colorMap = new()
        {
            { LogLevel.Debug, "\u001b[36m" },      // Ciano
            { LogLevel.Information, "\u001b[32m" }, // Verde
            { LogLevel.Warning, "\u001b[33m" },     // Giallo
            { LogLevel.Error, "\u001b[31m" },       // Rosso
            { LogLevel.Critical, "\u001b[1;31m" }   // Rosso Brillante
        };

        public static string GetColor(LogLevel level) =>
            _colorMap.TryGetValue(level, out var color) ? color : string.Empty;
    }
}