package specifyipandportscan

import (
	"bytes"
	"fmt"
	parameterinit "httpscanner/parameter_init"
	"httpscanner/process"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gookit/color"
)

func SpecifyipandportscanInit() {
	tmp := getURL()
	parameterinit.AllPara.Wg.Add(2)
	go product(tmp)
	go process.Process()
	parameterinit.AllPara.Wg.Wait()
}

// 将xxx.xxx.xxx.xxx:xx变成合法的http://xxx.xx.xxx.xx:xx
func getURL() []string {

	var tmp []string
	var fPath string
	color.Blue.Println("本程序用来扫描<IP:PORT>列表的http服务")
	fmt.Println("请将文件拖入:")
	fmt.Scanln(&fPath)
	fHandle, err := os.Open(fPath)
	if err != nil {
		panic(fmt.Sprintf("打开文件失败:%s", err.Error()))
	}
	fByte, _ := ioutil.ReadAll(fHandle)
	fByte = bytes.ReplaceAll(fByte, []byte("\r"), []byte(""))
	for _, ipPort := range strings.Split(string(fByte), "\n") {
		tmp = append(tmp, "http://"+ipPort)
	}
	// 去重
	var rmTmp map[string]interface{} = make(map[string]interface{})
	for _, t := range tmp {
		rmTmp[t] = nil
	}
	tmp = nil
	for u := range rmTmp {
		tmp = append(tmp, u)
	}
	return tmp
}

func product(targets []string) {
	defer parameterinit.AllPara.Wg.Done()
	defer close(parameterinit.AllPara.Target)
	for _, t := range targets {
		parameterinit.AllPara.Target <- t
	}
}
