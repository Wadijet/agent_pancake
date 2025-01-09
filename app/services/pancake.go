package services

import (
	"agent_pancake/global"
	"agent_pancake/utility/httpclient"
	"errors"
	"log"
	"time"
)

// Hàm GetPageAccessToken sẽ gửi yêu cầu lấy access token của trang Facebook từ server
func PanCake_GeneratePageAccessToken(page_id string, access_token string) (err error) {
	// Số lần thử request
	requestCount := 0
	for {
		requestCount++

		if global.ApiToken == "" {
			// trả về lỗi
			return errors.New("Chưa đăng nhập. Thoát vòng lặp.")
		}

		// Khởi tạo client
		client := httpclient.NewHttpClient(global.GlobalConfig.PancakeBaseUrl, 10*time.Second)

		// Chuẩn bị dữ liệu cần gửi
		params := map[string]string{
			"access_token": access_token,
		}

		// Gửi yêu cầu POST
		resp, err := client.POST("/v1/pages/"+page_id+"/generate_page_access_token", nil, params)
		if err != nil {
			log.Fatal("Lỗi khi gọi API:", err)
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Lấy access token thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Fatal("Lỗi khi phân tích phản hồi:", err)
		}

		if result["success"] == true {
			page_access_token := result["page_access_token"].(string)
			FolkForm_UpdatePageAccessToken(page_id, page_access_token)

			return nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

		// Dừng 30s trước khi tiếp tục
		time.Sleep(30 * time.Second)
	}
}

// Hàm PanCake_GetFbPages sẽ gửi yêu cầu lấy danh sách trang Facebook từ server
func PanCake_GetFbPages(access_token string) (err error) {
	// Số lần thử request
	requestCount := 0
	for {
		requestCount++

		if global.ApiToken == "" {
			// trả về lỗi
			return errors.New("Chưa đăng nhập. Thoát vòng lặp.")
		}

		// Khởi tạo client
		client := httpclient.NewHttpClient(global.GlobalConfig.PancakeBaseUrl, 60*time.Second)
		// Thiết lập header
		params := map[string]string{
			"access_token": access_token,
		}

		// Gửi yêu cầu GET
		resp, err := client.GET("/v1/pages", params)
		if err != nil {
			log.Fatal("Lỗi khi gọi API:", err)
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Lấy danh sách trang Facebook thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Fatal("Lỗi khi phân tích phản hồi:", err)
		}

		if result["success"] == true {
			// Lấy dữ liệu từ phản hồi lưu ở data.items
			data := result["categorized"].(map[string]interface{})["activated"].([]interface{})
			for _, item := range data {
				FolkForm_SendFbPage(item)
			}
			return nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

		// Dừng 30s trước khi tiếp tục
		time.Sleep(30 * time.Second)
	}
}
