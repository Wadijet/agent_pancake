package main

import (
	"agent_pancake/app/integrations"
	"agent_pancake/app/jobs"
	"agent_pancake/app/scheduler"
	"agent_pancake/config"
	"agent_pancake/global"
	"context"
	"log"
	"time"
)

// Các Scheduler
var Scheduler = scheduler.NewScheduler() // Scheduler chứa các jobs

func main() {
	// Đọc dữ liệu từ file .env
	global.GlobalConfig = config.NewConfig()
	log.Println("Đã đọc cấu hình từ file .env")

	// Khởi tạo scheduler
	s := scheduler.NewScheduler()

	// Tạo một job mới chạy mỗi 5 phút
	syncNewJob := jobs.NewSyncNewJob("sync-new-job", "0 */5 * * * *")

	// Chạy job ngay lập tức
	ctx := context.Background()
	if err := syncNewJob.Execute(ctx); err != nil {
		log.Printf("Lỗi khi thực thi job lần đầu: %v", err)
	}

	// Thêm job vào scheduler để chạy theo lịch
	err := s.AddJob(syncNewJob.GetName(), syncNewJob.GetSchedule(), func() {
		ctx := context.Background()
		if err := syncNewJob.Execute(ctx); err != nil {
			log.Printf("Lỗi khi thực thi job %s: %v", syncNewJob.GetName(), err)
		}
	})
	if err != nil {
		log.Fatalf("Lỗi khi thêm job: %v", err)
	}

	// Khởi động scheduler
	s.Start()
	log.Println("Scheduler đã được khởi động")

	// Giữ chương trình chạy
	// Trong thực tế, bạn có thể thêm các logic khác ở đây
	select {}
}

func SyncBaseAuth() {

	// Nếu chưa đăng nhập thì đăng nhập
	_, err := integrations.FolkForm_CheckIn()
	if err != nil {
		log.Println("Chưa đăng nhập, tiến hành đăng nhập...")
		integrations.FolkForm_Login()
		integrations.FolkForm_CheckIn()
	}

	// Đồng bộ danh sách các pages từ pancake sang folkform
	err = integrations.Bridge_SyncPages()
	if err != nil {
		log.Println("Lỗi khi đồng bộ trang:", err)
	} else {
		log.Println("Đồng bộ trang thành công")
	}

	// Đồng bộ danh sách các pages từ folkform sang local
	err = integrations.Local_SyncPagesFolkformToLocal()
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
		err := integrations.Bridge_SyncConversationsFromCloud()
		if err != nil {
			log.Println("Lỗi khi đồng bộ cuộc trò chuyện:", err)
		} else {
			log.Println("Đồng bộ cuộc trò chuyện thành công")
		}

		// Đồng bộ danh sách các tin nhắn từ pancake sang folkform
		err = integrations.Bridge_SyncMessages()
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
		_, err := integrations.FolkForm_CheckIn()
		if err != nil {
			log.Println("Chưa đăng nhập, tiến hành đăng nhập...")
			integrations.FolkForm_Login()
			integrations.FolkForm_CheckIn()
		}

		log.Println("Bắt đầu đồng bộ dữ liệu mới nhất")
		integrations.Sync_NewMessagesOfAllPages()
		log.Println("Đồng bộ dữ liệu mới nhất thành công")

		// Dừng sleepSeconds giây trước khi tiếp tục
		log.Println("Dừng", sleepSeconds, "giây trước khi tiếp tục")
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}
}
