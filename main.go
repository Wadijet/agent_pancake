package main

import (
	"agent_pancake/app/services"
	"agent_pancake/config"
	"agent_pancake/global"
	"fmt"
	"time"
)

func main() {
	// Đọc dữ liệu từ file .env
	global.GlobalConfig = config.NewConfig()

	// Vòng lặp vô hạn chạy 5 phút một lần
	for {
		// Công việc cần thực hiện
		fmt.Println("Thực hiện công việc tại:", time.Now())

		// Nếu chưa đăng nhập thì đăng nhập
		checkin := services.FolkForm_CheckIn()
		if checkin != nil {
			services.FolkForm_Login()
			services.FolkForm_CheckIn()
		}

		// Cập nhật page access token của tất cả các trang
		err := services.FolkForm_UpdatePagesAccessToken()
		if err != nil {
			fmt.Println("Lỗi khi cập nhật page access token:", err)
		}

		// Lấy tất cả các Conversations
		err = services.FolkForm_UpdateAllConversations()
		if err != nil {
			fmt.Println("Lỗi khi lấy danh sách Conversations:", err)
		}

		// Dừng 5 phút trước khi tiếp tục
		time.Sleep(5 * time.Minute)
	}
}
