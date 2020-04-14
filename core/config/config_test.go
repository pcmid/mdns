package config

import "testing"

func TestNewConfig(t *testing.T) {
	c := NewConfig("../../config.sample.d/config.json")

	print(c)
}
