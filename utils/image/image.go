package image

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/disintegration/imaging"
)

type ImageType string

const JPG ImageType = "jpg"
const PNG ImageType = "png"
const Unknown ImageType = "unknown"

func IsImage(filename string) bool {
	return strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") || strings.HasSuffix(filename, ".png")
}

func getImageType(filename string) ImageType {
	if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
		return JPG
	} else if strings.HasSuffix(filename, ".png") {
		return PNG
	} else {
		return Unknown
	}
}

func CreateThumbnail(filename string, data []byte) ([]byte, error) {
	imageType := getImageType(filename)
	if imageType == Unknown {
		return nil, fmt.Errorf("unknown file extension, supports .jpg, .jpeg and .png, filename=%s", filename)
	}

	dataBuffer := bytes.NewBuffer(data)
	img, _, err := image.Decode(dataBuffer)
	if err != nil {
		return nil, nil
	}

	thumbnail := imaging.Thumbnail(img, 100, 100, imaging.CatmullRom)
	dst := imaging.New(100, 100, color.NRGBA{0, 0, 0, 0})
	dst = imaging.Paste(dst, thumbnail, image.Pt(0, 0))

	outBuffer := &bytes.Buffer{}

	switch imageType {
	case JPG:
		err = jpeg.Encode(outBuffer, dst, nil)
		if err != nil {
			return nil, fmt.Errorf("error while encoding thumbnail, err=%v", err.Error())
		}
	case PNG:
		err = png.Encode(outBuffer, dst)

		if err != nil {
			return nil, fmt.Errorf("error while encoding thumbnail, err=%v", err.Error())
		}
	default:
		return nil, fmt.Errorf("unknown file extension, support .jpg, .jpeg and .png")
	}

	return outBuffer.Bytes(), nil
}
