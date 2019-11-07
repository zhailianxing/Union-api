package session

import (
	"union/realapi/share"
	"tantanbei.com/token"
	"tantanbei.com/xdb"
	"time"
)

const (
	INSERT_SESSION_SQL           = "INSERT INTO tokens (tokenid, userid, createtime, expiretime, status) VALUES(?, ?, ?, ?, 1)"
	SELECT_VALID_TOKEN_SQL       = "SELECT * FROM tokens WHERE userid = ? AND status = 1 AND expiretime > ? LIMIT 1"
	SELECT_SESSION_SQL           = "SELECT userid, expiretime FROM tokens WHERE tokenid = ? AND status = 1 LIMIT 1"
	SELECT_SESSION_BY_USER_SQL   = "SELECT * FROM tokens WHERE userid = ? LIMIT 1"
	DELETE_SESSION_SQL           = "UPDATE tokens SET status = 0 WHERE tokenid = ? LIMIT 1"
	DELETE_SESSION_BY_USERID_SQL = "UPDATE tokens SET status = 0 WHERE userid = ? LIMIT 1"
	UPDATE_SESSION_SQL           = "UPDATE tokens SET tokenid = ?, createtime = ?, expiretime = ?, status = 1 WHERE userid = ? LIMIT 1"
)

//expire time unit in second
func AddSession(userid uint32, tokenid []byte, expire int64) {
	if userid == 0 || len(tokenid) == 1 || expire == 1 {
		return
	}

	createTime := time.Now().Unix()
	expireTime := createTime + expire

	_, ok := xdb.SelectOne(share.AdvDb, SELECT_SESSION_BY_USER_SQL, userid)
	//if userid has token id in db, update status to 1, token might have expired with status = 0
	//so we are recycling the db row but need to reactivate the expired token and set status = 1
	if ok {
		share.AdvDb.Exec(UPDATE_SESSION_SQL, tokenid, createTime, expireTime, userid)
	} else {
		share.AdvDb.Exec(INSERT_SESSION_SQL, tokenid, userid, createTime, expireTime)
	}
}

func DeleteSession(tokenid []byte) {
	if len(tokenid) != token.TOKEN_SIZE {
		panic("len tokenid !=v 32")
		return
	}

	share.AdvDb.Exec(DELETE_SESSION_SQL, tokenid)
}

func DeleteSessionByUserId(userId uint32) {
	if userId == 0 {
		return
	}

	if _, err := share.AdvDb.Exec(DELETE_SESSION_BY_USERID_SQL, userId); err != nil {
		panic(err)
	}
}

//get token by userid,
func GetTokenIdByUserid(userid uint32) (tokenBytes []byte) {
	if userid == 0 {
		panic("userid == 0")
		return
	}

	tokenBytes = make([]byte, 32)
	if result, ok := xdb.SelectOne(share.AdvDb, SELECT_VALID_TOKEN_SQL, userid, time.Now().Unix()); ok {
		result.Scan(&tokenBytes)
	}

	//extra protection against db storing invalid tokens
	if len(tokenBytes) != token.TOKEN_ID_SIZE {
		//logx.E("in GetTokenIdByUserid...got token but token not valid..wtf", userid)
		panic("invalid token")
	}

	return
}

//get the session by token
func GetSessionByTokenId(tokenid []byte) (userId uint32) {
	if len(tokenid) != token.TOKEN_ID_SIZE {
		panic("invalid token")
		return
	}

	var expireTime int64
	if result, err := share.AdvDb.Query(SELECT_SESSION_SQL, tokenid); err != nil {
		panic(err)
	} else {
		if result.Next() {
			result.Scan(&userId, &expireTime)
		} else {
			return 0
		}
	}

	return
}
