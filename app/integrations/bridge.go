package integrations

import (
	"errors"
	"log"
	"time"

	"agent_pancake/global"

	"go.mongodb.org/mongo-driver/bson"
)

// ========================================================================================================
// Hàm xử lý logic trên server FolkForm
// ========================================================================================================

// Hàm Bridge_SyncPages(access_token string) sẽ đồng bộ danh sách trang Facebook từ server Pancake về server FolkForm
// - Lấy danh sách trang từ server Pancake
// - Đẩy danh sách trang vào server FolkForm
func bridge_SyncPagesOfAccessToken(access_token string) (resultErr error) {

	log.Println("Đang đồng bộ trang với access token:", access_token)

	// Lấy danh sách trang từ server Pancake
	resultPages, err := PanCake_GetFbPages(access_token)
	if err != nil {
		log.Println("Lỗi khi lấy danh sách trang Facebook:", err)
		return errors.New("Lỗi khi lấy danh sách trang Facebook")
	}

	// Lấy data lưu trong resultPages dạng []interface{} ở categorizedactivated
	activePages := resultPages["categorized"].(map[string]interface{})["activated"].([]interface{})
	for _, page := range activePages {

		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		log.Println("Đang tạo trang trên FolkForm với access token:", access_token)

		FolkForm_CreateFbPage(access_token, page)

	}

	log.Println("Đồng bộ trang với access token thành công:", access_token)

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
			log.Println("Lỗi khi lấy danh sách access token:", err)
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
					log.Println("Đang đồng bộ trang với access token:", access_token)
					bridge_SyncPagesOfAccessToken(access_token)
				}
			}
			page++
			continue
		} else {
			break
		}
	}

	log.Println("Đồng bộ trang Facebook từ server Pancake về server FolkForm thành công")

	return nil
}

// Hàm FolkForm_UpdarePageAccessToken sẽ cập nhật page_access_token của trang Facebook trên server FolkForm bằng cách:
// - Gửi yêu cầu tạo page_access_token lên server PanCake
// - Lấy page_access_token từ phản hồi và cập nhật lên server FolkForm
func Bridge_UpdatePagesAccessToken_toFolkForm() (resultErr error) {

	limit := 50
	page := 0

	for {
		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// Lấy danh sách các pages từ server FolkForm
		resultPages, err := FolkForm_GetFbPages(page, limit)
		if err != nil {
			log.Println("Lỗi khi lấy danh sách trang Facebook:", err)
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

					log.Println("Đang kiểm tra page_access_token cho trang:", page_id)

					// Gọi hàm PanCake_GetFbPages để test page_access_token có hợp lệ không
					_, err := Pancake_GetConversations_v2(page_id, old_page_access_token)
					if err == nil {
						log.Println("Page_access_token vẫn còn hiệu lực cho trang:", page_id)
						continue
					}

					// Gọi hàm PanCake_GeneratePageAccessToken để lấy page_access_token
					resultGeneratePageAccessToken, err := PanCake_GeneratePageAccessToken(page_id, access_token)
					if err != nil {
						log.Println("Lỗi khi lấy page access token:", err)
						continue
					}

					// chuyển resultGeneratePageAccessToken từ interface{} sang dạng map[string]interface{}
					page_access_token := resultGeneratePageAccessToken["page_access_token"].(string)
					// Gọi hàm FolkForm_UpdatePageAccessToken để cập nhật page_access_token
					_, err = FolkForm_UpdatePageAccessToken(page_id, page_access_token)
					if err != nil {
						log.Println("Lỗi khi cập nhật page access token:", err)
						continue
					}
					log.Println("Cập nhật page access token thành công cho trang:", page_id)
				}
			}

			page++
			continue
		} else {
			break
		}
	}

	log.Println("Cập nhật page access token cho tất cả các trang thành công")

	return nil
}

// Hàm Bridge_SyncPagesFolkformToLocal sẽ đồng bộ danh sách trang Facebook từ server FolkForm về server local
// - Lấy danh sách trang từ server FolkForm
// - Đẩy danh sách trang vào server local
func Bridge_SyncPagesFolkformToLocal() (resultErr error) {
	limit := 50
	page := 0

	for {
		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// Lấy danh sách các pages từ server FolkForm
		resultPages, err := FolkForm_GetFbPages(page, limit)
		if err != nil {
			log.Println("Lỗi khi lấy danh sách trang Facebook:", err)
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
						log.Println("Lỗi khi chuyển đổi dữ liệu trang:", err)
						return err
					}

					err = bson.Unmarshal(bsonBytes, &cloudFbPage)
					if err != nil {
						log.Println("Lỗi khi chuyển đổi dữ liệu trang:", err)
						return err
					}

					// Append cloudFbPage to global.PanCake_FbPages
					global.PanCake_FbPages = append(global.PanCake_FbPages, cloudFbPage)
				}
			}
			log.Println("Đồng bộ danh sách trang từ FolkForm về local thành công")
		}

	}

	return nil
}

