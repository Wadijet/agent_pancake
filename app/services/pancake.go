package services

import (
	"agent_pancake/global"
	"agent_pancake/utility/httpclient"
	"errors"
	"log"
	"strconv"
	"time"
)

// Hàm PanCake_GetFbPages lấy danh sách pages từ server Pancake
func PanCake_GetFbPages(access_token string) (result map[string]interface{}, err error) {

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.PancakeBaseUrl, 60*time.Second)

	// Thiết lập header
	params := map[string]string{
		"access_token": access_token,
	}

	// Số lần thử request
	requestCount := 0
	for {
		requestCount++
		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return nil, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

		// Dừng 30s trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// Gửi yêu cầu GET
		resp, err := client.GET("/v1/pages", params)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Lấy danh sách trang Facebook thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		if result["success"] == true {
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm PanCake_GeneratePageAccessToken tạo page_access_token từ server Pancake
func PanCake_GeneratePageAccessToken(page_id string, access_token string) (result map[string]interface{}, err error) {

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.PancakeBaseUrl, 10*time.Second)

	// Chuẩn bị dữ liệu cần gửi
	params := map[string]string{
		"access_token": access_token,
	}

	// Số lần thử request
	requestCount := 0
	for {
		requestCount++
		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return nil, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

		// Dừng 30s trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// Gửi yêu cầu POST
		resp, err := client.POST("/v1/pages/"+page_id+"/generate_page_access_token", nil, params)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Lấy page_access_token thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		if result["success"] == true {
			return result, nil
		} else {
			log.Println("Lấy page_access_token thất bại: ", result["message"])
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm Pancake_GetConversations_v1 lấy danh sách Conversations từ server Pancake
func Pancake_GetConversations_v1(page_id string, page_access_token string, since int64, until int64, page_number int) (result map[string]interface{}, err error) {
	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.PancakeBaseUrl, 60*time.Second)

	// Thiết lập header
	params := map[string]string{
		"page_access_token": page_access_token,
		"since":             strconv.FormatInt(since, 10),
		"until":             strconv.FormatInt(until, 10),
		"page_number":       strconv.Itoa(page_number),
	}

	// Số lần thử request
	requestCount := 0
	for {
		requestCount++

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return nil, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

		// Dừng 30s trước khi tiếp tục
		time.Sleep(500 * time.Millisecond)

		// Gửi yêu cầu GET
		resp, err := client.GET("/public_api/v1/pages/"+page_id+"/conversations", params)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Lấy danh sách cuộc trò chuyện thất bại. StatusCode=", resp.StatusCode, "Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		if result["success"] == true {
			return result, nil
		} else {
			log.Println("Lấy danh sách cuộc trò chuyện thất bại: ", result["message"])
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

	}
}

// Hàm Pancake_GetConversations_v1 lấy danh sách Conversations từ server Pancake
func Pancake_GetConversations_v2(page_id string, last_conversation_id string) (result map[string]interface{}, err error) {
	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.PancakeBaseUrl, 60*time.Second)

	// Số lần thử request
	requestCount := 0
	for {
		requestCount++

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return nil, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

	Start:

		page_access_token, err := Local_GetPageAccessToken(page_id)
		if err != nil {
			log.Println("Lỗi khi lấy page_access_token")
			return nil, err
		}

		// Thiết lập header
		params := map[string]string{
			"page_access_token":    page_access_token,
			"last_conversation_id": last_conversation_id,
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return nil, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

		// Dừng 30s trước khi tiếp tục
		time.Sleep(300 * time.Millisecond)

		// Gửi yêu cầu GET
		resp, err := client.GET("/public_api/v2/pages/"+page_id+"/conversations", params)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Lấy danh sách cuộc trò chuyện thất bại. StatusCode = ", resp.StatusCode, "Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		if result["success"] == true {
			return result, nil
		} else {
			errCode, _ := result["error_code"].(float64)
			if errCode == 105 {
				err = Local_UpdatePagesAccessToken(page_id)
				if err != nil {
					log.Println("Lỗi khi cập nhật page_access_token")
				}
				goto Start
			}

			log.Println("Lấy danh sách cuộc trò chuyện thất bại: message = ", result["message"])
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

	}
}

// Hàm Pancake_GetConversations lấy danh sách Conversations từ server Pancake
func Pancake_GetMessages(page_id string, conversation_id string, customer_id string) (result map[string]interface{}, err error) {
	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.PancakeBaseUrl, 60*time.Second)

	// Số lần thử request
	requestCount := 0
	for {
		requestCount++

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return nil, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

	Start:

		// Dừng 30s trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		page_access_token, err := Local_GetPageAccessToken(page_id)
		if err != nil {
			log.Println("Lỗi khi lấy page_access_token")
			return nil, err
		}

		// Thiết lập header
		params := map[string]string{
			"page_access_token": page_access_token,
			"customer_id":       customer_id,
		}

		// Gửi yêu cầu GET
		resp, err := client.GET("/public_api/v1/pages/"+page_id+"/conversations/"+conversation_id+"/messages", params)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Lấy danh sách tin nhắn thất bại. StatusCode=", resp.StatusCode, "Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		if result["success"] == true {
			return result, nil
		} else {
			errCode, _ := result["error_code"].(float64)
			if errCode == 105 {
				err = Local_UpdatePagesAccessToken(page_id)
				if err != nil {
					log.Println("Lỗi khi cập nhật page_access_token")
				}
				goto Start
			}
			log.Println("Lấy danh sách tin nhắn thất bại: ", result["message"])
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

	}
}
