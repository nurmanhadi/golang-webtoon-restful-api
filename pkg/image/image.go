package image

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/chai2010/webp"
)

func Validate(image multipart.FileHeader) error {
	filenames := []string{".jpg", ".jpeg", ".webp", ".png"}
	ext := strings.ToLower(filepath.Ext(image.Filename))
	if slices.Contains(filenames, ext) {
		return nil
	}
	return errors.New("file most be one of: .jpg, .jpeg, .webp, .png")
}
func CompressToCwebp(obj multipart.FileHeader) (*os.File, error) {
	src, err := obj.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var img image.Image
	switch strings.ToLower(filepath.Ext(obj.Filename)) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(src)
	case ".png":
		img, err = png.Decode(src)
	case ".webp":
		img, err = webp.Decode(src)
	default:
		return nil, errors.New("file most be one of: .jpg, .jpeg, .webp, .png")
	}
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp("", "*.webp")
	if err != nil {
		return nil, err
	}

	if err := webp.Encode(tmpFile, img, &webp.Options{Quality: 75}); err != nil {
		defer tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, err
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		defer tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, err
	}

	return tmpFile, nil
}
