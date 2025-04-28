package integrations

import (
	"agent_pancake/global"
	"agent_pancake/utility/httpclient"
	"agent_pancake/utility/hwid"
	"errors"
	"log"
	"strconv"
	"time"
)

// Hàm FolkForm_CreateMessage sẽ gửi yêu cầu tạo tin nhắn lên server
func FolkForm_CreateMessage(pageId string, pageUsername string, conversationId string, customerId string, messageData interface{}) (result map[string]interface{}, err error) {

	if global.ApiToken == "" {
		// trả về lỗi
		return nil, errors.New("Chưa đăng nhập. Thoát vòng lặp.")
	}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 60*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)
	// Chuẩn bị dữ liệu cần gửi
	data := map[string]interface{}{
		"pageId":         pageId,
		"pageUsername":   pageUsername,
		"conversationId": conversationId,
		"customerId":     customerId,
		"panCakeData":    messageData,
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
		//time.Sleep(100 * time.Millisecond)

		// Gửi yêu cầu POST
		resp, err := client.POST("/fb_messages", data, nil)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Gửi tin nhắn thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		if result["status"] == "success" {
			log.Println("Gửi tin nhắn thành công")
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm FolkForm_GetConversations sẽ gửi yêu cầu lấy danh sách hội thoại từ server
func FolkForm_GetConversations(page int, limit int) (result map[string]interface{}, err error) {

	if global.ApiToken == "" {
		// trả về lỗi
		return nil, errors.New("Chưa đăng nhập. Thoát vòng lặp.")
	}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)
	// thêm param vào url với key là "page" và value là 0, limit là 10
	params := map[string]string{
		"page":  strconv.Itoa(page),
		"limit": strconv.Itoa(limit),
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
		resp, err := client.GET("/fb_conversations", params)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		// Lấy dữ liệu từ phản hồi
		if result["status"] == "success" {
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm FolkForm_GetConversations sẽ gửi yêu cầu lấy danh sách hội thoại từ server
func FolkForm_GetConversationsWithPageId(page int, limit int, pageId string) (result map[string]interface{}, err error) {

	if global.ApiToken == "" {
		// trả về lỗi
		return nil, errors.New("Chưa đăng nhập. Thoát vòng lặp.")
	}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)
	// thêm param vào url với key là "page" và value là 0, limit là 10
	params := map[string]string{
		"page":   strconv.Itoa(page),
		"limit":  strconv.Itoa(limit),
		"pageId": pageId,
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
		resp, err := client.GET("/fb_conversations/newest", params)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		// Lấy dữ liệu từ phản hồi
		if result["status"] == "success" {
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm Folkform_CreateConversation sẽ gửi yêu cầu tạo hội thoại lên server
func FolkForm_CreateConversation(pageId string, pageUsername string, conversation_data interface{}) (result map[string]interface{}, err error) {

	if global.ApiToken == "" {
		// trả về lỗi
		return nil, errors.New("Chưa đăng nhập. Thoát vòng lặp.")
	}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 60*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)
	// Chuẩn bị dữ liệu cần gửi
	data := map[string]interface{}{
		"pageId":       pageId,
		"pageUsername": pageUsername,
		"panCakeData":  conversation_data,
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
		//time.Sleep(100 * time.Millisecond)

		// Gửi yêu cầu POST
		resp, err := client.POST("/fb_conversations", data, nil)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Gửi hội thoại thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		if result["status"] == "success" {
			log.Println("Gửi hội thoại thành công")
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm FolkForm_GetFbPageById sẽ gửi yêu cầu lấy thông tin trang Facebook từ server
func FolkForm_GetFbPageById(id string) (result map[string]interface{}, err error) {

	if global.ApiToken == "" {
		// trả về lỗi
		return nil, errors.New("Chưa đăng nhập. Thoát vòng lặp.")
	}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)

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
		resp, err := client.GET("/fb_pages/"+id, nil)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		// Lấy dữ liệu từ phản hồi
		if result["status"] == "success" {
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm FolkForm_GetFbPageByPageId sẽ gửi yêu cầu lấy thông tin trang Facebook từ server
func FolkForm_GetFbPageByPageId(pageId string) (result map[string]interface{}, err error) {

	if global.ApiToken == "" {
		// trả về lỗi
		return nil, errors.New("Chưa đăng nhập. Thoát vòng lặp.")
	}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)

	// Số lần thử request
	requestCount := 0
	for {
		requestCount++
		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return nil, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

		// Dừng 30s trước khi tiếp tục
		//time.Sleep(100 * time.Millisecond)

		// Gửi yêu cầu GET
		resp, err := client.GET("/fb_pages/pageId/"+pageId, nil)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		// Lấy dữ liệu từ phản hồi
		if result["status"] == "success" {
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm FolkForm_GetFbPages sẽ gửi yêu cầu lấy danh sách trang Facebook từ server
func FolkForm_GetFbPages(page int, limit int) (result map[string]interface{}, err error) {

	if global.ApiToken == "" {
		// trả về lỗi
		return nil, errors.New("Chưa đăng nhập. Thoát vòng lặp.")
	}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)
	// thêm param vào url với key là "page" và value là 0, limit là 10

	// Chuẩn bị params cho yêu cầu GET
	params := map[string]string{
		"page":  strconv.Itoa(page),
		"limit": strconv.Itoa(limit),
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
		resp, err := client.GET("/fb_pages", params)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		// Lấy dữ liệu từ phản hồi
		if result["status"] == "success" {
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm FolkForm_UpdatePageAccessToken sẽ gửi yêu cầu cập nhật access token của trang Facebook lên server
func FolkForm_UpdatePageAccessToken(page_id string, page_access_token string) (result map[string]interface{}, err error) {

	if global.ApiToken == "" {
		// trả về lỗi
		return nil, errors.New("Chưa đăng nhập. Thoát vòng lặp.")
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
		resp, err := client.POST("/fb_pages/update_token", data, nil)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Cập nhật page_access_token thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		if result["status"] == "success" {
			log.Println("Cập nhật page_access_token thành công")
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm FolkForm_CreateFbPage sẽ gửi yêu cầu lưu trang Facebook lên server
func FolkForm_CreateFbPage(access_token string, page_data interface{}) (result map[string]interface{}, err error) {

	if global.ApiToken == "" {
		// trả về lỗi
		return nil, errors.New("Chưa đăng nhập. Thoát vòng lặp.")
	}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 60*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)
	// Chuẩn bị dữ liệu cần gửi
	data := map[string]interface{}{
		"accessToken": access_token,
		"panCakeData": page_data,
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
		//time.Sleep(100 * time.Millisecond)

		// Gửi yêu cầu POST
		resp, err := client.POST("/fb_pages", data, nil)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Gửi trang Facebook thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		if result["status"] == "success" {
			log.Println("Gửi trang Facebook thành công")
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm FolkForm_GetAccessTokens sẽ gửi yêu cầu lấy danh sách access token từ server
func FolkForm_GetAccessTokens(page int, limit int) (result map[string]interface{}, err error) {

	if global.ApiToken == "" {
		// trả về lỗi
		return nil, errors.New("Chưa đăng nhập. Thoát vòng lặp.")
	}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)

	// thêm param vào url với key là "page" và value là 0, limit là 10
	params := map[string]string{
		"page":  strconv.Itoa(page),
		"limit": strconv.Itoa(limit),
	}

	// Số lần thử request
	requestCount := 0
	for {
		requestCount++

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

		// Dừng 30s trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// Gửi yêu cầu GET
		resp, err := client.GET("/access_tokens", params)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		// Lấy dữ liệu từ phản hồi
		if result["status"] == "success" {
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm FolkForm_Login để Agent login vào hệ thốnga
func FolkForm_Login() (result map[string]interface{}, resultError error) {
	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)

	// Số lần thử đăng nhập
	requestCount := 0
	for {
		// Tăng số lần thử lên 1
		requestCount++

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return nil, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}

		// Dừng 30s trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// lấy hardware ID
		hwid, err := hwid.GenerateHardwareID()
		if err != nil {
			log.Println("Lỗi khi lấy Hardware ID:", err)
			continue
		}

		// Chuẩn bị dữ liệu cần gửi
		data := map[string]interface{}{
			"email":    global.GlobalConfig.Email,
			"password": global.GlobalConfig.Password,
			"hwid":     hwid,
		}

		// Gửi yêu cầu POST
		resp, err := client.POST("/users/login", data, nil)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Đăng nhập thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		// Lưu token vào biến toàn cục
		if result["status"] == "success" {
			log.Println("Đăng nhập thành công")
			// Lưu token vào biến toàn cục
			global.ApiToken = result["data"].(map[string]interface{})["token"].(string)
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}

// Hàm Điểm danh sẽ gửi thông tin điểm danh lên server
func FolkForm_CheckIn() (result map[string]interface{}, err error) {

	if global.ApiToken == "" {
		// trả về lỗi
		return nil, errors.New("Chưa đăng nhập. Thoát vòng lặp.")
	}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)

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
		resp, err := client.POST("/agents/checkin/"+global.GlobalConfig.AgentId, nil, nil)
		if err != nil {
			log.Println("Lỗi khi gọi API:", err)
			continue
		}

		// Kiểm tra mã trạng thái, nếu không phải 200 thì thử lại
		if resp.StatusCode != 200 {
			log.Println("Điểm danh thất bại. Thử lại lần thứ", requestCount)
			continue
		}

		// Đọc dữ liệu từ phản hồi
		var result map[string]interface{}
		if err := httpclient.ParseJSONResponse(resp, &result); err != nil {
			log.Println("Lỗi khi phân tích phản hồi:", err)
			continue
		}

		if result["status"] == "success" {
			log.Println("Điểm danh thành công")
			return result, nil
		}

		// Nếu số lần thử vượt quá 5 lần thì thoát vòng lặp
		if requestCount > 5 {
			return result, errors.New("Đã thử quá nhiều lần. Thoát vòng lặp.")
		}
	}
}
