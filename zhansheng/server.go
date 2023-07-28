package main

import (
	"danfwing.com/m/zhansheng/config"
	"danfwing.com/m/zhansheng/config/global"
	"danfwing.com/m/zhansheng/models"
	"danfwing.com/m/zhansheng/service"
	"danfwing.com/m/zhansheng/utils"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"path/filepath"
)

var (
	version  string
	gport    string
	tlsmodel string
	ginport  string
)

func init() {
	flag.StringVar(&gport, "gport", "17660", "grpc端口")
	flag.StringVar(&tlsmodel, "tlsmodel", "insecure", "TLS加密") // mutual 双端加密  insecure 不认证
	flag.StringVar(&ginport, "ginport", "16555", "web服务端口")
	flag.Parse()
	global.Version = version
	global.TLSMODEL = tlsmodel
	//检查目录
	err2 := utils.CheckFile()
	if err2 != nil {
		fmt.Println("err: ")
		os.Exit(1)
	}
	config.GetEnv()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	writer2 := os.Stdout
	if global.Env == "dev" {
		logrus.SetLevel(logrus.InfoLevel)
		global.RsaPrivateKey = global.DevRsaPrivateKey
	} else if global.Env == "prod" {
		logrus.SetLevel(logrus.ErrorLevel)
		global.RsaPrivateKey = global.ProdRsaPrivateKey
	}
	writer3, err := os.OpenFile(fmt.Sprint(filepath.Dir(filepath.Dir(config.Exedir))+"/logs/server.log"), os.O_WRONLY|os.O_CREATE,
		0777)
	if err != nil {
		return
	}
	logrus.SetOutput(io.MultiWriter(writer2, writer3))
	logrus.Infof("init is finish")
}

func main() {
	err := models.InitDB()
	if err != nil {
		panic(err)
	}
	listen, err := net.Listen("tcp", "127.0.0.1:"+gport)
	if err != nil {
		logrus.Errorln("listen tcp is fail ! err : ", err)
		return
	}
	defer listen.Close()
	logrus.Infoln("grpc  success listen port is : ", gport)
	err = service.RegisterServer(listen)
	if err != nil {
		logrus.Errorln("grpc registerServer is fail ! err: ", err)
		panic(err)
	}
}
