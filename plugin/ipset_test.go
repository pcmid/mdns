package plugin

import "testing"

func TestIpSet_Init(t *testing.T) {
	err := (&IpSet{}).Init("conf.d/")

	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
