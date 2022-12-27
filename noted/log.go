package noted

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitializeLogger ensures that logging is set up for the project.
// currently supports console logging only.
func InitializeLogger() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return "[" + file + ":" + strconv.Itoa(line) + "]"
		// } else {
		// 	return "[" + file + "]"
		// }
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Caller().Logger()
}