package test

import (
	"context"
	"danfwing.com/m/zhansheng/models/vpnChain"
	"danfwing.com/m/zhansheng/service/UserService"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io/ioutil"
	"testing"
)

func TestName(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:17660", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	client := UserService.NewUserServiceClient(conn)
	if client == nil {
		t.Error("client is nil")
		return
	}
	ctx, _ := context.WithCancel(context.Background())
	defer conn.Close()
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
	t.Run("user login", func(t *testing.T) {
		client.UserLogin(ctx, &UserService.UserLoginRequest{
			Username: "admin",
			Password: "admin123456",
		})
	})
	t.Run("get vpn list ", func(t *testing.T) {
		rs, err2 := client.GetVpnNode(ctx, &UserService.GetVpnNodeRequest{})
		if err2 != nil {
			logrus.Error("err: ", err2)
			return
		}
		logrus.Println("rs len : ", len(rs.VpnNodes), "rs msg", rs.Msg)
	})

	t.Run("create chain", func(t *testing.T) {
		rs, err2 := client.VpnChainList(ctx, &UserService.VpnChainListRequest{
			ChainName:      "danfwing-test",
			VpnChainIdList: []string{"cf7ca09d-2d18-11ee-b30f-080027dbb8e1", "cf7ca0e0-2d18-11ee-b30f-080027dbb8e1", "cf7ca132-2d18-11ee-b30f-080027dbb8e1"},
		})
		if err2 != nil {
			logrus.Error("err: ", err2)
			return
		}
		logrus.Println("rs status : ", rs.Status, "rs msg", rs.Msg)

	})
}
