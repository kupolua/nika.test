package files

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

const thumbnailSize = 640

// Thumbnail creates thumbnail images
func Thumbnail(src string, dest string) (err error) {
	exist, err := isExist(dest)
	if err != nil {
		return err
	}
	if exist {
		log.Printf("thumbnail already exist %v", dest)
		return nil
	}

	imgFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("can't open source file %s  error: %v", src, err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("can't decode image file %s  error: %v", imgFile.Name(), err)
	}

	newImg := resize.Thumbnail(thumbnailSize, thumbnailSize, img, resize.MitchellNetravali)

	targetFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("can't create target file %s  errror: %v", targetFile.Name(), err)
	}
	defer targetFile.Close()

	err = jpeg.Encode(targetFile, newImg, nil)
	if err != nil {
		return fmt.Errorf("cant' write image to file %s,  error: %v", targetFile.Name(), err)
	}
	return nil
}

func isExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
