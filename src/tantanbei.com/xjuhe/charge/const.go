package charge

/*
CREATE TABLE chargerecord(
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	phone VARCHAR(64) NOT NULL ,
	card VARCHAR(64) NOT NULL ,
	orderid VARCHAR(512) NOT NULL ,
	errorcode INT UNSIGNED NOT NULL ,
	reason VARCHAR(512) NOT NULL ,
	result VARCHAR(512) NOT NULL ,
	PRIMARY KEY (id)
)ENGINE = MyISAM DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
*/

const (
	APP_KEY = "f6afef3dd7b731c0b02d582b9f10dc6f"

	TOP_UP_URL = "http://op.juhe.cn/ofpay/mobile/onlineorder"

	CARD_NUM_10  = "10"
	CARD_NUM_20  = "20"
	CARD_NUM_30  = "30"
	CARD_NUM_50  = "50"
	CARD_NUM_100 = "100"
	CARD_NUM_300 = "300"

	ADD_RECORD_SQL = "INSERT INTO chargerecord SET phone=?, card=?, orderid=?, errorcode=?, reason=?, result=?"
)
