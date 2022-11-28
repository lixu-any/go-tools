package img

import (
	"github.com/disintegration/imaging"
	lfile "github.com/lixu-any/go-tools/file"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
)

// GetImgAttr 获取图片属性
func GetImgAttr(filename string) (width, height int, md5str string, err error) {

	file, err := os.Open(filename)

	if err != nil {
		return
	}

	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		return
	}

	b := img.Bounds()
	width = b.Max.X
	height = b.Max.Y

	md5str, err = lfile.FileMD5(filename)

	return
}

//生成缩略图
func CreateThumbnail(imgname string, w, h int) (sfilename string, err error) {

	//生成缩略图
	shortfile, err := imaging.Open(imgname)
	if err != nil {

		return
	}

	ext := filepath.Ext(imgname)

	//生成缩略图，尺寸150*200，并保持到为文件2.jpg
	dsc := imaging.Resize(shortfile, w, h, imaging.Lanczos)

	tempshortname := "CreateThumbnail.*" + ext

	tempshowfile, err := ioutil.TempFile("", tempshortname)
	if err != nil {
		return
	}

	sfilename = tempshowfile.Name()

	err = imaging.Save(dsc, tempshowfile.Name())
	if err != nil {
		return
	}

	return
}
