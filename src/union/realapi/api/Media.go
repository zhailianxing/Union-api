package api

import (
	"github.com/gocraft/web"
	"io/ioutil"
	"tantanbei.com/xjson"
	"union/realapi/data/media"
	"union/realapi/packet"
	log "github.com/alecthomas/log4go"
	"encoding/json"

)

//func (self *Apn) getOneMeida(rw web.ResponseWriter, req *web.Request) {
//	if req.Method != "POST" {
//		Fail(rw)
//		return
//	}
//
//	okPacket := &packet.OkPacket{}
//
//	body, err := ioutil.ReadAll(req.Body)
//	if err != nil {
//		log.Error("read request error", err)
//		Fail(rw)
//		return
//	}
//	log.Debug("cbody:")
//	log.Debug(string(body))
//
//	bookPacket := &packet.Book{}
//	err = json.Unmarshal(body, &bookPacket)
//	if err != nil {
//		log.Error("decode the json packet error:", err)
//		Fail(rw)
//		return
//	}
//
//	if bookPacket.Id <= 0{
//		okPacket.Code = 0
//		okPacket.Message = "request parms is invalid"
//		rw.Write(xjson.Encode(okPacket))
//		return
//	}
//
//	result := &packet.Book{}
//	if result, err = dbTest.GetBook(bookPacket.Id); err != nil {
//		okPacket.Code = 0
//		okPacket.Message = err.Error()
//		rw.Write(xjson.Encode(okPacket))
//		return
//	}
//	okPacket.Code = 1
//	okPacket.Message = "success"
//	okPacket.Data = result
//	Success(rw)
//	rw.Write(xjson.Encode(okPacket))
//
//}


func (self *Apn) GetAllMeida(rw web.ResponseWriter, req *web.Request) {
	okPacket := &packet.OkPacket{}

	result := make([]*packet.Meida, 0)
	var err error
	if result, err = media.GetMedia(); err != nil {
		okPacket.Code = 0
		okPacket.Message = err.Error()
		rw.Write(xjson.Encode(okPacket))
		return
	}

	okPacket.Code = 1
	okPacket.Message = "success"
	okPacket.Data = result
	Success(rw)
	rw.Write(xjson.Encode(okPacket))

}


func (self *Apn) AddOneMeida(rw web.ResponseWriter, req *web.Request) {
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
	log.Debug("body:")
	log.Debug(string(body))

	meidaPacket := &packet.Meida{}
	err = json.Unmarshal(body, &meidaPacket)
	if err != nil {
		log.Error("decode the json packet error:", err)
		Fail(rw)
		return
	}

	if len(meidaPacket.Name) <= 0{
		okPacket.Code = 0
		okPacket.Message = "request parms is invalid"
		rw.Write(xjson.Encode(okPacket))
		return
	}

	if _, err = media.AddMedia(meidaPacket.Name, meidaPacket.Status); err != nil {
		okPacket.Code = 0
		okPacket.Message = err.Error()
		rw.Write(xjson.Encode(okPacket))
		return
	}
	okPacket.Code = 1
	okPacket.Message = "success"
	Success(rw)
	rw.Write(xjson.Encode(okPacket))

}