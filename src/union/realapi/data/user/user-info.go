package user

import (
	"union/realapi/packet"
	"union/realapi/share"
)

func GetUserById(userId uint32) (user *packet.User) {
	if result, err := share.AdvDb.Query(SQL_SELECT_USER_INFO, userId); err != nil {
		panic(err)
	} else {
		user = &packet.User{}
		if result.Next() {
			result.Scan(&user.UserId, &user.UserPhone, &user.UserName, &user.DataSignup,
				&user.DisplayName, &user.Email, &user.AccType)
		}
	}

	return
}

func ChangeUsername(userId uint32, username string) {
	_, err := share.AdvDb.Exec(SQL_MODIFY_USERNAME, username, userId)
	if err != nil {
		panic(err)
	}
}
