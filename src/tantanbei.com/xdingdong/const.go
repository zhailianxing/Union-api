package xdingdong

import "errors"

const APIKEY = "6c016ff2f803650588d032e8a231530a"
const CONTENT = "【拍牌宝】尊敬的用户，你的验证码是：%s。请勿告诉其他人。"
const URL = "https://api.dingdongcloud.com/v1/sms/sendyzm"
const URL_TZ = "https://api.dingdongcloud.com/v1/sms/sendtz"

var ERROR_INVALID error = errors.New("code is invalid")
var ERROR_SEND_FAILED error = errors.New("request verification failed")
