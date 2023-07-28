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
	VpnNodes   list.List `json:"vpnNodes"` //node uuid
}

var UserVpnChain = make(map[string]VpnChain, 2) // 一共两个key user和auto
var UserVpnNode = make(map[string]VpnNode)      // key = uuid
