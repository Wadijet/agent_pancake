package config

import (
	"fmt"
	"log"

	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Configuration chứa thông tin tĩnh cần thiết để chạy ứng dụng
// Nó chứa thông tin cơ sở dữ liệu
type Configuration struct {
	Email      string `env:"EMAIL,required"`        // Chế độ khởi tạo
	Password   string `env:"PASSWORD,required"`     // Địa chỉ server
	AgentId    string `env:"AGENT_ID,required"`     // Bí mật JWT
	ApiBaseUrl string `env:"API_BASE_URL,required"` // Địa chỉ server
}

// NewConfig sẽ đọc dữ liệu cấu hình từ file .env được cung cấp
func NewConfig(files ...string) *Configuration {
	err := godotenv.Load(filepath.Join(".env")) // Tải cấu hình từ file .env
	if err != nil {
		log.Printf("Không tìm thấy file .env %q\n", files)
	}

	cfg := Configuration{}

	// Phân tích env vào cấu hình
	err = env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg
}
