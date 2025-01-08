package main

import (
	"agent_pancake/config"
	"agent_pancake/global"
	"agent_pancake/utility/httpclient"
	"agent_pancake/utility/hwid"
	"errors"
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

		// Nếu chưa đăng nhập thì đăng nhập
		checkin := CheckIn()
		if checkin != nil {
			Login()
			CheckIn()
		}

		accessTokens := GetAccessTokens()
		if len(accessTokens) > 0 {
			fmt.Println("Danh sách access token:", accessTokens)

		} else {
			fmt.Println("Không có access token nào.")
		}

		// Dừng 5 phút trước khi tiếp tục
		time.Sleep(5 * time.Minute)
	}
}

// Hàm GetAccessTokens sẽ gửi yêu cầu lấy danh sách access token từ server
func GetAccessTokens() []string {
	// Khởi tạo mảng chứa access token
	accessTokens := []string{}

	// Khởi tạo client
	client := httpclient.NewHttpClient(global.GlobalConfig.ApiBaseUrl, 10*time.Second)
	// Thiết lập header
	client.SetHeader("Authorization", "Bearer "+global.ApiToken)
	// thêm param vào url với key là "page" và value là 0, limit là 10

	params := map[string]string{
		"page":  "0",
		"limit": "10",
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
		// Lấy dữ liệu từ phản hồi lưu ở data.items
		data := result["data"].(map[string]interface{})["items"].([]interface{})
		for _, item := range data {
			accessTokens = append(accessTokens, item.(map[string]interface{})["value"].(string))
		}
	}

	return accessTokens
}

// Hàm Login để Agent login vào hệ thống
func Login() {
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
func CheckIn() (err error) {

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
