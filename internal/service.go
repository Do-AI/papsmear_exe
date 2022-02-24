package internal

import (
	"fmt"
	"github.com/doai/papsmear/pkg"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 실제 사용하는 서비스 관련 로직을 담고있는 코드

// SendSlideInfo 함수는 Slide 관련 정보를 DTO에 담아 spring papsmear 서버로 전송해주는 역할을 한다.
func SendSlideInfo(name string, size string, slide *pkg.Slide) {

	imgBuff := ReadThumbnail(slide)

	slideInfo := struct {
		OrgNo   int    `json:"orgNo"`
		AgentNo int    `json:"agentNo"`
		SlideId string `json:"slideId"`
		Size    string `json:"size"`
	}{
		OrgNo:   CONFIG.Credential.OrgNo,
		AgentNo: CONFIG.Credential.AgentNo,
		SlideId: name,
		Size:    size,
	}

	signedUrl := postServer(slideInfo, "slide")

	UploadGcp(signedUrl, imgBuff)
}

// SendTileInfo 함수는 Tile 관련 정보를 DTO에 담아 spring papsmear 서버로 전송해주는 역할을 한다.
// Tile 개수가 많으므로 goroutine을 통해 병렬처리로 진행한다.
func SendTileInfo(slideNo int, slideId string, slide *pkg.Slide, info Coordinate, wg *sync.WaitGroup) {
	defer wg.Done()

	tileInfo := ReadPatch(slideNo, slideId, slide, info)
	imgBuff := tileInfo.ImageBuffer
	sendInfo := struct {
		SlideNo  int    `json:"slideNo"`
		SlideId  string `json:"slideId"`
		Level    int32  `json:"level"`
		Position string `json:"position"`
		Size     string `json:"size"`
	}{
		SlideNo:  tileInfo.SlideNo,
		SlideId:  tileInfo.SlideId,
		Level:    info.Level,
		Position: tileInfo.Position,
		Size:     tileInfo.Size,
	}
	coordinate := fmt.Sprintf("%d,%d", info.X, info.Y)
	signedUrl := postServer(sendInfo, "tile")

	UploadGcp(signedUrl, imgBuff)

	var tile Tile
	DB.Where("slide_no = ? AND level = ? AND coordinate = ?", slideNo, info.Level, coordinate).Find(&tile)
	// sqlite synced 1 저장
	DB.Model(&tile).Where("slide_no = ? AND level = ? AND coordinate = ?", slideNo, info.Level, coordinate).
		Update("synced", 1)
}

// SendSlideReady 함수는 Slide의 모든 Tile이 전송되었을 때 ready sign을 spring papsmear 서버로 보내준다.
func SendSlideReady(name string) {

	slideInfo := struct {
		OrgNo   int    `json:"orgNo"`
		AgentNo int    `json:"agentNo"`
		SlideId string `json:"slideId"`
		Size    string `json:"size"`
	}{
		OrgNo:   CONFIG.Credential.OrgNo,
		AgentNo: CONFIG.Credential.AgentNo,
		SlideId: name,
		Size:    "",
	}

	postServer(slideInfo, "slide/ready")
}

// GetNonSyncedSlide 함수는 spring papsmear 서버와 동기화 되지 않은 Slide 파일을 읽어온다.
func GetNonSyncedSlide() (int, string, *pkg.Slide, []Coordinate) {
	var slide Slide
	var tiles []Tile

	DB.Limit(1).Where("Synced = 0").Find(&slide)
	if slide.No == 0 {
		return 0, "", nil, nil
	}
	filePath := filepath.Join(CONFIG.Folder, slide.SlideId+slide.Extension)
	openSlide := ReadSlide(filePath)

	// 오류 slide 있는 경우
	DB.Where("Synced = 0").Find(&tiles)

	if len(tiles) != 0 {
		var errorSlide Slide
		DB.Where("no = ?", tiles[0].SlideNo).First(&errorSlide)
		var coordinates []Coordinate
		for _, tile := range tiles {
			xy := strings.Split(tile.Coordinate, ",")
			wh := strings.Split(tile.Size, ",")
			x, _ := strconv.Atoi(xy[0])
			y, _ := strconv.Atoi(xy[1])
			w, _ := strconv.Atoi(wh[0])
			h, _ := strconv.Atoi(wh[1])
			info := Coordinate{tile.Level, int64(x), int64(y), int64(w), int64(h)}
			coordinates = append(coordinates, info)
		}
		return errorSlide.No, errorSlide.SlideId, openSlide, coordinates
	}

	width, height := openSlide.LevelDimensions(0)
	slideSize := fmt.Sprintf("%d,%d", width, height)
	coordinates := MakeCoordinateList(openSlide)

	// sqlite에 tile 정보 저장
	for _, info := range coordinates {
		DB.Create(
			&Tile{
				SlideNo:    slide.No,
				Level:      info.Level,
				Coordinate: fmt.Sprintf("%d,%d", info.X, info.Y),
				Size:       fmt.Sprintf("%d,%d", info.W, info.H),
				Synced:     0,
			},
		)
	}

	SendSlideInfo(slide.SlideId, slideSize, openSlide)

	return slide.No, slide.SlideId, openSlide, coordinates
}

// SendTileService 함수는 spring papsmear 서버와 동기화 되지 않은 slide정보를 가져와
// 그 slide의 tile 정보를 전처리 한 후 전송해준다.
func SendTileService() {
	slideNo, slideId, slide, coordinates := GetNonSyncedSlide()

	defer func(slide *pkg.Slide) {
		if slide != nil {
			slide.Close()
		}
	}(slide)

	if slideNo == 0 && slide == nil {
		log.Println("There are no new slides.")
		time.Sleep(30 * time.Second)
		return
	}

	maxGoroutines := CONFIG.Goroutine
	guard := make(chan struct{}, maxGoroutines)

	var wg sync.WaitGroup

	for _, info := range coordinates {
		guard <- struct{}{}
		go func(info Coordinate) {
			wg.Add(1)
			SendTileInfo(slideNo, slideId, slide, info, &wg)
			<-guard
		}(info)
	}

	wg.Wait()

	// slide synced 완료 되었다고 쿼리 날리기
	var slideRow Slide
	DB.Model(&slideRow).Where("no = ?", slideNo).
		Update("synced", 1)

	// slide db에 잘 저장되었는지 가져오기
	var doneSlideName Slide
	DB.Select("slide_id").Where("no = ?", slideNo).Find(&doneSlideName)
	log.Println(doneSlideName.SlideId + " is done")

	// spring papsmear 서버에 slide 전송 끝났다고 보내기
	SendSlideReady(doneSlideName.SlideId)
}
