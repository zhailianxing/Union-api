package api

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"union/realapi/data/user"
	"union/realapi/packet"

	"union/realapi/security"

	log "github.com/alecthomas/log4go"
	"github.com/gocraft/web"
	"tantanbei.com/token"
	"tantanbei.com/xdingdong"
	"tantanbei.com/xjson"
	"tantanbei.com/xstring"
	"union/realapi/data/session"
)

type SignInUp struct {
	Phone       string `json:"phone"`
	UserName    string `json:"user_name"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	AccType     int    `json:"acc_type"` //0普通用户 1代理商 2管理员 3 代理子账户  4 充值管理员 5 审核管理员 6子管理账号
	Email       string `json:"email"`
}

func (self *SignInUp) IsValid() bool {
	if (len(self.Phone) != 11 && len(self.UserName) == 0) || len(self.Password) > 20 || len(self.Password) < 6 {
		return false
	}
	return true
}

//error code 0: already registered
func (self *Apn) SignUp(rw web.ResponseWriter, req *web.Request) {
	okPacket := &packet.OkPacket{}

	if req.Method != "POST" {
		log.Error("signUp must be post request")
		Fail(rw)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("read request error", err)
		Fail(rw)
		return
	}

	signInUpPacket := &SignInUp{}
	err = json.Unmarshal(body, &signInUpPacket)
	if err != nil {
		okPacket.Code = 0
		okPacket.Message = "请求参数错误"
		rw.Write(xjson.Encode(okPacket))
		log.Error("decode the json packet error:", err)
		return
	}

	if !signInUpPacket.IsValid() {
		okPacket.Code = 0
		okPacket.Message = "请求参数错误"
		rw.Write(xjson.Encode(okPacket))
		return
	}

	var userId uint32
	if userId, err = user.AddNewUser(signInUpPacket.Phone, signInUpPacket.Password, signInUpPacket.UserName, signInUpPacket.DisplayName, signInUpPacket.Email, signInUpPacket.AccType); err != nil {
		log.Error("add new user error:", err)
		if err == user.ERR_ALREADY_REGISTERED {
			Success(rw)
			okPacket.Code = 0
			okPacket.Data = "该用户名已经注册过"
			rw.Write(xjson.Encode(okPacket))
			return
		} else {
			Fail(rw)
			return
		}
	}

	tokenId := session.GetTokenIdByUserid(userId)

	var fullToken []byte

	if tokenId == nil || len(tokenId) == token.TOKEN_ID_SIZE {
		//full 64 byte token
		fullToken = token.Generate()
		session.AddSession(userId, fullToken[32:64], SESSION_EXPIRE_TIME)
	} else {
		fullToken = make([]byte, token.TOKEN_SIZE)

		//if we have tokenid for this user, re-use it..
		copy(fullToken[32:], tokenId)

		//regnerate
		token.RecreateHash(fullToken)
	}

	Success(rw)
	okPacket.Code = 1
	user := &packet.Session{
		UserName:    signInUpPacket.UserName,
		UserId:      userId,
		UserPhone:   signInUpPacket.Phone,
		TokenId:     fullToken,
		Email:       signInUpPacket.Email,
		DisplayName: signInUpPacket.DisplayName,
		AccType:     signInUpPacket.AccType,
	}
	okPacket.Data = user

	uid_cookie := &http.Cookie{
		Name:     HEADER_TOKEN,
		Value:    base64.URLEncoding.EncodeToString(fullToken),
		Path:     "/",
		HttpOnly: false,
		MaxAge:   SESSION_EXPIRE_TIME,
		Expires:  time.Now().Add(SESSION_EXPIRE_TIME),
	}
	http.SetCookie(rw, uid_cookie)
	rw.Write(xjson.Encode(okPacket))
}

type GetCode struct {
	Phone          string
	MustRegistered bool
	Code           string
}

//error code : 0: phone has registered
//             1: get verification code more times
//             2: the phone has not registered
func (self *Apn) GetVerificationCode(rw web.ResponseWriter, req *web.Request) {

	if req.Method != "POST" {
		log.Error("signUp must be post request")
		Fail(rw)
		return
	}

	remoteIP, _ := xstring.CheckIpAndPort(req.RemoteAddr)
	if !security.CheckIpValid(remoteIP, 10) {
		Success(rw)
		rw.Write(xjson.Encode(packet.OkPacket{Code: 0, Message: "get verification code more times"}))
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("read request error", err)
		Fail(rw)
		return
	}

	getCode := &GetCode{}
	err = json.Unmarshal(body, &getCode)
	if err != nil {
		log.Error("decode the json packet error:", err)
		Fail(rw)
		return
	}

	if user.IsUsernameRegistered(getCode.Phone) {
		if !getCode.MustRegistered {
			Success(rw)
			rw.Write(xjson.Encode(packet.OkPacket{Code: 0, Message: "phone has registered"}))
			return
		}
	} else {
		if getCode.MustRegistered {
			Success(rw)
			rw.Write(xjson.Encode(packet.OkPacket{Code: 0, Message: "the phone has not registered"}))
			return
		}
	}

	err = xdingdong.Send(getCode.Phone, getCode.Code)
	if err != nil {
		log.Error("send code error:", err)
		Fail(rw)
		return
	}

	Success(rw)
	rw.Write(xjson.Encode(packet.OkPacket{Code: 1}))
	return
}

//0:phone/password error
func (self *Apn) Login(rw web.ResponseWriter, req *web.Request) {
	okPacket := &packet.OkPacket{}


	if req.Method != "POST" {
		log.Debug("sign in must be post method")
		Fail(rw)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("read request error", err)
		Fail(rw)
		return
	}

	signInUppacket := &SignInUp{}
	err = json.Unmarshal(body, &signInUppacket)
	if err != nil {
		log.Error("decode the json packet error:", err)
		Fail(rw)
		return
	}

	sessionPacket := user.CheckUserPassword(signInUppacket.UserName, signInUppacket.Password)
	if sessionPacket == nil {
		//phone , password error
		Success(rw)
		okPacket.Code = 0
		okPacket.Message = "username or password is wrong!"
		rw.Write(xjson.Encode(okPacket))
		return
	}

	tokenId := session.GetTokenIdByUserid(sessionPacket.UserId)

	var fullToken []byte

	if tokenId == nil || len(tokenId) == token.TOKEN_ID_SIZE {
		//full 64 byte token
		fullToken = token.Generate()
		session.AddSession(sessionPacket.UserId, fullToken[32:64], SESSION_EXPIRE_TIME)
	} else {
		fullToken = make([]byte, token.TOKEN_SIZE)

		//if we have tokenid for this user, re-use it..
		copy(fullToken[32:], tokenId)

		//regnerate
		token.RecreateHash(fullToken)
	}

	sessionPacket.TokenId = fullToken

	u := user.GetUserById(sessionPacket.UserId)
	sessionPacket.AccType = u.AccType
	sessionPacket.DisplayName = u.DisplayName
	sessionPacket.Email = u.Email

	Success(rw)
	okPacket.Code = 1
	okPacket.Data = sessionPacket

	uid_cookie := &http.Cookie{
		Name:     HEADER_TOKEN,
		Value:    base64.URLEncoding.EncodeToString(fullToken),
		Path:     "/",
		HttpOnly: false,
		MaxAge:   SESSION_EXPIRE_TIME,
		Expires:  time.Now().Add(SESSION_EXPIRE_TIME),
	}
	http.SetCookie(rw, uid_cookie)
	rw.Write(xjson.Encode(okPacket))


}

//0: phone is not registered
func (self *Apn) ResetPassword(rw web.ResponseWriter, req *web.Request) {
	if req.Method != "POST" {
		Fail(rw)
		return
	}

	okPacket := &packet.OkPacket{}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("read request error", err)
		Fail(rw)
		return
	}

	modifyPacket := &SignInUp{}
	err = json.Unmarshal(body, &modifyPacket)
	if err != nil {
		log.Error("decode the json packet error:", err)
		Fail(rw)
		return
	}

	if !modifyPacket.IsValid() {
		okPacket.Code = 0
		okPacket.Message = "request parms is invalid"
		rw.Write(xjson.Encode(okPacket))
		return
	}

	if registered := user.IsUsernameRegistered(modifyPacket.UserName); !registered {
		okPacket.Code = 0
		okPacket.Message = "username has not registered"
		rw.Write(xjson.Encode(okPacket))
		return
	}

	user.ModifyPassword(modifyPacket.Phone, modifyPacket.Password)

	userId := user.GetUserIdByUserName(modifyPacket.UserName)
	if userId == 0 {
		Fail(rw)
		return
	}

	//delete the session, because user change the password
	session.DeleteSessionByUserId(userId)

	okPacket.Code = 1
	Success(rw)
	rw.Write(xjson.Encode(okPacket))

}

func (self *Apn) CheckCookie(rw web.ResponseWriter, req *web.Request) {
	okPacket := &packet.OkPacket{}

	cookie, err := req.Cookie(HEADER_TOKEN)
	if err != nil {
		log.Error("read request error", err)
		Fail(rw)
		return
	}
	token,err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Error("base64 DecodeString error", err)
		Fail(rw)
		return
	}
	if (cookie.Expires.Unix()  < time.Now().Unix()){
		okPacket.Code = 0
		okPacket.Message = "cookie expire"
		Success(rw)
		rw.Write(xjson.Encode(okPacket))
		return
	}
	uid  := session.GetSessionByTokenId(token)
	if uid <= 0 {
		okPacket.Code = 0
		okPacket.Message = "cookie invalid"
		Success(rw)
		rw.Write(xjson.Encode(okPacket))
		return
	}
	okPacket.Code = 1
	okPacket.Message = "success"
	Success(rw)
	rw.Write(xjson.Encode(okPacket))

}