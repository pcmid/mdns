package config

import "testing"

func TestNewConfig(t *testing.T) {
	c := NewConfig("../../config.sample.d/config.json")

	if c.PluginConfDir != "config.sample.d/" {
		t.Fail()
	}
}
