package services

import (
	"agent_pancake/global"
	"agent_pancake/utility/httpclient"
	"agent_pancake/utility/hwid"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Hàm tương tác với server FolkForm

// Hàm FolkForm_GetFbPages sẽ gửi yêu cầu lấy danh sách trang Facebook từ server
func FolkForm_GetFbPages() (pages []interface{}) {
	// Khởi tạo mảng chứa trang
	pages = []interface{}{}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)
	// thêm param vào url với key là "page" và value là 0, limit là 10

	// Chuẩn bị params cho yêu cầu GET
	limit := 10
	page := 0

	for {

		params := map[string]string{
			"page":  strconv.Itoa(page),
			"limit": strconv.Itoa(limit),
		}

		// Gửi yêu cầu GET
		resp, err := client.GET("/fb_pages", params)
		if err != nil {
			log.Fatal("Lỗi khi gọi API:", err)
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Fatal("Lỗi khi phân tích phản hồi:", err)
		}

		// Lấy dữ liệu từ phản hồi
		if result["status"] == "success" {
			if result["data"].(map[string]interface{})["itemCount"].(float64) > 0 {
				// Lấy dữ liệu từ phản hồi lưu ở data.items
				data := result["data"].(map[string]interface{})["items"].([]interface{})
				for _, item := range data {
					pages = append(pages, item)
				}

				// Tăng số trang lên 1
				page++
				continue
			} else {
				break
			}

		}
	}

	return pages
}

