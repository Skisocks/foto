package image

import (
	"github.com/rwcarlsen/goexif/exif"
	"github.com/skisocks/foto/pkg/helpers/files"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	Data     *os.File
	MetaData *exif.Exif
}

func Read(path string) (*Image, error) {
	img := new(Image)
	err := img.updateData(path)
	if err != nil {
		return nil, err
	}
	defer img.Data.Close()
	img.MetaData, err = exif.Decode(img.Data)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (i *Image) Move(dst string) error {
	dst = filepath.Join(dst, i.GetName())
	err := files.Move(i.GetPath(), dst)
	if err != nil {
		return err
	}
	return i.updateData(dst)
}

func (i *Image) Copy(dst string) error {
	dst = filepath.Join(dst, i.GetName())
	err := files.Copy(i.Data.Name(), dst)
	if err != nil {
		return err
	}
	return i.updateData(dst)
}

func (i *Image) GetPath() string {
	return i.Data.Name()
}

func (i *Image) GetName() string {
	return filepath.Base(i.GetPath())
}

func (i *Image) GetExt() string {
	split := strings.Split(i.GetName(), ".")
	return strings.ToLower("." + split[len(split)-1])
}

func (i *Image) updateData(path string) (err error) {
	i.Data, err = os.Open(path)
	return err
}
