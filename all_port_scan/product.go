package allportscan

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/gookit/color"
)

var (
	Threads      int
	IpChan       chan string = make(chan string)
	Target       chan string = make(chan string)
	mu           sync.Mutex
	WG           sync.WaitGroup
	UserAgent    = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36`
	F200, _      = os.OpenFile("200.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	MyHttpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Second * 15,
	}
)

func AllPortScanInit() {
	var userThreads string
	color.Blue.Println("本程序用来扫描http服务")
	fmt.Println("请设置线程数：")
	fmt.Scanln(&userThreads)
	threads, err := strconv.Atoi(userThreads)
	if err != nil {
		color.Red.Println("[-]默认50线程")
		threads = 50
	}
	Threads = threads

	WG.Add(3)
	go ReadIP()
	go AllPortScan()
	go Process()
	WG.Wait()
}

// 向通道写入全端口目标
func AllPortScan() {
	defer WG.Done()
	for port := 1; port <= 65535; port++ {
		for ip := range IpChan {
			Target <- fmt.Sprintf("http://%s:%d", ip, port)
		}
	}
	close(Target)
}

// 向IP通道写入IP
func ReadIP() {
	defer WG.Done()
	var tmp map[string]interface{} = make(map[string]interface{})
	var fPath string
	fmt.Println("请将文件拖入:")
	fmt.Scanln(&fPath)
	fHandle, err := os.Open(fPath)
	if err != nil {
		panic(fmt.Sprintf("打开文件失败:%s", err.Error()))
	}
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
		IpChan <- ip
	}
	close(IpChan)
}