// Hàm bridge_SyncConversationsOfPage sẽ đồng bộ danh sách hội thoại của trang Facebook từ server Pancake về server FolkForm
// - Lấy danh sách hội thoại của page từ server Pancake
// - Đẩy danh sách hội thoại vào server FolkForm
func bridge_SyncConversationsOfPage(page_id string, page_username string) (resultErr error) {

	last_conversation_id := ""
	for {
		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		resultGetConversations, err := Pancake_GetConversations_v2(page_id, last_conversation_id)
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

				new_last_conversation_id := conversations[len(conversations)-1].(map[string]interface{})["id"].(string)
				if new_last_conversation_id != last_conversation_id {
					last_conversation_id = new_last_conversation_id
					continue
				} else {
					break
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	log.Println("Đồng bộ hội thoại cho trang:", page_id, "thành công")

	return nil
}

// Hàm Bridge_SyncConversations sẽ đồng bộ danh sách hội thoại của trang Facebook từ server Pancake về server FolkForm
// - Lấy danh sách trang từ server FolkForm
// - Gọi hàm bridge_SyncConversationsOfPage để đồng bộ hội thoại của từng trang
func Bridge_SyncConversationsFromCloud() (resultErr error) {

	limit := 50
	page := 0

	for {

		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// Lấy danh sách các pages từ server FolkForm
		resultPages, err := FolkForm_GetFbPages(page, limit)
		if err != nil {
			log.Println("Lỗi khi lấy danh sách trang Facebook:", err)
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
						err = bridge_SyncConversationsOfPage(page_id, page_username)
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

	log.Println("Đồng bộ hội thoại từ server Pancake về server FolkForm thành công")

	return nil
}

// Hàm bridge_SyncMessageOfConversation sẽ đồng bộ danh sách tin nhắn của hội thoại từ server Pancake về server FolkForm
func bridge_SyncMessageOfConversation(page_id string, page_username string, conversation_id string, customer_id string) (resultErr error) {

	resultGetMessages, err := Pancake_GetMessages(page_id, conversation_id, customer_id)
	if err != nil {
		log.Println("Lỗi khi lấy danh sách tin nhắn từ server Pancake:", err)
		return errors.New("Lỗi khi lấy danh sách tin nhắn từ server Pancake")
	}

	_, err = FolkForm_CreateMessage(page_id, page_username, conversation_id, customer_id, resultGetMessages)
	if err != nil {
		log.Println("Lỗi khi tạo tin nhắn trên server FolkForm:", err)
		return errors.New("Lỗi khi tạo tin nhắn trên server FolkForm")
	}

	log.Println("Đồng bộ tin nhắn cho hội thoại:", conversation_id, "thành công")

	return nil

}

// Hàm Bridge_SyncMessages sẽ đồng bộ danh sách tin nhắn của trang Facebook từ server Pancake về server FolkForm
func Bridge_SyncMessages() (resultErr error) {

	limit := 50
	page := 0

	for {

		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// Lấy danh sách các Conversations từ server FolkForm
		resultGetConversations, err := FolkForm_GetConversations(page, limit)
		if err != nil {
			log.Println("Lỗi khi lấy danh sách trang Facebook:", err)
			return errors.New("Lỗi khi lấy danh sách trang Facebook")
		}

		data := resultGetConversations["data"].(map[string]interface{})
		itemCount := data["itemCount"].(float64)

		if itemCount > 0 {
			items := data["items"].([]interface{})

			if len(items) > 0 {
				for _, item := range items {

					// chuyển item từ interface{} sang dạng map[string]interface{}
					conversation := item.(map[string]interface{})
					pageId := conversation["pageId"].(string)
					pageUsername := conversation["pageUsername"].(string)
					conversationId := conversation["conversationId"].(string)
					customerId := conversation["customerId"].(string)

					resultGetPageByPageId, err := FolkForm_GetFbPageByPageId(pageId)
					if err != nil {
						log.Println("Lỗi khi lấy trang theo pageId:", err)
						continue
					}

					data := resultGetPageByPageId["data"].(map[string]interface{})
					page_access_token := data["pageAccessToken"].(string)

					if page_access_token != "" {
						// Gọi hàm bridge_SyncConversationsOfPage để đồng bộ hội thoại của từng trang
						err = bridge_SyncMessageOfConversation(pageId, pageUsername, conversationId, customerId)
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

	log.Println("Đồng bộ tin nhắn từ server Pancake về server FolkForm thành công")

	return nil
}

// ========================================================================================================
// Hàm đồng bộ dữ liệu mới nhất từ server Pancake về server FolkForm của 1 trang Facebook
func Sync_NewMessagesOfPage(page_id string, page_username string) (resultErr error) {

	conversation_id_updated := ""

	// Lấy Conversation mới nhất từ server Folkform
	resultGetConversations, err := FolkForm_GetConversationsWithPageId(0, 1, page_id)
	if err != nil {
		log.Println("Lỗi khi lấy danh sách trang Facebook:", err)
		return errors.New("Lỗi khi lấy danh sách trang Facebook")
	}

	data := resultGetConversations["data"].(map[string]interface{})
	itemCount := data["itemCount"].(float64)

	if itemCount > 0 {
		items := data["items"].([]interface{})
		if len(items) > 0 {
			item := items[0]
			conversation := item.(map[string]interface{})
			conversation_id_updated = conversation["conversationId"].(string)
		}
	}

	last_conversation_id := ""
	// Lấy danh sách conversation mới nhất từ server Pancake cho đến khi gặp conversation_id trùng với conversation_id_updated
	for {

		resultGetConversations, err := Pancake_GetConversations_v2(page_id, last_conversation_id)
		if err != nil {
			log.Println("Lỗi khi lấy danh sách hội thoại:", err)
			break
		}

		if resultGetConversations["conversations"] != nil {
			conversations := resultGetConversations["conversations"].([]interface{})
			if len(conversations) > 0 {
				for _, conversation := range conversations {
					conversation_id := conversation.(map[string]interface{})["id"].(string)
					customerId := conversation.(map[string]interface{})["customer_id"].(string)
					if conversation_id == conversation_id_updated {
						return nil
					}

					_, err = FolkForm_CreateConversation(page_id, page_username, conversation)
					if err != nil {
						log.Println("Lỗi khi tạo hội thoại:", err)
						continue
					}

					// Gọi hàm bridge_SyncConversationsOfPage để đồng bộ hội thoại của từng trang
					err = bridge_SyncMessageOfConversation(page_id, page_username, conversation_id, customerId)
					if err != nil {
						log.Println("Lỗi khi đồng bộ tin nhắn:", err)
						continue
					}

					// Dừng nửa giây trước khi tiếp tục
					time.Sleep(100 * time.Millisecond)
				}

				new_last_conversation_id := conversations[len(conversations)-1].(map[string]interface{})["id"].(string)
				if new_last_conversation_id != last_conversation_id {
					last_conversation_id = new_last_conversation_id
					continue
				} else {
					break
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	log.Println("Đồng bộ tin nhắn mới nhất cho trang:", page_id, "thành công")

	return nil

}

// Hàm Sync_NewMessages sẽ đồng bộ dữ liệu mới nhất từ server Pancake về server FolkForm
func Sync_NewMessagesOfAllPages() (resultErr error) {

	limit := 50
	page := 0

	for {

		// Dừng nửa giây trước khi tiếp tục
		time.Sleep(100 * time.Millisecond)

		// Lấy danh sách các pages từ server FolkForm
		resultPages, err := FolkForm_GetFbPages(page, limit)
		if err != nil {
			log.Println("Lỗi khi lấy danh sách trang Facebook:", err)
			return errors.New("Lỗi khi lấy danh sách trang Facebook")
		}

		data := resultPages["data"].(map[string]interface{})
		itemCount := data["itemCount"].(float64)

		if itemCount > 0 {
			items := data["items"].([]interface{})

			if len(items) > 0 {
				for _, item := range items {

					// chuyển item từ interface{} sang dạng map[string]interface{}
					page := item.(map[string]interface{})
					page_id := page["pageId"].(string)
					page_username := page["pageUsername"].(string)
					is_sync := page["isSync"].(bool)
					if is_sync == true {
						// Gọi hàm Sync_NewMessagesOfPage để đồng bộ tin nhắn của từng trang
						err = Sync_NewMessagesOfPage(page_id, page_username)
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

	log.Println("Đồng bộ tin nhắn mới nhất từ server Pancake về server FolkForm thành công")

	return nil
}
