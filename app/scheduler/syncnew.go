/*
Package scheduler cung cấp ví dụ về cách tạo một job mới.
File này minh họa cách triển khai một job cụ thể bằng cách kế thừa từ BaseJob.
*/

package scheduler

import (
	"agent_pancake/app/integrations"
	"context"
	"log"
	"sync"
	"time"
)

// SyncNewJob là một job mẫu minh họa cách triển khai một job cụ thể.
// Job này sẽ in ra thông báo mỗi khi được thực thi.
type SyncNewJob struct {
	*BaseJob
	mu        sync.Mutex // Mutex để kiểm soát trạng thái job
	isRunning bool       // Trạng thái đang chạy của job
}

// NewExampleJob tạo một instance mới của ExampleJob.
// Tham số:
// - name: Tên định danh của job
// - schedule: Biểu thức cron định nghĩa lịch chạy
// Trả về một instance của ExampleJob
func NewSyncNewJob(name, schedule string) *SyncNewJob {
	return &SyncNewJob{
		BaseJob:   NewBaseJob(name, schedule),
		mu:        sync.Mutex{},
		isRunning: false,
	}
}

// ExecuteInternal thực thi logic riêng của job.
// Trong ví dụ này, job sẽ:
// 1. In ra thông báo bắt đầu
// 2. Thực thi công việc đồng bộ
// 3. In ra thông báo kết thúc
// Tham số:
// - ctx: Context để kiểm soát thời gian thực thi
// Trả về error nếu có lỗi xảy ra
func (j *SyncNewJob) ExecuteInternal(ctx context.Context) error {
	log.Printf("Bắt đầu thực thi job: %s", j.GetName())

	// Thực thi công việc
	SyncNewData(5)

	log.Printf("Hoàn thành job: %s", j.GetName())
	return nil
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

// SyncNewData sẽ đồng bộ dữ liệu mới nhất, chỉ update những dữ liệu mới nhất kể từ lần cuối cùng đồng bộ
func SyncNewData(sleepSeconds int) {
	// Vòng lặp vô hạn chạy 5 phút một lần
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
