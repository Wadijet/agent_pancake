package hwid

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
)

// Hàm lấy thông tin MAC Address
func getMACAddress() (string, error) {
	cmd := exec.Command("getmac") // Chạy lệnh "getmac" trên Windows
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "-") { // MAC Address có định dạng chứa "-"
			mac := strings.Fields(line)[0]
			return mac, nil
		}
	}
	return "", fmt.Errorf("không tìm thấy MAC Address")
}

// Hàm tạo Hardware ID từ MAC Address
func GenerateHardwareID() (string, error) {
	macAddress, err := getMACAddress()
	if err != nil {
		return "", err
	}

	// Hash MAC Address bằng MD5
	hash := md5.New()
	hash.Write([]byte(macAddress))
	return hex.EncodeToString(hash.Sum(nil)), nil
}
