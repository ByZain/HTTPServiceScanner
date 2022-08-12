package parameterinit

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gookit/color"
)

type Parainit struct {
	Threads  int
	Target   chan string
	IpChan   chan string
	Wg       sync.WaitGroup
	Mu       sync.Mutex
	HttpPara HttpPara
	FilePara FilePara
}

type HttpPara struct {
	UserAgent    string
	MyHttpClient *http.Client
}

type FilePara struct {
	F200 *os.File
}

var AllPara Parainit

func init() {
	AllPara = Parainit{
		Threads: 50,
		Target:  make(chan string),
		IpChan:  make(chan string),
		Wg:      sync.WaitGroup{},
		Mu:      sync.Mutex{},
		HttpPara: HttpPara{
			UserAgent: `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36`,
			MyHttpClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
				Timeout: time.Second * 15,
			},
		},
	}
	AllPara.FilePara.F200, _ = os.OpenFile("200.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)

	var userThreads string
	color.Blue.Println("请设置线程:")
	fmt.Scanln(&userThreads)
	if therad, err := strconv.Atoi(userThreads); err == nil {
		AllPara.Threads = therad
	} else {
		color.Red.Println("[-]输入错误，默认使用50线程")
	}

}
