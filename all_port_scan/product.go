package allportscan

import (
	"bytes"
	"fmt"
	parameterinit "httpscanner/parameter_init"
	"httpscanner/process"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"

	"github.com/gookit/color"
)

func AllPortScanInit() {
	var userThreads string
	color.Blue.Println("本程序用来扫描IP列表的http服务")
	var fPath string
	fmt.Println("请将文件拖入:")
	fmt.Scanln(&fPath)
	fHandle, err := os.Open(fPath)
	if err != nil {
		panic(fmt.Sprintf("打开文件失败:%s", err.Error()))
	}

	fmt.Println("请设置线程数：")
	fmt.Scanln(&userThreads)
	threads, err := strconv.Atoi(userThreads)
	if err != nil {
		color.Red.Println("[-]默认50线程")
		threads = 50
	}
	parameterinit.AllPara.Threads = threads

	parameterinit.AllPara.Wg.Add(3)
	go ReadIP(fHandle)
	go GetAllPort()
	go process.Process()
	parameterinit.AllPara.Wg.Wait()
}

// 向通道写入全端口目标
func GetAllPort() {
	defer parameterinit.AllPara.Wg.Done()
	for ip := range parameterinit.AllPara.IpChan {
		for port := 1; port <= 65535; port++ {
			parameterinit.AllPara.Target <- fmt.Sprintf("http://%s:%d", ip, port)
		}
	}
	close(parameterinit.AllPara.Target)
}

// 向IP通道写入IP
func ReadIP(fHandle *os.File) {
	defer parameterinit.AllPara.Wg.Done()
	var tmp map[string]interface{} = make(map[string]interface{})

	fByte, _ := ioutil.ReadAll(fHandle)
	fByte = bytes.ReplaceAll(fByte, []byte("\r"), []byte(""))
	re := `\d+\.\d+\.\d+\.\d+`
	reg := regexp.MustCompile(re)
	// 去重以及合法性检测
	for _, ip := range reg.FindAllString(string(fByte), -1) {
		address := net.ParseIP(ip)
		if address != nil {
			tmp[address.String()] = nil
		}
	}
	for ip := range tmp {
		parameterinit.AllPara.IpChan <- ip
	}
	close(parameterinit.AllPara.IpChan)
}
