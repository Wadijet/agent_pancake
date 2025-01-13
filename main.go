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
		_, err := services.FolkForm_CheckIn()
		if err != nil {
			services.FolkForm_Login()
			services.FolkForm_CheckIn()
		}

		err = services.Bridge_SyncPages()
		if err != nil {
			fmt.Println("Lỗi khi đồng bộ trang:", err)
		}

		err = services.Bridge_UpdatePagesAccessToken()
		if err != nil {
			fmt.Println("Lỗi khi cập nhật page access token:", err)
		}

		err = services.Bridge_SyncConversations()
		if err != nil {
			fmt.Println("Lỗi khi đồng bộ cuộc trò chuyện:", err)
		}

		err = services.Bridge_SyncMessages()
		if err != nil {
			fmt.Println("Lỗi khi đồng bộ cuộc trò chuyện:", err)
		}

		// Dừng 5 phút trước khi tiếp tục
		time.Sleep(50 * time.Minute)
	}
}
