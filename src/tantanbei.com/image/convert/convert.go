package convert

import (
	"gopkg.in/gographics/imagick.v2/imagick"
)

//output types
const (
	OUTPUT_IMAGE = 0
	OUTPUT_JPEG  = 1
	OUTPUT_PNG   = 2
	OUTPUT_WEBP  = 3
)

func init() {
	imagick.Initialize()
}

func ConvertImage(value []byte, outputType int) (result []byte) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	//do not call this as library is using finalizer to release C memory...
	// only call manually if not using constructor defer mw.Destroy()

	//handle value

	if e := mw.ReadImageBlob(value); e != nil {
		panic(e)
	}

	format := mw.GetImageFormat()

	switch outputType {
	case OUTPUT_IMAGE:
		//		logx.D("should panic becasue of using imagick package")
		panic("Now use imagick package, not support this type")

	case OUTPUT_JPEG:
		if format == "JPEG" {
			return value
		}

		mw.SetImageFormat("JPEG")
		result = mw.GetImageBlob()

	case OUTPUT_PNG:
		if format == "PNG" {
			return value
		}

		mw.SetImageFormat("PNG")
		result = mw.GetImageBlob()

	case OUTPUT_WEBP:
		if format == "WEBP" {
			return value
		}

		mw.SetImageFormat("WEBP")
		//android < 4.3 has broken webp alpha channel support
		mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_DEACTIVATE)
		result = mw.GetImageBlob()
	default:
		//		logx.D("illegal out type")
		panic("illegal out type")
	}

	return result
}
