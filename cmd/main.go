package main

import (
	"fmt"
	allportscan "httpscanner/all_port_scan"
	specifyipandportscan "httpscanner/specify_ipandport_scan"
	"os"

	"github.com/gookit/color"
)

func main() {
	var userchoose string
	color.Green.Println("请选择模式：")
	color.Green.Println("1.输入<ip_or_domain>合集文件，进行全端口扫描，文件中地址使用换行分隔。")
	color.Green.Println("2.输入<ip_or_domain:port>合集文件，进行指定端口扫描，文件中地址使用换行分隔。")
	fmt.Scanln(&userchoose)

	switch userchoose {
	case "1":
		allportscan.AllPortScanInit()
	case "2":
		specifyipandportscan.SpecifyipandportscanInit()
	default:
		color.Red.Println("输入错误，程序退出。")
	}
	color.Green.Println("\n[+]完成！结果保存在./200.txt中")

	fmt.Println()
	fmt.Println("按任意键退出...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
