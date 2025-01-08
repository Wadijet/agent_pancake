package main

import (
	"agent_pancake/app/httpclient"
	"agent_pancake/config"
	"agent_pancake/global"
	"fmt"
	"log"
	"time"
)

func main() {
	// Đọc dữ liệu từ file .env
	global.GlobalConfig = config.NewConfig()

	// Vòng lặp vô hạn chạy 5 phút một lần
	for {
		// Công việc cần thực hiện
		fmt.Println("Thực hiện công việc tại:", time.Now())

		// Dừng 5 phút trước khi tiếp tục
		time.Sleep(5 * time.Minute)
	}
}

// hàm Điểm danh sẽ gửi thông tin điểm danh lên server
func CheckIn() {
	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)

	// Thiết lập header
	client.SetHeader("Authorization", "Bearer your-token")

	// Gửi yêu cầu GET
	resp, err := client.GET("/posts/1")
	if err != nil {
		log.Fatal("Lỗi khi gọi API:", err)
	}

	// Đọc dữ liệu từ phản hồi
	var result map[string]interface{}
	if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
		log.Fatal("Lỗi khi phân tích phản hồi:", err)
	}

	// Hiển thị kết quả
	fmt.Println("Kết quả:", result)
}
