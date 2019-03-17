package nipplebot

import (
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"os"
)

var log = &logrus.Logger{
	Out: os.Stderr,
	Formatter: &logrus.TextFormatter{
		//TimestampFormat: "2018-10-03 20:13:37",
		FullTimestamp: true,
	},
	Hooks: make(logrus.LevelHooks),
	Level: logrus.DebugLevel,
}

// CheckErr is a generic function to validate all errors.
func CheckErr(err error) {
	if err != nil {
		log.Fatalf("error: %s\n --- exiting ---", err)
		os.Exit(1337)
	}
}

// getMaxTargetCount return an "uncertain" maxTargetCount (try to make it more human)
func makeUncertain(inp int) int {
	return inp + funk.RandomInt(1, inp)
}
