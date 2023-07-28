package InfoService

import "context"

var info_proto = InfoService{}

type InfoService struct {
}

func (i InfoService) mustEmbedUnimplementedInfoServer() {
	//TODO implement me
	panic("implement me")
}

func (i InfoService) GetVersion(ctx context.Context, request *GetVersionRequest) (*GetVersionResponse, error) {
	var response = &GetVersionResponse{}
	response.Version = "1.0.1"
	return response, nil
}
