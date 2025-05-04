/*
Package jobs chứa các job cụ thể của ứng dụng.
Mỗi job sẽ kế thừa từ scheduler.BaseJob và triển khai logic riêng.
*/
package jobs

import (
	"agent_pancake/app/integrations"
	"agent_pancake/app/scheduler"
	"context"
	"log"
	"time"
)

// SyncNewJob là job đồng bộ dữ liệu mới.
// Job này sẽ đồng bộ các dữ liệu mới từ khi lần cuối đồng bộ.
type SyncNewJob struct {
	*scheduler.BaseJob
}

// NewSyncNewJob tạo một instance mới của SyncNewJob.
// Tham số:
// - name: Tên định danh của job
// - schedule: Biểu thức cron định nghĩa lịch chạy
// Trả về một instance của SyncNewJob
func NewSyncNewJob(name, schedule string) *SyncNewJob {
	return &SyncNewJob{
		BaseJob: scheduler.NewBaseJob(name, schedule),
	}
}

// ExecuteInternal thực thi logic đồng bộ dữ liệu mới.
// Tham số:
// - ctx: Context để kiểm soát thời gian thực thi
// Trả về error nếu có lỗi xảy ra
func (j *SyncNewJob) ExecuteInternal(ctx context.Context) error {
	log.Printf("Bắt đầu đồng bộ dữ liệu mới - Job: %s", j.GetName())

	// Thực hiện đồng bộ
	SyncNewData(5)

	log.Printf("Hoàn thành đồng bộ dữ liệu mới - Job: %s", j.GetName())
	return nil
}

// SyncBaseAuth thực hiện xác thực và đồng bộ dữ liệu cơ bản
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

// SyncNewData đồng bộ dữ liệu mới nhất
func SyncNewData(sleepSeconds int) {
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
