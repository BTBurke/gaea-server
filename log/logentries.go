package log

import (
	"fmt"
	"os"

	"github.com/bsphere/le_go"
)

func Debug(format string, args ...interface{}) {
	format = "level=debug " + format
	sendArbLogLine(fmt.Sprintf(format, args...))
}

func Info(format string, args ...interface{}) {
	format = "level=info " + format
	sendArbLogLine(fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	format = "level=warn " + format
	sendArbLogLine(fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	format = "level=error " + format
	sendArbLogLine(fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	format = "level=fatal " + format
	sendArbLogLine(fmt.Sprintf(format, args...))
}

func sendArbLogLine(line string) {
	apiToken := os.Getenv("LE_TOKEN")
	if len(apiToken) == 0 {
		fmt.Println("Warning: LE_TOKEN not set.  Unable to send application logs.")
	} else {
		le, err := le_go.Connect(apiToken)
		defer le.Close()
		if err != nil {
			fmt.Printf("Error connecting to logentries. Msg: %s", err)
		} else {
			fmt.Println(line)
			le.Printf(line)
		}
	}
}
