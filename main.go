package main

import (
	"agent_pancake/config"
	"agent_pancake/global"
)

func main() {
	// Đọc dữ liệu từ file .env
	global.GlobalConfig = config.NewConfig()

	// Tạo vòng lặp vô hạn chạy 5 phút một lần

}
