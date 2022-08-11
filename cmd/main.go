package main

import allportscan "httpscanner/all_port_scan"

func main() {
	AllPortScan()
}

// 向通道写入全端口目标
func AllPortScan() {
	allportscan.AllPortScanInit()
}
