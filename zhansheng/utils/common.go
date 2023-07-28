package utils

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	url2 "net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type CustomError struct {
	Message string
}

func (c CustomError) Error() string {
	return c.Message
}

func StructToMap(p interface{}) (m map[string]string) {
	j, _ := json.Marshal(p)
	_ = json.Unmarshal(j, &m)
	return
}

func PrintFileLog(filename string, v ...interface{}) {
	Exedir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file, err := os.OpenFile(fmt.Sprintf(filepath.Dir(filepath.Dir(Exedir)))+"/logs/"+filename+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		log.Println("log open filename:", filename, ", error:", err.Error())
		return
	}
	logger := log.New(file, "", log.LstdFlags)
	log.Println(fmt.Sprint(v...))
	logger.Println(fmt.Sprint(v...))
	file.Close()
}

func GetInfo(skip int) (info string) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("runtime.Caller() failed")
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file) // Base函数返回路径的最后一个元素
	info = funcName + " " + fileName + " " + strconv.Itoa(lineNo)
	return
}

func Socket5Proxy(url string) (rs []byte, err error) {
	dialer, err := proxy.SOCKS5("tcp", url, nil, proxy.Direct)
	if err != nil {
		return rs, err
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport, Timeout: 5 * time.Second}
	httpTransport.Dial = dialer.Dial
	if resp, err := httpClient.Get("https://httpbin.org/ip"); err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return body, err
	}
	return rs, nil
}

func HttpProxy(url string) (rs []byte, err error) {
	//获取ip
	getUrl := "https://httpbin.org/ip"
	u, _ := url2.Parse(url)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(u),
		},
	}
	response, err := client.Get(getUrl)
	if err != nil {
		return rs, err
	}
	defer response.Body.Close()

	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return rs, err
	}
	return res, err
}

func CheckFile() error {
	Exedir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	lastPath := filepath.Dir(filepath.Dir(Exedir))
	// log
	_, err := os.Stat(fmt.Sprintf(lastPath + "/logs"))
	if os.IsNotExist(err) {
		err := os.MkdirAll(fmt.Sprintf(lastPath+"/logs"), 0766)
		if err != nil {
			return err
		}
	}
	return nil
}

// 通过map主键唯一的特性过滤重复元素
func RemoveRepByMap(slc []int) []int {
	result := []int{}
	tempMap := map[int]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

// 通过两重循环过滤重复元素
func RemoveRepByLoop(slc []int) []int {
	result := []int{}
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, slc[i])
		}
	}
	return result
}

// 元素去重
func RemoveRep(slc []int) []int {
	if len(slc) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return RemoveRepByLoop(slc)
	} else {
		// 大于的时候，通过map来过滤
		return RemoveRepByMap(slc)
	}
}

// 切片实现分页
func SlicePage(page, pageSize, nums int64) (sliceStart, sliceEnd int64) {
	if page <= 0 {
		page = 1
	}
	if pageSize < 0 {
		pageSize = 20 //设置一页默认显示的记录数
	}
	if pageSize > nums {
		return 0, nums
	}
	// 总页数
	pageCount := int64(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}

func GetMacAddr() (rs []string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return rs, err
	}

	for _, inter := range interfaces {
		if inter.HardwareAddr != "" {
			rs = append(rs, inter.HardwareAddr)
		}
	}

	return rs, nil
}
