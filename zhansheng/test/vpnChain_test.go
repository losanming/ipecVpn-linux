package test

import (
	"danfwing.com/m/zhansheng/models/vpnChain"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"testing"
)

func TestName(t *testing.T) {
	t.Run("create data", func(t *testing.T) {
		type TmpData struct {
			Data []vpnChain.VpnNode `json:"data"`
		}
		var rs TmpData
		for i := 0; i < 100; i++ {
			t_uuid, _ := uuid.NewUUID()
			var tmp = vpnChain.VpnNode{
				Ip:   fmt.Sprintf("192.168.1.%d", i),
				Name: fmt.Sprintf("node%d", i),
				Uuid: t_uuid.String(),
			}
			rs.Data = append(rs.Data, tmp)
		}
		b, err := json.Marshal(rs)
		if err != nil {
			return
		}
		ioutil.WriteFile("/home/danfwing/Ipsec-linux/Ipsec-linux/path/zhansheng/test/data/vpnNode.json", b, 0666)
	})
	t.Run("set vpn ", func(t *testing.T) {

	})
}