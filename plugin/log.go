package plugin

import (
	"github.com/pcmid/mdns/core/common"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	Register(&Logger{})
}

type Logger struct {
	LogFile string `json:"log_file"`
	logger  *logrus.Logger
}

func (l *Logger) Name() string {
	return "log"
}

func (l *Logger) Init(config map[string]interface{}) error {
	l.logger = logrus.New()

	l.LogFile = config["log_file"].(string)

	if l.LogFile != "" {
		logFile, err := os.OpenFile(l.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			// open log file failed, log to stdout
			logrus.Errorf("Can not open log file: %s. Set to stdout", err)
			return nil
		}

		l.logger.SetOutput(logFile)
	}

	return nil
}

func (l *Logger) HandleDns(ctx *common.Context) {
	l.logger.Infof("Question from %s %s", ctx.Client.RemoteAddr(), ctx.Query.Question[0].String())
}

func (l *Logger) Where() uint8 {
	return IN
}
