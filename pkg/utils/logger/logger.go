package logger

import "github.com/rs/zerolog"

type severityGCP string

const (
	GCPAlertLevel    severityGCP = "ALERT"
	GCPCriticalLevel severityGCP = "CRITICAL"
	GCPErrorLevel    severityGCP = "ERROR"
	GCPWarningLevel  severityGCP = "WARNING"
	GCPInfoLevel     severityGCP = "INFO"
	GCPDebugLevel    severityGCP = "DEBUG"

	GCPFieldKeySeverity = "severity"
	GCPFieldKeyMsg      = "message"
	GCPFieldKeyTime     = "time"
)

var (
	zerologToGCPLevel = map[zerolog.Level]severityGCP{
		zerolog.PanicLevel: GCPAlertLevel,
		zerolog.FatalLevel: GCPCriticalLevel,
		zerolog.ErrorLevel: GCPErrorLevel,
		zerolog.WarnLevel:  GCPWarningLevel,
		zerolog.InfoLevel:  GCPInfoLevel,
		zerolog.DebugLevel: GCPDebugLevel,
		zerolog.TraceLevel: GCPDebugLevel,
	}
)

// SeverityHook adds GCP severity levels to zerolog output log messages.
type SeverityHook struct{}

// Run implements the zerolog.Hook interface.
func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		e.Str(GCPFieldKeySeverity, string(zerologToGCPLevel[level]))
	}
}
