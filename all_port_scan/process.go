package allportscan

import (
	"fmt"
	"net/http"
	"strconv"
)

// // 将xxx.xxx.xxx.xxx:xx变成合法的http://xxx.xx.xxx.xx:xx
// func GetURL() []string {
// 	var tmp []string
// 	var fPath string
// 	fmt.Println("请将文件拖入:")
// 	fmt.Scanln(&fPath)
// 	fHandle, err := os.Open(fPath)
// 	if err != nil {
// 		panic(fmt.Sprintf("打开文件失败:%s", err.Error()))
// 	}
// 	fByte, _ := ioutil.ReadAll(fHandle)
// 	fByte = bytes.ReplaceAll(fByte, []byte("\r"), []byte(""))
// 	for _, ipPort := range strings.Split(string(fByte), "\n") {
// 		tmp = append(tmp, "http://"+ipPort)
// 	}
// 	return tmp
// }

// 检测模块
func scanHttpTarget(s string) {

	httpNewRequest, _ := http.NewRequest("GET", s, nil)
	httpNewRequest.Header.Add(`user-agent`, UserAgent)
	resp, err := MyHttpClient.Do(httpNewRequest)
	if err != nil {
		fmt.Printf("\r" + err.Error())
		return
	}
	fmt.Printf("\r[*] " + s + "        code:" + strconv.Itoa(resp.StatusCode))
	if resp.StatusCode == 200 {
		mu.Lock()
		F200.WriteString(s + "\n")
		//countsOf200++
		//color.Blue.Printf("\r[+] 网站存活数量：%i", countsOf200)
		mu.Unlock()
		return
	}
}

// 并发验证
func Process() {
	defer WG.Done()
	// 并发控制
	WG.Add(Threads)
	for thread := 0; thread < Threads; thread++ {
		go func(thread int) {
			defer WG.Done()
			defer fmt.Println("[*] gorouting" + strconv.Itoa(thread+1) + "关闭!")
			fmt.Println("[*] gorouting" + strconv.Itoa(thread+1) + "开启!")
			for t := range Target {
				scanHttpTarget(t)
			}
		}(thread)
	}
}
