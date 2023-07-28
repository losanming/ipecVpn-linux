package config

import (
	"danfwing.com/m/zhansheng/config/global"
	"danfwing.com/m/zhansheng/utils"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
)

var goos, exepath, Exedir, runpath string

// 获取可用端口
func GetAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil

}

func GetEnv() {
	// 输出环境日志
	goos = runtime.GOOS
	exepath, _ = filepath.Abs(os.Args[0])
	Exedir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	runpath, _ = os.Getwd()
	utils.PrintFileLog(global.INFO, "操作系统：", goos)
	utils.PrintFileLog(global.INFO, "执行程序路径：", exepath)
	utils.PrintFileLog(global.INFO, "执行程序文件路径：", Exedir)
	utils.PrintFileLog(global.INFO, "运行目录路径：", runpath)
	utils.PrintFileLog(global.INFO, "Version ：", global.Version)
	utils.PrintFileLog(global.INFO, "TLSMODULE ：", global.TLSMODEL)
}
