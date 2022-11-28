package video

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	lfile "github.com/lixu-any/go-tools/file"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"io/ioutil"
	"os"
	"strconv"
)

// BoxHeader 信息头
type BoxHeader struct {
	Size       uint32
	FourccType [4]byte
	Size64     uint64
}

// GetMP4Duration 获取视频时长
func GetMP4Duration(filePath string) (duration int, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)

	var (
		info      = make([]byte, 0x10)
		boxHeader BoxHeader
		offset    int64 = 0
	)
	// 获取结构偏移
	for {
		_, err = file.ReadAt(info, offset)
		if err != nil {
			return
		}
		boxHeader = getHeaderBoxInfo(info)
		fourccType := getFourccType(boxHeader)
		if fourccType == "moov" {
			break
		}
		// 有一部分mp4 mdat尺寸过大需要特殊处理
		if fourccType == "mdat" {
			if boxHeader.Size == 1 {
				offset += int64(boxHeader.Size64)
				continue
			}
		}
		offset += int64(boxHeader.Size)
	}
	// 获取move结构开头一部分
	moveStartBytes := make([]byte, 0x100)
	_, err = file.ReadAt(moveStartBytes, offset)
	if err != nil {
		return
	}
	// 定义timeScale与Duration偏移
	timeScaleOffset := 0x1C
	durationOffset := 0x20
	timeScale := binary.BigEndian.Uint32(moveStartBytes[timeScaleOffset : timeScaleOffset+4])
	Duration := binary.BigEndian.Uint32(moveStartBytes[durationOffset : durationOffset+4])
	return int(Duration / timeScale), nil
}

// getHeaderBoxInfo 获取头信息
func getHeaderBoxInfo(data []byte) (boxHeader BoxHeader) {
	buf := bytes.NewBuffer(data)
	_ = binary.Read(buf, binary.BigEndian, &boxHeader)
	return
}

// getFourccType 获取信息头类型
func getFourccType(boxHeader BoxHeader) (fourccType string) {
	fourccType = string(boxHeader.FourccType[:])
	return
}

// ResolveTime 将秒转成时分秒格式
func ResolveTime(seconds uint32) string {
	var (
		h, m, s string
	)
	var day = seconds / (24 * 3600)
	hour := (seconds - day*3600*24) / 3600
	minute := (seconds - day*24*3600 - hour*3600) / 60
	second := seconds - day*24*3600 - hour*3600 - minute*60
	h = strconv.Itoa(int(hour))
	if hour < 10 {
		h = "0" + strconv.Itoa(int(hour))
	}
	m = strconv.Itoa(int(minute))
	if minute < 10 {
		m = "0" + strconv.Itoa(int(minute))
	}
	s = strconv.Itoa(int(second))
	if second < 10 {
		s = "0" + strconv.Itoa(int(second))
	}
	return fmt.Sprintf("%s:%s:%s", h, m, s)
}

//获取图片大小
func GetVideoSize(videopath string) (width, height int, md5str string, err error) {
	data, err := ffmpeg_go.Probe(videopath)
	if err != nil {
		return
	}

	type VideoInfo struct {
		Streams []struct {
			CodecType string `json:"codec_type"`
			Width     int
			Height    int
		} `json:"streams"`
	}
	vInfo := &VideoInfo{}
	err = json.Unmarshal([]byte(data), vInfo)
	if err != nil {
		return
	}
	for _, s := range vInfo.Streams {
		if s.CodecType == "video" {
			width = s.Width
			height = s.Height
			break
		}
	}

	md5str, err = lfile.FileMD5(videopath)

	return

}

// CutVideoByTimesForJpg 截取视频的第几帧作为图片保存
func CutVideoByTimesForJpg(videopath string, times, w, h int) (imgpath string, err error) {

	tmpfile, err := ioutil.TempFile("", "CutVideoByTimesForPng.*.jpg")

	if err != nil {
		return
	}

	imgpath = tmpfile.Name()

	tstr := fmt.Sprintf("%d", times)

	outKwArgs := ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}

	if w > 0 && h > 0 {
		outKwArgs["s"] = fmt.Sprintf("%dx%d", w, h)
	}

	err = ffmpeg_go.Input(videopath, ffmpeg_go.KwArgs{"ss": tstr}).Output(imgpath, outKwArgs).OverWriteOutput().ErrorToStdOut().Run()

	return
}
