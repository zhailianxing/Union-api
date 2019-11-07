package api

import (
	"github.com/gocraft/web"
	"tantanbei.com/xjson"
	"union/realapi/packet"
)

func (self *Api) TestGet(rw web.ResponseWriter, req *web.Request) {

	okPacket := &packet.OkPacket{}
	okPacket.Code = 1
	okPacket.Data = self.UserId
	Success(rw)
	rw.Write(xjson.Encode(okPacket))
}

func (self *Api) TestPost(rw web.ResponseWriter, req *web.Request) {
	okPacket := &packet.OkPacket{}
	okPacket.Code = 1
	okPacket.Data = self.UserId
	Success(rw)
	rw.Write(xjson.Encode(okPacket))
}





