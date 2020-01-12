package plugin

import (
	"encoding/json"
	"github.com/pcmid/mdns/core/common"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type Logger struct {
	Level   string `json:"level"`
	LogFile string `json:"log_file"`

	logger *logrus.Logger
}

type LoggerConfig struct {
	Level string `json:"level"`
}

func init() {
	Register(&Logger{})
}

func (l *Logger) Name() string {
	return "log"
}

func (l *Logger) Init(configDir string) error {
	l.logger = logrus.New()

	jsonData, _ := ioutil.ReadFile(configDir + "log.json")
	config := LoggerConfig{}

	err := json.Unmarshal(jsonData, &config)

	if err != nil {
		return err
	}

	l.logger.Level, _ = logrus.ParseLevel(config.Level)
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
