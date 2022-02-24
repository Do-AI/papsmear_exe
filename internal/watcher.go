package internal

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 컴퓨터 directory의 변화를 읽어 새로 생성된 파일들의 목록을 가져오는 코드.

// Contains 함수는 타겟 문자열(str)이 문자열 집합(s)에 속해있는지 확인하는 함수이다.
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// getTargetFiles 함수는 agent로 전송할 파일의 경로 목록을 가져온다, 만약 파일 이름에 공백이 들어가면 지워준다.
func getTargetFiles() []string {
	var targetFiles []string
	err := filepath.Walk(CONFIG.Folder,
		func(path string, info fs.FileInfo, err error) error {
			newPath := strings.ReplaceAll(path, " ", "")
			_ = os.Rename(path, newPath)

			if Contains(CONFIG.Extensions, filepath.Ext(newPath)) {
				targetFiles = append(targetFiles, newPath)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return targetFiles
}

// GetSystemSlide 함수는 컴퓨터에 저장되어 있는 경로를 뺀 slide 파일이름만 가져온다.
func GetSystemSlide() []string {
	var systemSlides []string
	for _, slide := range getTargetFiles() {
		systemSlides = append(systemSlides, filepath.Base(slide))
	}
	return systemSlides
}

// getDbSlide 함수는 db에 저장되어 있는 slide 파일 이름을 가져온다.
func getDbSlide() []string {
	var dbSlides []Slide
	DB.Find(&dbSlides)
	var dbSlideName []string
	for _, slide := range dbSlides {
		dbSlideName = append(dbSlideName, slide.SlideId+slide.Extension)
	}
	return dbSlideName
}

// getInsertSlide 함수는 db에 아직 들어가지 않은 컴퓨터에 저장된 slide 파일 이름을 가져온다.
func getInsertSlide() []string {
	systemSlideName := GetSystemSlide()
	dbSlideName := getDbSlide()
	var insertSlide []string
	for _, sysSlide := range systemSlideName {
		if !Contains(dbSlideName, sysSlide) {
			insertSlide = append(insertSlide, sysSlide)
		}
	}
	return insertSlide
}

// InsertNonTrackingSlides 함수는 db에 들어가야할 slide를 db에 삽입해주는 함수이다.
func InsertNonTrackingSlides() {
	for _, slide := range getInsertSlide() {
		extension := filepath.Ext(slide)
		baseName := strings.ReplaceAll(slide, extension, "")

		// sqlite 저장
		DB.Create(
			&Slide{
				SlideId:   baseName,
				OrgNo:     CONFIG.Credential.OrgNo,
				AgentNo:   CONFIG.Credential.AgentNo,
				Extension: extension,
				Synced:    0,
			},
		)
	}
}
