package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/doai/papsmear/internal"
	"gocv.io/x/gocv"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

const bucket = "dodev-svc-papsmear"
const CredentialPath = "/Users/p829911/Documents/platform-api-credencial.json"

func main12() {
	//slides := internal.GetSystemSlide()
	slides := []string{"LBC1111-20210122(1).bif"}
	for _, slide := range slides {
		log.Println(slide)
		extension := ".bif"
		slideBase := strings.ReplaceAll(slide, extension, "")
		log.Println(slideBase)
		filePath := filepath.Join(internal.CONFIG.Folder, slide)
		log.Println(filePath)
		openSlide := internal.ReadSlide(filePath)

		// slide image save
		thumbnailBuffer := internal.ReadThumbnail(openSlide)
		slideSaveName := getObjectName(slideBase, &internal.TileInfo{}, "slide")
		thumbnailUrl, _ := generateV4PutObjectSignedURL(slideSaveName)
		internal.UploadGcp(thumbnailUrl, thumbnailBuffer)

		maxGoroutines := internal.CONFIG.Goroutine
		guard := make(chan struct{}, maxGoroutines)

		var wg sync.WaitGroup

		// tile image save
		coordinates := internal.MakeCoordinateList(openSlide)
		for _, coordinate := range coordinates {
			guard <- struct{}{}
			go func(coordinate internal.Coordinate) {
				wg.Add(1)
				tileInfo := internal.ReadPatch(0, slideBase, openSlide, coordinate)
				tileBuffer := tileInfo.ImageBuffer
				tileSaveName := getObjectName(slideBase, tileInfo, "tile")
				tileUrl, _ := generateV4PutObjectSignedURL(tileSaveName)
				internal.UploadGcp(tileUrl, tileBuffer)
				<-guard
			}(coordinate)
		}
	}
}

func generateV4PutObjectSignedURL(saveName string) (string, error) {
	jsonKey, err := ioutil.ReadFile(CredentialPath)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}
	opts := &storage.SignedURLOptions{
		Scheme: storage.SigningSchemeV4,
		Method: "PUT",
		Headers: []string{
			"Content-Type:image/jpeg",
		},
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}
	u, err := storage.SignedURL(bucket, saveName, opts)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %v", err)
	}
	return u, nil
}

func uploadFile(image []byte, saveName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(CredentialPath))
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer func(client *storage.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(saveName).NewWriter(ctx)
	if _, err = io.Copy(wc, bytes.NewBuffer(image)); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

func downloadFile(object string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(CredentialPath))
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer func(client *storage.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer func(rc *storage.Reader) {
		err := rc.Close()
		if err != nil {

		}
	}(rc)

	data, err := ioutil.ReadAll(rc)
	decode, _ := gocv.IMDecode(data, gocv.IMReadUnchanged)
	gocv.IMWrite(object, decode)

	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return data, nil
}

func getObjectName(slide string, tileInfo *internal.TileInfo, imageType string) string {
	object := ""

	if imageType == "tile" {
		object = fmt.Sprintf("doai/2021-12/%s/%d/%s", slide, tileInfo.Level, strings.ReplaceAll(tileInfo.Position, ",", "_")+".jpg")
		log.Println(object)
	} else {
		object = fmt.Sprintf("doai/2021-12/%s/thumbnail.jpg", slide)
		log.Println(object)
	}
	return object
}
