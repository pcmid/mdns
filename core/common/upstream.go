package common

type DNSUpstream struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	Protocol      string `json:"protocol"`
	SOCKS5Address string `json:"socks5_address"`
	Timeout       int    `json:"timeout"`
}
