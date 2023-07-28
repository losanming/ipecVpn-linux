package test

import (
	"danfwing.com/m/zhansheng/utils"
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestCheckMac(t *testing.T) {
	t.Run("get mac addr", func(t *testing.T) {
		utils.GetMacAddr()
	})
	t.Run("check get powershell certs's list and remove certs", func(t *testing.T) {
		// 定义 PowerShell 命令
		cmd := exec.Command("powershell", "-Command", "Get-ChildItem -Path Cert:\\LocalMachine\\Root")

		// 执行命令并获取输出
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("执行 PowerShell 命令时发生错误:", err)
			return
		}
		if string(output) == "" {
			fmt.Println("执行 PowerShell 命令时发生错误:", err)
			return
		}
		//数据处理
		split := strings.Split(string(output), "\n")
		for i, v := range split {
			if strings.Contains(v, "CN=strongSwan Root CA") {
				tmp := strings.Split(split[i], "  ")
				if tmp[0] != "" {
					//删除证书
					exec.Command("powershell", "-Command", "Remove-Item -Path Cert:\\LocalMachine\\Root\\"+tmp[0])
					//删除之后的校验需要做
				}
			}
		}

	})

	t.Run("import certs ", func(t *testing.T) {
		// Import-Certificate -FilePath "C:\Users\Administrator.DESKTOP-S7F2FKM\Desktop\ipsec\strongswanCert.pem" -CertStoreLocation cert:\CurrentUser\Root
		exec.Command("powershell", "-Command", "Import-Certificate -FilePath \"C:\\Users\\Administrator.DESKTOP-S7F2FKM\\Desktop\\ipsec\\strongswanCert.pem\" -CertStoreLocation cert:\\CurrentUser\\Root")
		fmt.Println("over ")
	})
}
