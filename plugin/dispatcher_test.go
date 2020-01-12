package plugin

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestDispatcher_Init(t *testing.T) {
	d := &Dispatcher{}

	err := d.Init("/home/id/Desktop/mdns/conf.d/")

	if err != nil {
		t.Error(err)
	}

	log.Debug(d)
}
