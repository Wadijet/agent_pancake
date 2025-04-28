package global

import (
	"agent_pancake/config"

	"agent_pancake/app/scheduler"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var GlobalConfig *config.Configuration
var ApiToken string = ""

type FbPage struct {
	Id              primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`        // ID của quyền
	PageName        string                 `json:"pageName" bson:"pageName"`                 // Tên của trang
	PageUsername    string                 `json:"pageUsername" bson:"pageUsername"`         // Tên người dùng của trang
	PageId          string                 `json:"pageId" bson:"pageId" index:"unique;text"` // ID của trang
	IsSync          bool                   `json:"isSync" bson:"isSync"`                     // Trạng thái đồng bộ
	AccessToken     string                 `json:"accessToken" bson:"accessToken"`
	PageAccessToken string                 `json:"pageAccessToken" bson:"pageAccessToken"` // Mã truy cập của trang
	ApiData         map[string]interface{} `json:"apiData" bson:"apiData"`                 // Dữ liệu API
	CreatedAt       int64                  `json:"createdAt" bson:"createdAt"`             // Thời gian tạo quyền
	UpdatedAt       int64                  `json:"updatedAt" bson:"updatedAt"`             // Thời gian cập nhật quyền
}

var PanCake_FbPages []FbPage

// Các Scheduler
var Scheduler = scheduler.NewScheduler() // Scheduler chứa các jobs
