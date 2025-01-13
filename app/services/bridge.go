package services

import (
	"errors"
	"log"
	"time"
)

// ========================================================================================================
// Hàm xử lý logic trên server FolkForm
// ========================================================================================================

// Hàm Bridge_SyncPages(access_token string) sẽ đồng bộ danh sách trang Facebook từ server Pancake về server FolkForm
// - Lấy danh sách trang từ server Pancake
// - Đẩy danh sách trang vào server FolkForm
func bridge_SyncPagesOfAccessToken(access_token string) (resultErr error) {

	// Lấy danh sách trang từ server Pancake
	resultPages, err := PanCake_GetFbPages(access_token)
	if err != nil {
		return errors.New("Lỗi khi lấy danh sách trang Facebook")
	}

	// Lấy data lưu trong resultPages dạng []interface{} ở categorizedactivated
	activePages := resultPages["categorized"].(map[string]interface{})["activated"].([]interface{})
	for _, page := range activePages {

		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		FolkForm_CreateFbPage(access_token, page)
	}

	return nil
}

// Hàm Bridge_SyncPages sẽ đồng bộ danh sách trang Facebook từ server Pancake về server FolkForm
// - Lấy danh sách access token từ server FolkForm
// - Gọi hàm Bridge_SyncPagesOfAccessToken để đồng bộ trang của từng access token
func Bridge_SyncPages() (resultErr error) {

	log.Println("Bắt đầu đồng bộ trang Facebook từ server Pancake về server FolkForm...")

	limit := 50
	page := 0

	for {

		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// Lấy danh sách access token
		accessTokens, err := FolkForm_GetAccessTokens(page, limit)
		if err != nil {
			return errors.New("Lỗi khi lấy danh sách access token")
		}

		data := accessTokens["data"].(map[string]interface{})

		itemCount := data["itemCount"].(float64)
		if itemCount > 0 {
			items := data["items"].([]interface{})
			if len(items) > 0 {
				for _, item := range items {

					// Dừng nửa giây trước khi tiếp tục
					time.Sleep(100 * time.Millisecond)
					access_token := item.(map[string]interface{})["value"].(string)
					bridge_SyncPagesOfAccessToken(access_token)
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

// Hàm FolkForm_UpdarePageAccessToken sẽ cập nhật page_access_token của trang Facebook trên server FolkForm bằng cách:
// - Gửi yêu cầu tạo page_access_token lên server PanCake
// - Lấy page_access_token từ phản hồi và cập nhật lên server FolkForm
func Bridge_UpdatePagesAccessToken() (resultErr error) {

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
				for _, item := range items {

					// Dừng nửa giây trước khi tiếp tục
					time.Sleep(100 * time.Millisecond)

					// chuyển item từ interface{} sang dạng map[string]interface{}
					page := item.(map[string]interface{})
					page_id := page["pageId"].(string)
					access_token := page["accessToken"].(string)
					old_page_access_token := page["pageAccessToken"].(string)

					// Gọi hàm PanCake_GetFbPages để test page_access_token có hợp lệ không
					_, err := Pancake_GetConversations_v2(page_id, page_id, old_page_access_token)
					if err == nil {
						// Nếu page_access_token hợp lệ thì tiếp tục
						log.Println("Page_access_token vẫn còn hiệu lực")
						continue
					}

					// Gọi hàm PanCake_GeneratePageAccessToken để lấy page_access_token
					resultGeneratePageAccessToken, err := PanCake_GeneratePageAccessToken(page_id, access_token)
					if err != nil {
						log.Println("Lỗi khi lấy page access token: ", err)
						continue
					}

					// chuyển resultGeneratePageAccessToken từ interface{} sang dạng map[string]interface{}
					page_access_token := resultGeneratePageAccessToken["page_access_token"].(string)
					// Gọi hàm FolkForm_UpdatePageAccessToken để cập nhật page_access_token
					_, err = FolkForm_UpdatePageAccessToken(page_id, page_access_token)
					if err != nil {
						log.Println("Lỗi khi cập nhật page access token: ", err)
						continue
					}
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

// Hàm bridge_SyncConversationsOfPage_v1 sẽ đồng bộ danh sách hội thoại của trang Facebook từ server Pancake về server FolkForm
// - Lấy danh sách hội thoại của page từ server Pancake
// - Đẩy danh sách hội thoại vào server FolkForm
func bridge_SyncConversationsOfPage_v1(page_id string, page_username string, page_access_token string) (resultErr error) {

	// Lấy danh sách hội thoại từ server Pancake
	until := time.Now().Unix()
	since := until - 60*60*24*30

	// vòng lặp để lần lấy lần lượt từng khoảng since-until mỗi 30 ngày
	for {

		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		page_number := 1
		foundAnyConversation := false
		// vòng lặp để lấy page từ 1 đến khi hết dũ liệu
		for {

			// Dừng nửa giây trước khi tiếp tục
			time.Sleep(100 * time.Millisecond)

			resultGetConversations, err := Pancake_GetConversations_v1(page_id, page_access_token, since, until, page_number)
			if err != nil {
				log.Println("Lỗi khi lấy danh sách hội thoại:", err)
				break
			}

			total := resultGetConversations["total"].(float64)
			if total > 0 { // Nếu có dữ liệu thì tăng page lên 1
				// Lấy dữ liệu từ phản hồi lưu ở conversations
				conversations := resultGetConversations["conversations"].([]interface{})
				for _, conversation := range conversations {
					_, err = FolkForm_CreateConversation(page_id, page_username, conversation)
					if err != nil {
						log.Println("Lỗi khi tạo hội thoại:", err)
						continue
					}
				}

				foundAnyConversation = true
				page_number++
				continue
			} else { // Nếu không có dữ liệu thì thoát vòng lặp
				break
			}
		}

		if foundAnyConversation == true {
			until = since
			since = until - 60*60*24*30
			continue
		} else {
			break
		}
	}

	return nil
}

// Hàm bridge_SyncConversationsOfPage_v2 sẽ đồng bộ danh sách hội thoại của trang Facebook từ server Pancake về server FolkForm
// - Lấy danh sách hội thoại của page từ server Pancake
// - Đẩy danh sách hội thoại vào server FolkForm
func bridge_SyncConversationsOfPage_v2(page_id string, page_username string, page_access_token string) (resultErr error) {

	last_conversation_id := ""
	for {
		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		resultGetConversations, err := Pancake_GetConversations_v2(page_id, page_access_token, last_conversation_id)
		if err != nil {
			log.Println("Lỗi khi lấy danh sách hội thoại:", err)
			break
		}

		if resultGetConversations["conversations"] != nil {
			conversations := resultGetConversations["conversations"].([]interface{})
			if len(conversations) > 0 {
				for _, conversation := range conversations {
					_, err = FolkForm_CreateConversation(page_id, page_username, conversation)
					if err != nil {
						log.Println("Lỗi khi tạo hội thoại:", err)
						continue
					}
				}

				last_conversation_id = conversations[len(conversations)-1].(map[string]interface{})["id"].(string)
				continue
			} else {
				break
			}
		} else {
			break
		}
	}

	return nil
}

// Hàm Bridge_SyncConversations sẽ đồng bộ danh sách hội thoại của trang Facebook từ server Pancake về server FolkForm
// - Lấy danh sách trang từ server FolkForm
// - Gọi hàm bridge_SyncConversationsOfPage để đồng bộ hội thoại của từng trang
func Bridge_SyncConversations() (resultErr error) {

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
				for _, item := range items {

					// Dừng nửa giây trước khi tiếp tục
					time.Sleep(100 * time.Millisecond)

					// chuyển item từ interface{} sang dạng map[string]interface{}
					page := item.(map[string]interface{})
					page_id := page["pageId"].(string)
					page_access_token := page["pageAccessToken"].(string)
					page_username := page["pageUsername"].(string)
					is_sync := page["isSync"].(bool)
					if page_access_token != "" && is_sync == true {
						// Gọi hàm bridge_SyncConversationsOfPage để đồng bộ hội thoại của từng trang
						err = bridge_SyncConversationsOfPage_v2(page_id, page_username, page_access_token)
						if err != nil {
							log.Println("Lỗi khi đồng bộ hội thoại:", err)
							continue
						}
					}
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

func bridge_SyncMessageOfConversation(page_access_token string, page_id string, page_username string, conversation_id string, customer_id string) (resultErr error) {

	resultGetMessages, err := Pancake_GetMessages(page_id, page_access_token, conversation_id, customer_id)
	if err != nil {
		return errors.New("Lỗi khi lấy danh sách tin nhắn từ server Pancake")
	}

	_, err = FolkForm_CreateMessage(page_id, page_username, conversation_id, customer_id, resultGetMessages)
	if err != nil {
		return errors.New("Lỗi khi tạo tin nhắn trên server FolkForm")
	}

	return nil

}

func Bridge_SyncMessages() (resultErr error) {

	limit := 50
	page := 0

	for {

		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// Lấy danh sách các Conversations từ server FolkForm
		resultGetConversations, err := FolkForm_GetConversations(page, limit)
		if err != nil {
			return errors.New("Lỗi khi lấy danh sách trang Facebook")
		}

		data := resultGetConversations["data"].(map[string]interface{})
		itemCount := data["itemCount"].(float64)

		if itemCount > 0 {
			items := data["items"].([]interface{})
			if len(items) > 0 {
				for _, item := range items {

					// chuyển item từ interface{} sang dạng map[string]interface{}
					page := item.(map[string]interface{})
					pageId := page["pageId"].(string)
					pageUsername := page["pageUsername"].(string)
					conversationId := page["conversationId"].(string)
					customerId := page["customerId"].(string)

					resultGetPageByPageId, err := FolkForm_GetFbPageByPageId(pageId)
					if err != nil {
						log.Println("Lỗi khi lấy trang theo pageId:", err)
						continue
					}

					data := resultGetPageByPageId["data"].(map[string]interface{})
					page_access_token := data["pageAccessToken"].(string)

					if page_access_token != "" {
						// Gọi hàm bridge_SyncConversationsOfPage để đồng bộ hội thoại của từng trang
						err = bridge_SyncMessageOfConversation(page_access_token, pageId, pageUsername, conversationId, customerId)
						if err != nil {
							log.Println("Lỗi khi đồng bộ tin nhắn:", err)
							continue
						}
					}
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
