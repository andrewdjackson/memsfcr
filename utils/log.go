package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

// ECUCommandTrace color coding for Commands sent to the ECU
// const ECUCommandTrace = "\u001b[38;5;200mECU_TX>\u001b[0m"
const ECUCommandTrace = "ECU_TX>"

// ECUResponseTrace color coding for Reponses returned from the ECU
//const ECUResponseTrace = "\u001b[38;5;200mECU_RX<\u001b[0m"
const ECUResponseTrace = "ECU_RX<"

// ReceiveFromWebTrace color coding for Messages sent to the Web
//const ReceiveFromWebTrace = "\u001b[38;5;21mWEB_RX<\u001b[0m"
const ReceiveFromWebTrace = "WEB_RX<"

// SendToWebTrace color coding for Messages received from the Web
const SendToWebTrace = "WEB_TX>"

//const SendToWebTrace = "\u001b[38;5;21mWEB_TX>\u001b[0m"

// DiagnosticTrace color coding for Messages received from the Web
//const DiagnosticTrace = "\u001b[38;5;220mDIAGTR>\u001b[0m"
const DiagnosticTrace = "DIAGTR>"

// EmulatorTrace color coding
//const EmulatorTrace = "\u001b[38;5;106mEMU_TR>\u001b[0m"
const EmulatorTrace = "EMU_TR>"

var console = false

var (
	// LogE logs as an error
	LogCE = log.New(LogWriter{}, "\u001b[38;5;160mERROR: ", 0)
	// LogW logs as a warning
	LogCW = log.New(LogWriter{}, "\u001b[38;5;214mWARNING: ", 0)
	// LogI logs as an info, no prefix
	LogCI = log.New(LogWriter{}, "\u001b[38;5;70m", 0)
	// LogE logs as an error
	LogE = log.New(LogWriter{}, "ERROR: ", 0)
	// LogW logs as a warning
	LogW = log.New(LogWriter{}, "WARNING: ", 0)
	// LogI logs as an info, no prefix
	LogI = log.New(LogWriter{}, "", 0)
)

// LogWriter is used to format the log message
type LogWriter struct{}

// Write the log entry
func (f LogWriter) Write(p []byte) (n int, err error) {
	pc, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	format := "%s\r %s: %d %s"

	if console {
		format = "%s\r \u001b[38;5;38m↵ %s: %d %s\u001b[0m"
	}

	logEntry := fmt.Sprintf(format, p, filepath.Base(file), line, fnName)
	log.Printf("%s", logEntry)

	return len(p), nil
}
