package pkg

import (
	"bytes"
	"encoding/base64"
	"image"
	"log"
	"os"

	"github.com/disintegration/imaging"
)

func FromFileToImage(fileName string) (image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func FromImageToString(img image.Image) string {
	var buf bytes.Buffer
	err := imaging.Encode(&buf, img, imaging.PNG)
	if err != nil {
		log.Fatalf("Failed to encode blurred image: %v", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())

}
func FromStringToImage(img string) image.Image {
	imgData, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		log.Fatalf("Failed to decode base64 string: %v", err)
	}
	res, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}
	return res
}
