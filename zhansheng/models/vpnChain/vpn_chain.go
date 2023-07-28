package vpnChain

import (
	"container/list"
	"time"
)

type VpnNode struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

type Cert struct {
	VpnUuid string `json:"vpnUuid" required:"true"` //vpnNode uuid
	Cert    []byte `json:"cert" required:"true"`
}

type VpnChain struct {
	Uuid       string    `json:"uuid"`
	ChainName  string    `json:"chainName"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	VpnNodes   list.List `json:"vpnNodes"`
}

var UserVpnChain = make(map[string]VpnChain)
var UserVpnNode []VpnNode
