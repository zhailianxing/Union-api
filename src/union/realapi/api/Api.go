package api

import (
	"encoding/base64"

	"union/realapi/data/session"

	"github.com/gocraft/web"
	"tantanbei.com/token"
)

type RootContext struct{}

type Apn struct {
	*RootContext
}

type Web struct {
	*RootContext
}

type Api struct {
	*RootContext

	UserId uint32
}

func (self *Api) ApiSessionMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {

	//log.Debug(" in session middleware url:", r.URL, ", r.Method:", r.Method, ", r.Header:", r.Header)

	tokenStr := r.Header.Get(HEADER_TOKEN)
	if tokenStr == "" {
		//c, err := r.Cookie(HEADER_TOKEN)
		//if err != nil {
		//	FailWithErrCode(w, ERR_CODE_AUTH_FAILED)
		//	return
		//}
		//tokenStr = c.Value
	}

	//log.Debug("in session middleware tokenStr:", tokenStr)

	tokenBytes, e := base64.URLEncoding.DecodeString(tokenStr)

	//length of token byte must be 64
	if e != nil {
		//log.Debug("decode tokenid failed:", e, len(tokenBytes))
		FailWithErrCode(w, ERR_CODE_AUTH_FAILED)
		return
	}

	//validate token
	if !token.Validate(tokenBytes) {
		//log.Debug("token validation failed:", e, len(tokenBytes))
		FailWithErrCode(w, ERR_CODE_AUTH_FAILED)

		//TODO this part is a attack vector on the db...
		//we need to remove invalid tokens from db...fix a bug when token hash has changed and all tokens in db
		//should be invalidated..we need to protect this.
		if len(tokenBytes) == token.TOKEN_SIZE {
			//log.Debug("about to delete an invalid token from db...if exists")
			go session.DeleteSession(tokenBytes)
		}

		return
	}

	//tokenid is last 32 bytes...
	//first 32 bytes is for tokenId forgery protection..
	tokenId := tokenBytes[32:64]

	//check to see if DDOS_KEY is overlimit..
	//if overLimit := context.DDOS_Key_Check(string(tokenId), 100); overLimit {
	//	w.Disconnect()
	//	return
	//}

	userId := session.GetSessionByTokenId(tokenId)

	if userId == 0 {
		//log.Debug("ApiSessionMiddleware check tokenid failed:", tokenId)
		FailWithErrCode(w, ERR_CODE_AUTH_FAILED)
		return
	}

	self.UserId = userId

	//allow possible early gc of this items...next() make take a long time
	tokenStr = ""
	tokenBytes = nil
	tokenId = nil

	next(w, r)
}

func (self *RootContext) Acao(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	next(w, r)
}
