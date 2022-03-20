package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Logger is a type that defines the logger
type Logger struct {
	mu   sync.Mutex
	Logs *LogContainer
	flag string // Can define that state of the logger wheather to log or not
}

// NewLogger is a function that returns a new logger
func NewLogger(logContainer *LogContainer, flag string) *Logger {
	logger := &Logger{Logs: logContainer}
	logger.SetFlag(flag)
	return logger
}

// SetFlag is a method that set the logger's flag to given state
func (l *Logger) SetFlag(state string) {
	if state != Debug && state != Normal && state != None {
		state = None
	}

	l.flag = state
}

// Log is a method that will log the given statement to the selected log file
func (l *Logger) Log(stmt, logFile string) {

	// Checking the status of the logger
	if l.flag == None {
		return
	} else if l.flag == Debug {
		l.LogToParent(stmt)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	var isValidLogFile bool

	// Checking the validity of the given log file
	validLogFiles := []string{l.Logs.ServerLogFile, l.Logs.BotLogFile, l.Logs.ErrorLogFile}

	for _, validLogFile := range validLogFiles {
		if validLogFile == logFile {
			isValidLogFile = true
		}
	}

	if !isValidLogFile {
		logFile = l.Logs.ServerLogFile
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()

		stmt = fmt.Sprintf("[ %s ] %s", time.Now(), stmt)

		fmt.Fprintln(file, stmt)

	}
}

// LogToParent is a method that will log the given statement to the program starter
func (l *Logger) LogToParent(stmt string) {

	// Checking the status of the logger
	if l.flag == None {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	stmt = fmt.Sprintf("[ %s ] %s", time.Now(), stmt)
	fmt.Println(stmt)
}

// LogToErrorFile is a method that will log the given statement as an error to the error log file
func (l *Logger) LogToErrorFile(stmt string) {

	// Checking the status of the logger
	if l.flag == None {
		return
	} else if l.flag == Debug {
		l.LogToParent(stmt)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	file, err := os.OpenFile(l.Logs.ErrorLogFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()

		stmt = fmt.Sprintf("[ %s ] Error: %s", time.Now(), stmt)

		fmt.Fprintln(file, stmt)

	}
}

// LogToArchiveFile is a method that will log the given statement to archive log file
func (l *Logger) LogToArchiveFile(stmt string) {

	// Checking the status of the logger
	if l.flag == None {
		return
	} else if l.flag == Debug {
		l.LogToParent(stmt)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	file, err := os.OpenFile(l.Logs.ArchiveLogFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()

		stmt = fmt.Sprintf("[ %s ] %s", time.Now(), stmt)

		fmt.Fprintln(file, stmt)

	}
}