// Hàm FolkForm_UpdatePageAccessToken sẽ gửi yêu cầu cập nhật access token của trang Facebook lên server
func FolkForm_UpdatePageAccessToken(page_id string, page_access_token string) (err error) {
	// Số lần thử request
	requestCount := 0
	for {
		requestCount++

		if global.ApiToken == "" {
			// trả về lỗi
			return errors.New("Chưa đăng nhập. Thoát vòng lặp.")
		}

		// Khởi tạo client
		client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
		// Thiết lập header
		client.SetHeader("Authorization", "Bearer "+global.ApiToken)

		// Chuẩn bị dữ liệu cần gửi
		data := map[string]interface{}{
			"pageId":          page_id,
			"pageAccessToken": page_access_token,
		}

		// Gửi yêu cầu POST
		resp, err := client.POST("/fb_pages/update_token", data, nil)
		if err != nil {
			log.Fatal("Lỗi khi gọi API:", err)
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Cập nhật access token thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Fatal("Lỗi khi phân tích phản hồi:", err)
		}

		if result["status"] == "success" {
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

// Hàm FolkForm_CreateFbPage sẽ gửi yêu cầu lưu trang Facebook lên server
func FolkForm_CreateFbPage(page_data interface{}) (err error) {
	// Số lần thử request
	requestCount := 0
	for {
		requestCount++

		if global.ApiToken == "" {
			// trả về lỗi
			return errors.New("Chưa đăng nhập. Thoát vòng lặp.")
		}

		// Khởi tạo client
		client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 60*time.Second)
		// Thiết lập header
		client.SetHeader("Authorization", "Bearer "+global.ApiToken)

		// Chuẩn bị dữ liệu cần gửi
		data := map[string]interface{}{
			"apiData": page_data,
		}

		// Gửi yêu cầu POST
		resp, err := client.POST("/fb_pages", data, nil)
		if err != nil {
			log.Fatal("Lỗi khi gọi API:", err)
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Gửi trang Facebook thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Fatal("Lỗi khi phân tích phản hồi:", err)
		}

		if result["status"] == "success" {
			log.Println("Gửi trang Facebook thành công")
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

// Hàm FolkForm_GetAccessTokens sẽ gửi yêu cầu lấy danh sách access token từ server
func FolkForm_GetAccessTokens() []string {
	// Khởi tạo mảng chứa access token
	accessTokens := []string{}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)
	// thêm param vào url với key là "page" và value là 0, limit là 10

	// Chuẩn bị params cho yêu cầu GET
	limit := 10
	page := 0

	for {

		params := map[string]string{
			"page":  strconv.Itoa(page),
			"limit": strconv.Itoa(limit),
		}

		// Gửi yêu cầu GET
		resp, err := client.GET("/access_tokens", params)
		if err != nil {
			log.Fatal("Lỗi khi gọi API:", err)
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Fatal("Lỗi khi phân tích phản hồi:", err)
		}

		// Lấy dữ liệu từ phản hồi
		if result["status"] == "success" {
			if result["data"].(map[string]interface{})["itemCount"].(float64) > 0 {
				// Lấy dữ liệu từ phản hồi lưu ở data.items
				data := result["data"].(map[string]interface{})["items"].([]interface{})
				for _, item := range data {
					accessTokens = append(accessTokens, item.(map[string]interface{})["value"].(string))
				}

				// Tăng số trang lên 1
				page++
				continue
			} else {
				break
			}

		}
	}

	return accessTokens
}

// Hàm FolkForm_Login để Agent login vào hệ thốnga
func FolkForm_Login() {
	// Số lần thử đăng nhập
	requestCount := 0
	for {
		// Tăng số lần thử lên 1
		requestCount++

		// lấy hardware ID
		hwid, err := hwid.GenerateHardwareID()
		if err != nil {
			log.Fatal("Lỗi khi lấy Hardware ID:", err)
		}

		// Khởi tạo client
		client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
		//client.SetHeader("Content-Type", "application/json")

		// Chuẩn bị dữ liệu cần gửi
		data := map[string]interface{}{
			"email":    global.GlobalConfig.Email,
			"password": global.GlobalConfig.Password,
			"hwid":     hwid,
		}

		// Gửi yêu cầu POST
		resp, err := client.POST("/users/login", data, nil)
		if err != nil {
			log.Fatal("Lỗi khi gọi API:", err)
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Đăng nhập thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Fatal("Lỗi khi phân tích phản hồi:", err)
		}

		// Lưu token vào biến toàn cục
		if result["status"] == "success" {
			log.Println("Đăng nhập thành công")
			// Lưu token vào biến toàn cục
			global.ApiToken = result["data"].(map[string]interface{})["token"].(string)
			return
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			log.Fatal("Đã thử quá nhiều lần. Thoát vòng lặp.")
			break
		}

		// Dừng 30s trước khi tiếp tục
		time.Sleep(30 * time.Second)
	}

}

// Hàm Điểm danh sẽ gửi thông tin điểm danh lên server
func FolkForm_CheckIn() (err error) {

	requestCount := 0
	for {
		requestCount++

		if global.ApiToken == "" {
			// trả về lỗi
			return errors.New("Chưa đăng nhập. Thoát vòng lặp.")
		}

		// Khởi tạo client
		client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
		// Thiết lập header
		client.SetHeader("Authorization", "Bearer "+global.ApiToken)

		// Gửi yêu cầu POST
		resp, err := client.POST("/agents/checkin/"+global.GlobalConfig.AgentId, nil, nil)
		if err != nil {
			log.Fatal("Lỗi khi gọi API:", err)
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Điểm danh thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Fatal("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		if result["status"] == "success" {
			log.Println("Điểm danh thành công")
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

// ========================================================================================================
// Hàm xử lý logic trên server FolkForm

// Hàm FolkForm_UpdarePageAccessToken sẽ cập nhật page access token của trang Facebook lên server bằng cách:
// - Gửi yêu cầu tạo page access token từ server PanCake
// - Lấy page access token từ phản hồi và cập nhật lên server FolkForm
func FolkForm_UpdatePagesAccessToken() (err error) {
	// Số lần thử request
	requestCount := 0
	for {
		requestCount++

		// Lấy danh sách access token
		accessTokens := FolkForm_GetAccessTokens()
		if len(accessTokens) > 0 {

			// duyệt qua từng access token để lấy danh sách trang
			for _, access_token := range accessTokens {
				// lấy danh sách Pages từ server PanCake, đưa vào server FolkForm
				PanCake_GetFbPages(access_token)

				// Cập nhật page access token cho từng page
				pages := FolkForm_GetFbPages()
				if len(pages) > 0 {
					// duyệt qua từng page để lấy access token
					for _, page := range pages {
						// chuyển page từ interface{} sang dạng map[string]interface{}
						page := page.(map[string]interface{})
						page_id := page["pageId"].(string)

						page_access_token, err := PanCake_GeneratePageAccessToken(page_id, access_token)
						if page_access_token != "" {
							err = FolkForm_UpdatePageAccessToken(page_id, page_access_token)
							if err != nil {
								log.Fatal("Lỗi khi cập nhật page access token:", err)
							}
						} else {
							log.Fatal("Lỗi khi lấy page access token:", err)
						}
					}

				} else {
					fmt.Println("Không có trang nào.")
				}

				// Lấy danh sách hội thoại từ server PanCake, đưa vào server FolkForm

			}

		} else {
			fmt.Println("Không có access token nào.")
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

		// Dừng 30s trước khi tiếp tục
		time.Sleep(30 * time.Second)
	}
}

// Hàm FolkForm_GetConversations sẽ gửi yêu cầu lấy danh sách hội thoại từ server
func FolkForm_UpdateAllConversations() (err error) {
	// Số lần thử request
	requestCount := 0
	for {
		requestCount++

		// Cập nhật page access token cho từng page
		pages := FolkForm_GetFbPages()
		if len(pages) > 0 {
			// duyệt qua từng page để lấy access token
			for _, page := range pages {
				// chuyển page từ interface{} sang dạng map[string]interface{}
				page := page.(map[string]interface{})
				page_id := page["pageId"].(string)
				page_access_token := page["pageAccessToken"].(string)

				// in ra thông tin page
				fmt.Println("Page ID:", page_id)
				fmt.Println("Page Access Token:", page_access_token)

			}
		} else {
			fmt.Println("Không có trang nào.")
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

		// Dừng 30s trước khi tiếp tục
		time.Sleep(30 * time.Second)
	}
}
