package services

import (
	"errors"
	"log"
	"time"

	"agent_pancake/global"

	"go.mongodb.org/mongo-driver/bson"
)

// Hàm Bridge_SyncPagesFolkformToLocal sẽ đồng bộ danh sách trang Facebook từ server FolkForm về server local
// - Lấy danh sách trang từ server FolkForm
// - Đẩy danh sách trang vào server local
func Local_SyncPagesFolkformToLocal() (resultErr error) {
	limit := 50
	page := 0

	for {
		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// Lấy danh sách các pages từ server FolkForm
		resultPages, err := FolkForm_GetFbPages(page, limit)
		if err != nil {
			return errors.New("Lỗi khi lấy danh sách trang Facebook")
		}

		data := resultPages["data"].(map[string]interface{})
		itemCount := data["itemCount"].(float64)

		if itemCount > 0 {
			items := data["items"].([]interface{})
			if len(items) > 0 {
				// Clear all data in global.PanCake_FbPages
				global.PanCake_FbPages = nil

				for _, item := range items {

					// chuyển item từ interface{} sang dạng global.FbPage
					var cloudFbPage global.FbPage
					bsonBytes, err := bson.Marshal(item)
					if err != nil {
						return err
					}

					err = bson.Unmarshal(bsonBytes, &cloudFbPage)
					if err != nil {
						return err
					}

					// Append cloudFbPage to global.PanCake_FbPages
					global.PanCake_FbPages = append(global.PanCake_FbPages, cloudFbPage)
				}
			}

			page++
			continue
		} else {
			break
		}

	}

	return nil
}

func local_UpdatePagesAccessToken(pageId string, page_access_token string) (resultErr error) {
	// Find page in global.PanCake_FbPages
	for index, page := range global.PanCake_FbPages {
		if page.PageId == pageId {
			global.PanCake_FbPages[index].PageAccessToken = page_access_token
		}
	}
	return nil
}

func Local_UpdatePagesAccessToken(pageId string) (resultErr error) {

	for i, page := range global.PanCake_FbPages {
		if page.PageId == pageId {
			access_token := page.AccessToken

			// Gọi hàm PanCake_GeneratePageAccessToken để lấy page_access_token
			resultGeneratePageAccessToken, err := PanCake_GeneratePageAccessToken(page.PageId, access_token)
			if err != nil {
				log.Println("Lỗi khi lấy page access token: ", err)
				continue
			}

			// chuyển resultGeneratePageAccessToken từ interface{} sang dạng map[string]interface{}
			page_access_token := resultGeneratePageAccessToken["page_access_token"].(string)
			global.PanCake_FbPages[i].PageAccessToken = page_access_token

			return nil
		}
	}

	return errors.New("Không tìm thấy page")
}

func Local_GetPageAccessToken(pageId string) (pageAccessToken string, resultErr error) {
	// Find page in global.PanCake_FbPages
	for _, page := range global.PanCake_FbPages {
		if page.PageId == pageId {
			return page.PageAccessToken, nil
		}
	}
	return "", errors.New("Không tìm thấy page")
}
