package log

// Debug is a constant that indicates the logger is in debug mode
const Debug = "Debug"

// Normal is a constant that indicates the logger is normally logging in the log files
const Normal = "Normal"

// None is a constant that indicates the logger isn't logging
const None = "None"

// ILogger is an interface that defines the logging style
type ILogger interface {
	SetFlag(state string)
	Log(stmt, logFile string)
	LogToParent(stmt string)
	LogToErrorFile(stmt string)
	LogToArchiveFile(stmt string)
}

// LogContainer is a type that defines all the available logs
type LogContainer struct {
	UserLogFile             string
	ServiceProviderLogFile  string
	ProjectLogFile          string
	SubscriptionPlanLogFile string
	SubscriptionLogFile     string
	TransactionLogFile      string
	DeletedLogFile          string
	ServerLogFile           string
	BotLogFile              string
	ErrorLogFile            string
	ArchiveLogFile          string
}
