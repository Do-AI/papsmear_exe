package internal

import (
	"fmt"
	"github.com/doai/papsmear/config"
	"github.com/doai/papsmear/pkg"
	"gocv.io/x/gocv"
	"log"
	"math"
)

// OpenSlide를 이용해 Tile 정보를 생성하는 코드
// Agent 프로그램에서 가장 핵심적인 로직을 담당한다.

// Coordinate 는 tile level, 좌표 정보를 담고 있는 구조체이다.
type Coordinate struct {
	Level int32
	X     int64
	Y     int64
	W     int64
	H     int64
}

// SlideInfo 는 Slide 정보를 담고 있는 구조체이다.
type SlideInfo struct {
	OrgNo     int
	AgentNo   int
	SlideId   string
	Size      string
	Thumbnail []byte
}

// TileInfo 는 Tile 정보를 담고 있는 구조체이다.
type TileInfo struct {
	SlideNo     int
	SlideId     string
	ImageBuffer *gocv.NativeByteBuffer
	Level       int32
	Position    string
	Size        string
}

// ReadSlide 함수는 슬라이드 파일 경로를 받아 OpenSlide 객체로 변환해주는 함수이다.
func ReadSlide(path string) *pkg.Slide {
	slide, err := pkg.Open(path)
	if err != nil {
		message := fmt.Sprintf("```path: %s\n"+
			"Can't read slide using openslide\n"+
			"Check the file extensions.```", path)
		sendSlackMessage(message)
		return nil
	}
	return &slide
}

var CONFIG = config.Config

// MakeCoordinateList 함수는 CONFIG 변수에 저장되어있는 정보를 바탕으로
// Slide 정보에서 Tile 정보를 생성해내는 역할을 한다.
func MakeCoordinateList(slide *pkg.Slide) []Coordinate {
	var coordinateList []Coordinate
	for _, level := range CONFIG.Patch.Level {
		width, height := slide.LevelDimensions(0)

		var interval int64
		// 아래 코드는 OpenSlide 라이브러리의 로직을 이해하고 있어야한다.
		// OpenSlide는 Slide를 level에 따라 읽어들일 수 있는데
		// level 0 은 level 1 보다 가로가 2배, 세로가 2배인 이미지이다.
		// level 0 기준의 좌표를 저장하고 level 1의 이미지를 가져와야 하기 때문에 아래와 같은 처리를 해준다.
		// 실제 ReadPatch 함수에서 level 1의 이미지를 가져오기 위해 divider로 나눠주는 부분을 볼 수 있다.
		interval = CONFIG.Patch.SizeFront * int64(math.Pow(2, float64(level)))

		// 아래 코드는 tile 정보를 overlap 하게 가져오는 역할을 한다.
		// 지금은 overlap 하지 않은 tile 들을 가져오지만, 추후에 변경될 수도 있기에 남겨놓는다.
		//if level == 0 {
		//	interval = CONFIG.Patch.Size - CONFIG.Patch.Overlap
		//} else {
		//	interval = CONFIG.Patch.SizeFront * int64(math.Pow(2, float64(level)))
		//}

		for w := int64(0); w < width; w += interval {
			for h := int64(0); h < height; h += interval {
				xSize, ySize := interval, interval
				if w+interval > width {
					xSize = width - w
				}
				if h+interval > height {
					ySize = height - h
				}
				coordinateList = append(
					coordinateList,
					Coordinate{level, w, h, xSize, ySize},
				)
			}
		}
	}

	return coordinateList
}

// imgToByte 함수는 이미지를 전송하기 위해 byte 객체로 만들어 주는 역할을 한다.
func imgToByte(slide *pkg.Slide, x int64, y int64, level int32, w int64, h int64) *gocv.NativeByteBuffer {
	byteImg, _ := slide.ReadRegion(x, y, level, w, h)
	matImg, _ := gocv.NewMatFromBytes(int(h), int(w), gocv.MatTypeCV8UC4, byteImg)
	jpegImg := gocv.NewMat()
	gocv.CvtColor(matImg, &jpegImg, gocv.ColorBGRAToBGR)
	buffer, _ := gocv.IMEncodeWithParams(gocv.JPEGFileExt, jpegImg, []int{gocv.IMWriteJpegQuality, 70})

	defer func() {
		err := matImg.Close()
		if err != nil {
			log.Println("matImg.Close() error")
		}
		err = jpegImg.Close()
		if err != nil {
			log.Println("jpegImg.Close() error")
		}
	}()

	return buffer
}

// ReadThumbnail 함수는 slide에서 thumbnail 이미지를 가져오는 역할을 한다.
func ReadThumbnail(slide *pkg.Slide) *gocv.NativeByteBuffer {
	w, h := slide.LevelDimensions(6)
	return imgToByte(slide, 0, 0, 6, w, h)
}

// ReadPatch 함수는 tile 좌표를 바탕으로 slide에서 tile 정보를 읽어오는 역할을 한다.
func ReadPatch(slideNo int, slideId string, slide *pkg.Slide, info Coordinate) *TileInfo {
	divider := int64(math.Pow(2, float64(info.Level)))
	x := info.X / divider
	y := info.Y / divider
	w := info.W / divider
	h := info.H / divider
	imgBuff := imgToByte(slide, info.X, info.Y, info.Level, w, h)
	position := fmt.Sprintf("%d,%d", x, y)
	size := fmt.Sprintf("%d,%d", w, h)

	return &TileInfo{slideNo, slideId, imgBuff, info.Level, position, size}
}
