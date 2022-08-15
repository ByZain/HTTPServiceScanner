package process

import (
	"fmt"
	parameterinit "httpscanner/parameter_init"
	"net/http"
	"strconv"
)

// 检测模块
func scanHttpTarget(s string) {
	httpNewRequest, _ := http.NewRequest("GET", s, nil)

	// 关闭长连接，及时释放资源,默认是keep-alive
	// 也可以这样关：httpNewRequest.Header.Add("Connection", "close")
	httpNewRequest.Close = true

	httpNewRequest.Header.Add(`user-agent`, parameterinit.AllPara.HttpPara.UserAgent)

	resp, err := parameterinit.AllPara.HttpPara.MyHttpClient.Do(httpNewRequest)
	if err != nil {
		fmt.Printf("\r" + err.Error())
		return
	}
	fmt.Printf("\r[*] " + s + "        code:" + strconv.Itoa(resp.StatusCode))
	if resp.StatusCode == 200 {
		parameterinit.AllPara.Mu.Lock()
		parameterinit.AllPara.FilePara.F200.WriteString(s + "\n")
		parameterinit.AllPara.Mu.Unlock()
		return
	}
}

// 并发验证
func Process() {
	defer parameterinit.AllPara.Wg.Done()
	// 并发控制
	parameterinit.AllPara.Wg.Add(parameterinit.AllPara.Threads)
	for thread := 0; thread < parameterinit.AllPara.Threads; thread++ {
		go func(thread int) {
			defer parameterinit.AllPara.Wg.Done()
			defer fmt.Printf("\r[*] gorouting" + strconv.Itoa(thread+1) + "关闭!")
			fmt.Printf("\r[*] gorouting" + strconv.Itoa(thread+1) + "开启!")
			for t := range parameterinit.AllPara.Target {
				scanHttpTarget(t)
			}
		}(thread)
	}
}
