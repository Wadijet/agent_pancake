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

		// Lấy danh sách access token
		accessTokens := services.FolkForm_GetAccessTokens()
		if len(accessTokens) > 0 {

			// duyệt qua từng access token để lấy danh sách trang
			for _, access_token := range accessTokens {
				// lấy danh sách Pages từ server PanCake, đưa vào server FolkForm
				services.PanCake_GetFbPages(access_token)

				// Cập nhật page access token cho từng page
				pages := services.FolkForm_GetFbPages()
				if len(pages) > 0 {
					// duyệt qua từng page để lấy access token
					for _, page := range pages {
						// chuyển page từ interface{} sang dạng map[string]interface{}
						page := page.(map[string]interface{})

						// lấy access token từ server PanCake, đưa vào server FolkForm
						services.PanCake_GeneratePageAccessToken(page["pageId"].(string), access_token)
					}
				} else {
					fmt.Println("Không có trang nào.")
				}

				// Lấy danh sách hội thoại từ server PanCake, đưa vào server FolkForm

			}
		} else {
			fmt.Println("Không có access token nào.")
		}

		// Dừng 5 phút trước khi tiếp tục
		time.Sleep(5 * time.Minute)
	}
}
