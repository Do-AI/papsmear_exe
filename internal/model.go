package internal

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

// db와 관련된 코드

// Slide db table
type Slide struct {
	No        int    `gorm:"primary_key;autoIncrement"`
	SlideId   string `gorm:"unique"`
	OrgNo     int
	AgentNo   int
	Extension string
	Synced    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Tile db table
type Tile struct {
	No         int `gorm:"primary_key;autoIncrement"`
	SlideNo    int
	Level      int32
	Coordinate string
	Size       string
	Synced     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ConnectDB 함수는 DB와 연결하고 필요 table을 생성해준다.
func ConnectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(CONFIG.Db.Path), &gorm.Config{PrepareStmt: true})
	//sqlDB, err := db.DB()
	//sqlDB.SetMaxIdleConns(CONFIG.Goroutine)
	//sqlDB.SetMaxOpenConns(CONFIG.Goroutine)
	//sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&Slide{})
	if err != nil {
		return nil
	}
	err = db.AutoMigrate(&Tile{})
	if err != nil {
		return nil
	}
	return db
}

var DB *gorm.DB

// 처음 model.go 파일을 읽어들일 때 db와 연결하여 DB 변수에 갖고있는다.
func init() {
	DB = ConnectDB()
}
