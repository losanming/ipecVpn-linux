package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"danfwing.com/m/zhansheng/config/global"
	"danfwing.com/m/zhansheng/service/InfoService"
	"danfwing.com/m/zhansheng/service/UserService"
	"danfwing.com/m/zhansheng/service/configService"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"net"
)

func RegisterServer(lis net.Listener) (err error) {
	//TLS 双端加密
	//从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	if global.TLSMODEL == "mutual" {
		certificate, err := tls.X509KeyPair([]byte(global.TLSSERVERPEM), []byte(global.TLSSERVERKEY))
		if err != nil {
			logrus.Println("err : ", err)
		}

		// 创建CertPool，后续就用池里的证书来校验客户端证书有效性
		// 所以如果有多个客户端 可以给每个客户端使用不同的 CA 证书，来实现分别校验的目的
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM([]byte(global.TLSCAPEM)); !ok {
			logrus.Println("failed to append certs")
		}

		// 构建基于 TLS 的 TransportCredentials
		creds := credentials.NewTLS(&tls.Config{
			// 设置证书链，允许包含一个或多个
			Certificates: []tls.Certificate{certificate},
			// 要求必须校验客户端的证书 可以根据实际情况选用其他参数
			ClientAuth: tls.RequireAndVerifyClientCert, // NOTE: this is optional!
			// 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
			ClientCAs: certPool,
		})

		if creds != nil {
		}
		//TODO 数据验证
		gs := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(LoggingInterceptor))
		configService.RegisterConfigServer(gs, &configService.SystemConfigService{})
		InfoService.RegisterInfoServer(gs, &InfoService.InfoService{})
		UserService.RegisterUserServiceServer(gs, &UserService.UserService{})
		err = gs.Serve(lis)
		if err != nil {
			logrus.Println(" grpc is err ", err)
		}
		return err
	} else if global.TLSMODEL == "server-side" {
		certificate, err := tls.X509KeyPair([]byte(global.TLSSERVERPEM), []byte(global.TLSSERVERKEY))
		if err != nil {
			logrus.Println("err : ", err)
		}
		creds := credentials.NewServerTLSFromCert(&certificate)
		gs := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(LoggingInterceptor))
		configService.RegisterConfigServer(gs, &configService.SystemConfigService{})
		InfoService.RegisterInfoServer(gs, &InfoService.InfoService{})
		UserService.RegisterUserServiceServer(gs, &UserService.UserService{})
		err = gs.Serve(lis)
		if err != nil {
			logrus.Println(" grpc is err ", err)
		}
		return err
	} else if global.TLSMODEL == "insecure" {
		gs := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor))
		configService.RegisterConfigServer(gs, &configService.SystemConfigService{})
		InfoService.RegisterInfoServer(gs, &InfoService.InfoService{})
		UserService.RegisterUserServiceServer(gs, &UserService.UserService{})
		reflection.Register(gs)
		err = gs.Serve(lis)
		if err != nil {
			logrus.Println(" grpc is err ", err)
		}
		return err
	}
	return err
}

// 拦截器 - 打印日志
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	logrus.Infof("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	logrus.Infof("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}
