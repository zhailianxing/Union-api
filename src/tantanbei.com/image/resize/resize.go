package resize

import (
	"gopkg.in/gographics/imagick.v2/imagick"
)

func init() {
	imagick.Initialize()
}

func Resize(orgImage []byte, width, height uint) []byte {
	var err error

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err = mw.ReadImageBlob(orgImage)
	if err != nil {
		panic(err)
	}

	// Get original logo size
	widthOrignal := mw.GetImageWidth()
	heightOrignal := mw.GetImageHeight()

	widthScale := widthOrignal / width
	heightScale := heightOrignal / height

	var scale uint
	if widthScale > heightScale {
		scale = heightScale
	} else {
		scale = widthScale
	}

	if scale < 2 {
		return mw.GetImageBlob()
	}

	// Calculate half the size
	hWidth := uint(widthOrignal / scale)
	hHeight := uint(heightOrignal / scale)

	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		panic(err)
	}

	// Set the compression quality to 95 (high quality = low compression)
	err = mw.SetImageCompressionQuality(95)
	if err != nil {
		panic(err)
	}

	return mw.GetImageBlob()
}
