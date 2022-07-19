package image

import (
	"github.com/skisocks/foto/pkg/helpers/files"
	"github.com/xiam/exif"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Image struct {
	File *os.File
	Info *exif.Data
}

func Read(path string) (*Image, error) {
	img := new(Image)
	err := img.updateData(path)
	if err != nil {
		return nil, err
	}
	defer img.File.Close()

	img.Info, err = exif.Read(path)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (i *Image) Move(dst string) error {
	err := files.Move(i.GetPath(), dst)
	if err != nil {
		return err
	}
	return i.updateData(dst)
}

func (i *Image) Copy(dst string) error {
	err := files.Copy(i.File.Name(), dst)
	if err != nil {
		return err
	}
	return i.updateData(dst)
}

func (i *Image) GetPath() string {
	return i.File.Name()
}

func (i *Image) GetName() string {
	return filepath.Base(i.GetPath())
}

func (i *Image) GetExt() string {
	split := strings.Split(i.GetName(), ".")
	return strings.ToLower("." + split[len(split)-1])
}

func (i *Image) GetDateTime() (time.Time, error) {
	return time.Parse("2006:01:02 15:04:05", i.Info.Tags["Date and Time"])
}

func (i *Image) updateData(path string) (err error) {
	i.File, err = os.Open(path)
	return err
}
