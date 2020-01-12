package main

import (
	"flag"
	"github.com/pcmid/mdns/core"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

//  go build -ldflags "-X main.version=version"
var version = "0.2.0"

func main() {

	var (
		configPath    string
		logPath       string
		isLogVerbose  bool
		isShowVersion bool
	)

	flag.StringVar(&configPath, "c", "./config.json", "config file path")
	flag.StringVar(&logPath, "l", "", "log file path")
	flag.BoolVar(&isLogVerbose, "v", false, "verbose mode")
	flag.BoolVar(&isShowVersion, "V", false, "current version of overture")
	flag.Parse()

	if isShowVersion {
		println(version)
		return
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if isLogVerbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if logPath != "" {
		lf, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
		if err != nil {
			println("Logfile error: Please check your log file path")
		} else {
			log.SetOutput(io.MultiWriter(lf, os.Stdout))
		}
	}

	core.InitServer(configPath)
}
