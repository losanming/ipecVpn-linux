package UserService

import (
	"context"
	"danfwing.com/m/zhansheng/models/vpnChain"
	"danfwing.com/m/zhansheng/utils"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"time"
)

var user_proto = UserService{}

type UserService struct {
}

func (u UserService) GetVpnNode(ctx context.Context, request *GetVpnNodeRequest) (*GetVpnNodeResponse, error) {
	//本地数据测试
	file, err := ioutil.ReadFile("D:\\workAll\\workspace\\ipecVpn-linux-main\\zhansheng\\test\\data\\vpnNode.json")
	if err != nil {
		return nil, status.Error(codes.Internal, "read file error")
	}
	var tmp_data = struct {
		Data []vpnChain.VpnNode `json:"data"`
	}{}
	json.Unmarshal(file, &tmp_data)
	if tmp_data.Data == nil {
		return nil, status.Error(codes.Internal, "read file error")
	}
	//保存到内存 并返回rs
	var rs = GetVpnNodeResponse{}
	rs.VpnNodes = make([]*VpnNode, 0)
	for _, v := range tmp_data.Data {
		vpnChain.UserVpnNode[v.Uuid] = v
		rs.VpnNodes = append(rs.VpnNodes, &VpnNode{
			Uuid: v.Uuid,
			Name: v.Name,
			Ip:   v.Ip,
		})
	}
	rs.Msg = "success"
	rs.Status = 0
	return &rs, nil
}

func (u UserService) VpnChainList(ctx context.Context, request *VpnChainListRequest) (*VpnChainListResponse, error) {
	if request.VpnChainIdList == nil {
		return nil, status.Error(codes.InvalidArgument, "params is wrong")
	}
	if request.ChainName == "" {
		return nil, status.Error(codes.InvalidArgument, "params is wrong")
	}

	//uuid的存在性校验
	for _, v := range request.VpnChainIdList {
		if _, ok := vpnChain.UserVpnNode[v]; !ok {
			return nil, status.Error(codes.InvalidArgument, "params is wrong")
		}
	}
	//创建节点模式链路
	var u_chain vpnChain.VpnChain
	for i, _ := range request.VpnChainIdList {
		u_chain.VpnNodes.PushBack(request.VpnChainIdList[i])
	}
	u_chain.CreateTime = time.Now()
	u_chain.UpdateTime = time.Now()
	u_chain.Uuid = uuid.New().String()
	u_chain.ChainName = request.ChainName

	//保存到内存
	vpnChain.UserVpnChain["user"] = u_chain
	//@@TODO 同步到云管

	var rs = VpnChainListResponse{}
	rs.Msg = "ok"
	rs.Status = 0
	return &rs, nil
}

func (u UserService) UserLogin(ctx context.Context, request *UserLoginRequest) (*UserLoginResponse, error) {
	if request.Username == "admin" && request.Password == "admin123456" {
		return &UserLoginResponse{Status: 0, Msg: "login ok"}, nil
	}
	//@@TODO 云管校验
	rs, err := utils.GetMacAddr()
	if err != nil {
		return nil, status.Error(codes.Internal, "get mac addr error")
	}
	logrus.Println("mac addr is ", rs)

	return &UserLoginResponse{Status: 1, Msg: "login failed"}, status.Error(codes.InvalidArgument, "params is wrong")
}

func (u UserService) mustEmbedUnimplementedUserServiceServer() {
	return
}
