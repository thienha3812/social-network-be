package main

// var profile Profile
// 	var friends Friends
// 	err := db.Where("account_id = ?", request["id"]).Find(&profile).Error
// 	if err != nil {
// 		return nil
// 	}
// 	account_id, err := DecodeToken(c, "id")
// 	db.Where("user_1 = ? ", int(request["id"].(float64))).Where("user_2 =  ? ", account_id).Find(&friends)
// 	db.Where("user_1 = ? ", account_id).Where("user_2 =  ? ", int(request["id"].(float64))).Find(&friends)
// 	response["success"] = true
// 	response["profile"] = echo.Map{
// 		"full_name": profile.FullName,
// 		"live_in":   profile.LiveIn,
// 		"avatar":    profile.Avatar,
// 	}
// 	if friends.Status == 0 {
// 		response["request_status"] = "Thêm bạn bè"
// 		response["request_status_code"] = 0
// 	}
// 	if friends.Status == 1 {
// 		response["request_status"] = "Đang chờ kết bạn"
// 		response["request_status_code"] = 1
// 	}
// 	if friends.Status == 2 {
// 		response["request_status"] = "Bạn bè"
// 		response["request_status_code"] = 2
// 	}
// 	if int(account_id.(float64)) == int(friends.User_1) {
// 		response["self_request"] = true
// 	}
