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

	"github.com/gookit/color"
)

func AllPortScanInit() {
	// var userThreads string
	color.Blue.Println("本程序用来扫描IP列表的http服务")
	var fPath string
	fmt.Println("请将文件拖入:")
	fmt.Scanln(&fPath)
	fHandle, err := os.Open(fPath)
	if err != nil {
		panic(fmt.Sprintf("打开文件失败:%s", err.Error()))
	}

	parameterinit.AllPara.Wg.Add(3)
	go ReadIP(fHandle)
	go GetAllPort()
	go process.Process()
	parameterinit.AllPara.Wg.Wait()
}

// 向通道写入全端口目标
func GetAllPort() {
	defer parameterinit.AllPara.Wg.Done()
	for ipAndDmain := range parameterinit.AllPara.IpAndDomainChan {
		for port := 1; port <= 65535; port++ {
			parameterinit.AllPara.Target <- fmt.Sprintf("http://%s:%d", ipAndDmain, port)
		}
	}
	close(parameterinit.AllPara.Target)
}

// 向IP通道写入IP和域名
func ReadIP(fHandle *os.File) {
	defer parameterinit.AllPara.Wg.Done()
	var tmp map[string]interface{} = make(map[string]interface{})

	fByte, _ := ioutil.ReadAll(fHandle)
	fByte = bytes.ReplaceAll(fByte, []byte("\r"), []byte(""))
	// ip正则
	re := `\d+\.\d+\.\d+\.\d+`
	// 域名正则
	re1 := `([\w]([\w]{0,63}[\w])?\.)+[a-zA-Z]{2,6}`
	reg := regexp.MustCompile(re)
	reg1 := regexp.MustCompile(re1)
	// ip去重以及合法性检测
	for _, ip := range reg.FindAllString(string(fByte), -1) {
		address := net.ParseIP(ip)
		if address != nil {
			tmp[address.String()] = nil
		}
	}
	// 域名去重
	for _, domain := range reg1.FindAllString(string(fByte), -1) {
		tmp[domain] = nil
	}
	// 写入通道
	for ip := range tmp {
		parameterinit.AllPara.IpAndDomainChan <- ip
	}
	close(parameterinit.AllPara.IpAndDomainChan)
}
