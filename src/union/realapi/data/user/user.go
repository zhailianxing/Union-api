package user

import (
	"union/realapi/packet"
	"union/realapi/share"
	"database/sql"
	"tantanbei.com/xcipher"
	"tantanbei.com/xdb"
)

func AddNewUser(phone, password, username, displayName, email string, accType int) (uint32, error) {
	if _, ok := xdb.SelectOne(share.AdvDb, SQL_FIND_USER_BY_USERNAME, username); ok {
		return 0, ERR_ALREADY_REGISTERED
	}

	hashPassword := xcipher.HashPassword(password)

	result, err := share.AdvDb.Exec(SQL_ADD_NEW_USER, phone, hashPassword, username, displayName, email, accType)
	if err != nil {
		panic(err)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return uint32(userId), nil
}

func IsUsernameRegistered(username string) bool {
	if _, ok := xdb.SelectOne(share.AdvDb, SQL_FIND_USER_BY_USERNAME, username); ok {
		return true
	} else {
		return false
	}
}

func CheckUserPassword(username, password string) (session *packet.Session) {
	hashPassword := xcipher.HashPassword(password)

	var result *sql.Rows
	var err error
	session = &packet.Session{}

	if result, err = share.AdvDb.Query(SQL_CHECK_USER_PASSWORD, username, hashPassword); err != nil {
		panic(err)
	} else {
		if result.Next() {
			result.Scan(&session.UserId, &session.UserPhone, &session.UserName)
		} else {
			return nil
		}
	}

	return
}

func ModifyPassword(phone, password string) {
	hashPassword := xcipher.HashPassword(password)
	var err error

	if _, err = share.AdvDb.Exec(SQL_MODIFY_USER_PASSWORD, hashPassword, phone); err != nil {
		panic(err)
	}
}

func GetUserIdByUserName(phone string) (userId uint32) {
	if result, err := share.AdvDb.Query(SQL_FIND_USER_BY_PHONE, phone); err != nil {
		panic(err)
	} else {
		if result.Next() {
			result.Scan(&userId)
		}
	}

	return userId
}
