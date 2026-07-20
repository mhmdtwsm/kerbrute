package util

import (
	"os"
	"sync"

	"github.com/op/go-logging"
)

type Logger struct {
	Log       *logging.Logger
	cleanFile *os.File
	mu        sync.Mutex
}

func NewLogger(verbose bool, logFileName string, outputMode string) Logger {
	log := logging.MustGetLogger("kerbrute")
	format := logging.MustStringFormatter(
		`%{color}%{time:2006/01/02 15:04:05} >  %{message}%{color:reset}`,
	)
	formatNoColor := logging.MustStringFormatter(
		`%{time:2006/01/02 15:04:05} >  %{message}`,
	)
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	var cleanFile *os.File

	if logFileName != "" {
		outputFile, err := os.Create(logFileName)
		if err != nil {
			panic(err)
		}

		if outputMode == "default" {
			// Standard go-logging behavior for default mode
			fileBackend := logging.NewLogBackend(outputFile, "", 0)
			fileFormatter := logging.NewBackendFormatter(fileBackend, formatNoColor)
			logging.SetBackend(backendFormatter, fileFormatter)
		} else {
			// Bypass go-logging for the file if we want clean mode
			logging.SetBackend(backendFormatter)
			cleanFile = outputFile
		}
	} else {
		logging.SetBackend(backendFormatter)
	}

	if verbose {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}
	return Logger{Log: log, cleanFile: cleanFile}
}

// WriteClean writes a raw string to the file cleanly and thread-safely
func (l *Logger) WriteClean(msg string) {
	if l.cleanFile != nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.cleanFile.WriteString(msg + "\n")
	}
}