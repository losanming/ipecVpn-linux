package configService

import (
	"context"
	config2 "danfwing.com/m/zhansheng/models/config"
	"danfwing.com/m/zhansheng/utils/sqlWithTransaction"
	"errors"
	"github.com/jinzhu/gorm"
	"google.golang.org/protobuf/runtime/protoimpl"
	"strings"
	"time"
)

var config_proto = SystemConfigService{}

type SystemConfigService struct{}

func (s SystemConfigService) UpdateConfigById(ctx context.Context, request *UpdateConfigByIdRequest) (*UpdateConfigByIdResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s SystemConfigService) GetInfoByKey(ctx context.Context, request *GetInfoByKeyRequest) (*GetInfoByKeyResponse, error) {
	if request.Key == "" {
		return nil, errors.New("key is empty")
	}
	var config config2.SystemConfig
	config.Key = request.Key
	if err := config.Exit(); gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("key is not found")
	}

	rs, err := config.GetInfoByKey()
	if err != nil {
		return nil, err
	}
	var resp GetInfoByKeyResponse
	var config_list = &ConfigList{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		Id:            int32(rs.ID),
		Key:           rs.Key,
		Name:          rs.Name,
		ValueExplain:  rs.ValueExplain,
		CreatedAt:     rs.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     rs.UpdatedAt.Format("2006-01-02 15:04:05"),
		Value:         rs.Value,
	}
	resp.Info = config_list
	return &resp, err
}

func (s SystemConfigService) CreateConfig(ctx context.Context, request *CreateConfigRequest) (*CreateConfigResponse, error) {
	if request.Key == "" || request.Name == "" {
		return nil, errors.New("param is wrong")
	}
	var config config2.SystemConfig
	config.Key = request.Key
	if err := config.Exit(); !gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("键名称重复")
	}
	//去除空格
	ts := strings.TrimSpace(request.Key)
	//去除中间
	rp := strings.Replace(ts, " ", "", -1)
	if rp == "" {
		return nil, errors.New("key have space,space replace is fail")
	}
	request.Key = rp
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	config.Name = request.Name
	config.ValueExplain = request.ValueExplain
	config.Value = request.Value
	err := sqlWithTransaction.D_Transaction(func(DB **gorm.DB) (err error) {
		_, err = config.Create()
		return err
	})
	if err != nil {
		return nil, errors.New("config create is fail ")
	}
	var resp CreateConfigResponse
	resp.Id = int32(config.ID)
	resp.Message = "ok"
	return &resp, err
}
func (s SystemConfigService) GetConfigList(ctx context.Context, request *GetConfigListRequest) (*GetConfigListResponse, error) {
	if request.PageSize <= 0 || request.Page <= 0 {
		return nil, errors.New("param is wrong")
	}
	var config config2.SystemConfig
	getlist, err := config.Getlist(int(request.Page), int(request.PageSize))
	if err != nil {
		return nil, err
	}
	var resp GetConfigListResponse
	var rs []*ConfigList

	for _, v := range getlist {
		var temp = &ConfigList{
			state:         protoimpl.MessageState{},
			sizeCache:     0,
			unknownFields: nil,
			Id:            int32(v.ID),
			Key:           v.Key,
			Name:          v.Name,
			Value:         v.Value,
			ValueExplain:  v.ValueExplain,
			CreatedAt:     v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		rs = append(rs, temp)
	}
	count, err := config.GetCount()
	if err != nil {
		count = 0
	}
	resp.ConfigList = rs
	resp.Total = int32(count)
	return &resp, err
}

func (s SystemConfigService) DeleteConfigById(ctx context.Context, request *DeleteConfigByIdRequest) (*DeleteConfigByIdResponse, error) {
	//TODO implement me
	if request.Id == 0 {
		return nil, errors.New("id is wrong")
	}
	var config config2.SystemConfig
	config.ID = uint(request.Id)
	value, err2 := config.GetInfoById()
	if err2 != nil {
		return nil, errors.New("get config is fail")
	}
	if value.Key == "frpc" {
		return nil, errors.New("frpc can't be deleted")
	}
	if err := config.Exit(); gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("id is wrong")
	}
	err := sqlWithTransaction.D_Transaction(func(DB **gorm.DB) (err error) {
		_, err = config.Delete()
		return err
	})
	if err != nil {
		return nil, err
	}
	var resp DeleteConfigByIdResponse
	resp.Id = request.Id
	resp.Message = "delete ok"
	return &resp, err
}

func (s SystemConfigService) mustEmbedUnimplementedConfigServer() {
	//TODO implement me
}
