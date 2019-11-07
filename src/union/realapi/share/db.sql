
CREATE TABLE IF NOT EXISTS user(
    userid INT UNSIGNED NOT NULL AUTO_INCREMENT,
    phone VARCHAR(20) NOT NULL,
    password VARCHAR(512) NOT NULL,
    username VARCHAR(512) NOT NULL,
    datasignup BIGINT UNSIGNED NOT NULL DEFAULT 663004800,
    PRIMARY KEY (userid)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS purchase_intent(
    userid INT UNSIGNED NOT NULL PRIMARY KEY,
    intent_type SMALLINT UNSIGNED NOT NULL,
    produce_type SMALLINT UNSIGNED NOT NULL
);

CREATE TABLE IF NOT EXISTS car_status(
    userid INT UNSIGNED NOT NULL,
    car_type SMALLINT UNSIGNED NOT NULL,
    car_birthday SMALLINT UNSIGNED NOT NULL
);

CREATE TABLE auction_result (
    date INT UNSIGNED PRIMARY KEY NOT NULL,
    limitation INT UNSIGNED NOT NULL DEFAULT 0,
    people_number INT UNSIGNED NOT NULL DEFAULT 0,
    minimum_price INT UNSIGNED NOT NULL DEFAULT 0,
    average_price INT UNSIGNED NOT NULL DEFAULT 0,
    caution_price INT UNSIGNED NOT NULL DEFAULT 0
);

CREATE TABLE user_point (
  userid INT UNSIGNED NOT NULL ,
  point INT UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (userid)
)ENGINE=MyISAM;

CREATE TABLE keep_sign_in(
  userid INT UNSIGNED NOT NULL ,
  keep_day SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  last_signin_time BIGINT UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (userid)
)ENGINE=MyISAM;

CREATE TABLE feed_back(
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  userid INT UNSIGNED NOT NULL ,
  content VARCHAR(512) NOT NULL ,
  drive_message VARCHAR(512) NOT NULL ,
  time BIGINT UNSIGNED NOT NULL DEFAULT 0 ,
  PRIMARY KEY (id)
)ENGINE = MyISAM DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE detail_199101(
  distance INT UNSIGNED NOT NULL ,
  price INT UNSIGNED NOT NULL ,
  PRIMARY KEY (distance)
)ENGINE=MyISAM;

CREATE TABLE images (
  `imageid` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `data` MEDIUMBLOB NOT NULL,
  PRIMARY KEY (`imageid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;

CREATE TABLE image_bind (
  `imageid` INT unsigned NOT NULL,
  `sourceid` INT unsigned NOT NULL,
  `sourcetype` TINYINT unsigned NOT NULL,
  PRIMARY KEY (`sourceid`,`sourcetype`,`imageid`),
  KEY `imageid` (`imageid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;

CREATE TABLE documents (
  `docid` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(125) NOT NULL DEFAULT '',
  `content` LONGBLOB NOT NULL,
  `datesubmit` BIGINT UNSIGNED NOT NULL DEFAULT 663004800,
  PRIMARY KEY (`docid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;

CREATE TABLE papers (
  `paperid` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '',
  `categoryid` SMALLINT UNSIGNED NOT NULL,
  `keywordid1` INT UNSIGNED NOT NULL,
  `keywordid2` INT UNSIGNED,
  `keywordid3` INT UNSIGNED,
  `keywordid4` INT UNSIGNED,
  `datesubmit` BIGINT UNSIGNED NOT NULL DEFAULT 663004800,
  `reprintid` INT UNSIGNED NOT NULL DEFAULT 1,
  PRIMARY KEY (`paperid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE paper_text(
  `paperid` INT UNSIGNED NOT NULL,
  `content` LONGBLOB NOT NULL,
  PRIMARY KEY (`paperid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;

/*
  1: 公告 announcement
  2: 新闻 news
  3: 活动 activity
  4: 文章 article
  5: 攻略 tip
  6: 广告 advertisement
  7: 转载 reproduce
*/
CREATE TABLE categories (
  `categoryid` SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `categorytext` VARCHAR(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`categoryid`, `categorytext`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE reprints(
  `reprintid` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `reprinttext` VARCHAR(128) NOT NULL DEFAULT '拍牌宝',
  PRIMARY KEY (`reprintid`, `reprinttext`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE keywords (
  `keywordid` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `keywordtext` VARCHAR(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`keywordid`, `keywordtext`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `tokens` (
  `tokenid` BINARY(32) NOT NULL,
  `userid` INT UNSIGNED NOT NULL,
  `expiretime` BIGINT UNSIGNED NOT NULL,
  `createtime` BIGINT UNSIGNED NOT NULL,
  `status` TINYINT UNSIGNED NOT NULL,
  PRIMARY KEY `tokenid`(`tokenid`),
  KEY `userid`  (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8, ROW_FORMAT=DYNAMIC;

CREATE TABLE pointproduct (
  `productid` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(512) NOT NULL DEFAULT '',
  `producttype` TINYINT UNSIGNED NOT NULL,
  `value` VARCHAR(128) NOT NULL DEFAULT '',
  `point` INT UNSIGNED NOT NULL DEFAULT 1,
  `stock` INT UNSIGNED NOT NULL DEFAULT 0,
  `expirydate` BIGINT UNSIGNED NOT NULL DEFAULT 663004800,
  PRIMARY KEY (`productid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE videocontent (
  `videoid` INT UNSIGNED NOT NULL,
  `data` LONGBLOB NOT NULL,
  PRIMARY KEY (`videoid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED;

CREATE TABLE video (
  `videoid` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '',
  `summary` varchar(1024) NOT NULL DEFAULT '',
  `userid` INT UNSIGNED NOT NULL DEFAULT 0,
  `keywordid1` INT UNSIGNED NOT NULL,
  `keywordid2` INT UNSIGNED,
  `keywordid3` INT UNSIGNED,
  `keywordid4` INT UNSIGNED,
  `datesubmit` BIGINT UNSIGNED NOT NULL DEFAULT 663004800,
  PRIMARY KEY (`videoid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;