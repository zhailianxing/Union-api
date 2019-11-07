package packet

type Meida struct {
	Id      uint32 `json:"id"`
	Name    string `json:"name"`
	Status  int32 `json:"status"`
	UserId  uint32 `json:"user_id"`
	Type    uint32 `json:"type"`
	Config  string `json:"config"`
	ExpConfig string `json:"exp_config"`
	Experiment string `json:"experiment"`
	Comment string `json:"comment"`
	CreateTime string `json:"create_time"` // TODO 时间戳
	ModifiedTime string `json:"modified_time"`
	MediaClass string `json:"media_class"`
	DownloadUrl string `json:"download_url"`
	CategoryIds string `json:"category_ids"`
	Rate float32 `json:"rate"`
	IsRate int32 `json:"is_rate"`
	CityClassBlackWhite string `json:"city_class_black_white"`
	AccClass string `json:"acc_class"`
	BlackClass string `json:"black_class"`
	ApkId int32 `json:"apk_id"`
	DomainId int32 `json:"domain_id"`
	WhiteCityMaterialLevel string `json:"white_city_material_level"`
	MediaBdRateMax int32 `json:"media_bd_rate_max"`
	MediaMdRateMin int32 `json:"media_bd_rate_min"`
}
//
//`id` int(11) unsigned NOT NULL COMMENT '媒体ID',
//`name` varchar(64) COLLATE utf8_unicode_ci NOT NULL COMMENT '媒体名称',
//`status` int(11) NOT NULL COMMENT '0为有效，非0无效，1管理员禁用，2反作弊单日临时禁用',
//`user_id` int(10) unsigned NOT NULL COMMENT '关联用户ID',
//`type` int(10) unsigned NOT NULL COMMENT '0: 应用, 1: 移动网页',
//`config` text COLLATE utf8_unicode_ci NOT NULL COMMENT '媒体定制配置',
//`exp_config` text COLLATE utf8_unicode_ci COMMENT '实验配置',
//`experiment` text COLLATE utf8_unicode_ci COMMENT '实验配置',
//`comment` varchar(512) COLLATE utf8_unicode_ci NOT NULL,
//`create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//`modified_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
//`media_class` varchar(512) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '媒体分类 多个用逗号隔开',
//`download_url` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '下载地址',
//`category_ids` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '媒体类别如下，可多选,用英文逗号分隔。',
//`rate` decimal(11,2) NOT NULL DEFAULT '0.00' COMMENT '在is_rate =1 时，比例值',
//`is_rate` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否是按比例分成0 否，1 是',
//
//
//`city_class_black_white` text COLLATE utf8_unicode_ci NOT NULL COMMENT '城市下的行业类型黑白名单',
//`acc_class` text COLLATE utf8_unicode_ci NOT NULL COMMENT '广告位可接受的广告分类，数字用逗号隔开',
//`black_class` text COLLATE utf8_unicode_ci NOT NULL COMMENT '广告位不可接受的广告分类，数字用逗号隔开，黑名单',
//`apk_id` int(11) NOT NULL DEFAULT '0' COMMENT 'type =1 下载地址',
//`domain_id` int(11) NOT NULL DEFAULT '0' COMMENT '关联的域名',
//`white_city_material_level` varchar(1000) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '素材级别城市定向\r\nmaterial_level\r\n  1,2, 3,4 \r\n  很正规,无敏感词擦边球,有少量敏感词,尺度较大\r\n城市级别：\r\n	1,2,3,4,5,6\r\n	一线城市，新一线城市，二线城市，三线城市， 四线城市，五线城市\r\n',
//`media_bd_rate_max` int(11) NOT NULL DEFAULT '0' COMMENT '媒体BD分成最大值',
//`media_bd_rate_min` int(11) NOT NULL DEFAULT '0' COMMENT '媒体BD分成最小值',