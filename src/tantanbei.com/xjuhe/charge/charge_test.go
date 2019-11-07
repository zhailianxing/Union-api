package charge

import (
	"fmt"
	"tantanbei.com/xjson"
	"testing"
	"time"
)

func TestTopUp(t *testing.T) {
	TopUp("18117541072", CARD_NUM_10)
}

func TestDecode(t *testing.T) {
	str := `{
    "reason": "允许充值的手机号码及金额",
    "result": null,
    "error_code": 0
}`
	packet := &ResultPacket{}
	xjson.Decode([]byte(str), &packet)

	t.Log(packet)
}

func TestDecode2(t *testing.T) {
	str := `{"reason":"订单提交成功，等待充值","result":{"cardid":"10012","cardnum":"1","ordercash":10.08,"cardname":"安徽电信话费10元","sporder_id":"J17011023535310201871460","uorderid":"HFCZ1484063633076263457","game_userid":"18056561786","game_state":"0"},"error_code":0}`
	packet := &ResultPacket{}
	xjson.Decode([]byte(str), &packet)

	t.Log(packet)

	AddRecord("12345678901", "100", "HFCZ", fmt.Sprintf("%v", packet.Reason), fmt.Sprintf("%v", packet.Result), packet.ErrorCode)

	<-time.After(time.Second)
}
