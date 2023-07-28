package strongswan

import (
	"github.com/bronze1man/goStrongswanVici"
	"github.com/sirupsen/logrus"
)

var vpn_client *goStrongswanVici.Client

func GetClient() *goStrongswanVici.Client {
	if vpn_client == nil {
		vpn_client = goStrongswanVici.NewClientFromDefaultSocket()
	}
	out, err := vpn_client.Version()
	if err != nil {
		panic("strongswan client is wrong")
	}
	logrus.Infof("strongswan version :%s\n", out.Version)
	return vpn_client
}

func ListConns() (list []goStrongswanVici.VpnConnInfo, err error) {
	list, err = vpn_client.ListAllVpnConnInfo()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func ConnVpnChain() {
	panic("实现我")
}

func KillAllConns() (rs string, err error) {
	list, err := vpn_client.ListAllVpnConnInfo()
	if err != nil {
		return "get conn lists are failed", err
	}
	if len(list) == 0 {
		return "conn lists are empty", nil
	}
	for _, v := range list {
		logrus.Infof("kill connection id %s\n", v.Uniqueid)
		err := vpn_client.Terminate(&goStrongswanVici.TerminateRequest{Ike_id: v.Uniqueid})
		if err != nil {
			logrus.Errorf("kill connection is failed ike id %s\n", v.Uniqueid)
			return "failed", err
		}
	}
	return "all connections are be killed", nil
}

func NewVpnIkeConfig(child_map_name, ike_map_name string) map[string]goStrongswanVici.IKEConf {
	childConfMap := make(map[string]goStrongswanVici.ChildSAConf)
	childSAConf := goStrongswanVici.ChildSAConf{
		Local_ts:      []string{""},
		Remote_ts:     []string{""},
		ESPProposals:  []string{""},
		StartAction:   "",
		CloseAction:   "",
		Mode:          "tunnel",
		ReqID:         "10",
		RekeyTime:     "10m",
		InstallPolicy: "no",
	}
	childConfMap[child_map_name] = childSAConf
	localAuthConf := goStrongswanVici.AuthConf{
		AuthMethod: "pubkey",
		PubKeys:    []string{""},
		ID:         "",
	}
	remoteAuthConf := goStrongswanVici.AuthConf{
		AuthMethod: "pubkey",
		PubKeys:    []string{""},
	}

	ikeConfMap := make(map[string]goStrongswanVici.IKEConf)

	ikeConf := goStrongswanVici.IKEConf{
		LocalAddrs:  []string{""},
		RemoteAddrs: []string{""},
		Proposals:   []string{""},
		Version:     "",
		LocalAuth:   localAuthConf,
		RemoteAuth:  remoteAuthConf,
		Children:    childConfMap,
		Encap:       "no",
	}

	ikeConfMap[ike_map_name] = ikeConf
	return ikeConfMap
}

func test() {

}
