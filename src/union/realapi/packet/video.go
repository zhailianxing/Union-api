package packet

type UploadVideo struct {
	Title     string
	Summary   string
	Thumbnail []byte
	Video     []byte
}

type VideoInfo struct {
	VideoId    uint32
	Title      string
	Summary    string
	UserId     uint32
	UserName   string
	ImageId    uint32
	DataSubmit int64
}

type VideoInfos struct {
	Data []*VideoInfo
}
