package user

import "errors"

const (
	SQL_ADD_NEW_USER          string = "INSERT INTO user(phone, password, username, display_name, email, acc_type) VALUES(?,?,?,?,?,?)"
	SQL_MODIFY_USER_PASSWORD  string = "UPDATE user SET password=? WHERE phone=? LIMIT 1"
	SQL_FIND_USER_BY_PHONE    string = "SELECT id FROM user WHERE phone=? LIMIT 1"
	SQL_FIND_USER_BY_USERNAME string = "SELECT id FROM user WHERE username=? LIMIT 1"
	SQL_CHECK_USER_PASSWORD   string = "SELECT id, phone, username FROM user WHERE username=? AND password=? LIMIT 1"
	SQL_SELECT_USER_INFO      string = "SELECT id, phone, username, create_time, display_name, email, acc_type FROM user WHERE id=? LIMIT 1"
	SQL_MODIFY_USERNAME       string = "UPDATE user SET username = ? WHERE id = ? LIMIT 1"
)

var ERR_ALREADY_REGISTERED error = errors.New("this phone number is already registered")
