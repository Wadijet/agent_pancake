package main

import (
	"agent_pancake/config"
	"agent_pancake/global"
	"log"
)

func main() {
	// Đọc dữ liệu từ file .env
	global.GlobalConfig = config.NewConfig()
	log.Println("Đã đọc cấu hình từ file .env")

}

/*

func SyncBaseAuth() {

	// Nếu chưa đăng nhập thì đăng nhập
	_, err := services.FolkForm_CheckIn()
	if err != nil {
		log.Println("Chưa đăng nhập, tiến hành đăng nhập...")
		services.FolkForm_Login()
		services.FolkForm_CheckIn()
	}

	// Đồng bộ danh sách các pages từ pancake sang folkform
	err = services.Bridge_SyncPages()
	if err != nil {
		log.Println("Lỗi khi đồng bộ trang:", err)
	} else {
		log.Println("Đồng bộ trang thành công")
	}

	// Đồng bộ danh sách các pages từ folkform sang local
	err = services.Local_SyncPagesFolkformToLocal()
	if err != nil {
		log.Println("Lỗi khi đồng bộ trang:", err)
	} else {
		log.Println("Đồng bộ trang từ folkform sang local thành công")
	}
}

// SyncAllData sẽ đồng bộ tất cả dữ liệu vào cuối ngày, mỗi ngày chạy 1 lần
func SyncAllData(sleepMinutes int) {

	for {
		SyncBaseAuth()

		// Công việc cần thực hiện
		log.Println("Thực hiện công việc tại:", time.Now())

		// Đồng bộ danh sách các hội thoại từ pancake sang folkform
		err := services.Bridge_SyncConversationsFromCloud()
		if err != nil {
			log.Println("Lỗi khi đồng bộ cuộc trò chuyện:", err)
		} else {
			log.Println("Đồng bộ cuộc trò chuyện thành công")
		}

		// Đồng bộ danh sách các tin nhắn từ pancake sang folkform
		err = services.Bridge_SyncMessages()
		if err != nil {
			log.Println("Lỗi khi đồng bộ tin nhắn:", err)
		} else {
			log.Println("Đồng bộ tin nhắn thành công")
		}

		// Dừng 5 phút trước khi tiếp tục
		log.Println("Dừng", sleepMinutes, "phút trước khi tiếp tục")
		time.Sleep(time.Duration(sleepMinutes) * time.Minute)
	}
}

// SyncNewData sẽ đồng bộ dữ liệu mới nhất, chỉ update những dữ liệu mới nhất kể từ lần cuối cùng đồng bộ
func SyncNewData(sleepSeconds int) {
	// Vòng lặp vô hạn chạy 5 phút một lần
	for {
		SyncBaseAuth()

		// Nếu chưa đăng nhập thì đăng nhập
		_, err := services.FolkForm_CheckIn()
		if err != nil {
			log.Println("Chưa đăng nhập, tiến hành đăng nhập...")
			services.FolkForm_Login()
			services.FolkForm_CheckIn()
		}

		log.Println("Bắt đầu đồng bộ dữ liệu mới nhất")
		services.Sync_NewMessagesOfAllPages()
		log.Println("Đồng bộ dữ liệu mới nhất thành công")

		// Dừng sleepSeconds giây trước khi tiếp tục
		log.Println("Dừng", sleepSeconds, "giây trước khi tiếp tục")
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}
}
*/
