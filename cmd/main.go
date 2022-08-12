package main

import (
	"fmt"
	allportscan "httpscanner/all_port_scan"
	specifyipandportscan "httpscanner/specify_ipandport_scan"

	"github.com/gookit/color"
)

func main() {
	var userchoose string
	color.Green.Println("请选择模式：")
	color.Green.Println("1.输入<ip>合集文件，进行全端口扫描，文件中地址使用换行分隔。")
	color.Green.Println("1.输入<ip:port>合集文件，进行全端口扫描，文件中地址使用换行分隔。")
	fmt.Scanln(&userchoose)

	switch userchoose {
	case "1":
		allportscan.AllPortScanInit()
	case "2":
		specifyipandportscan.SpecifyipandportscanInit()
	default:
		color.Red.Println("输入错误，程序退出。")
	}

}
