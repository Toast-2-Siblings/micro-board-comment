package database

import (
	"os"
	"testing"

	"github.com/Toast-2-Siblings/micro-board-comment/models"
)

func TestMain(m *testing.M) {
	os.Setenv("Mode", "development")
}

// 데이터베이스 초기화 및 마이그레이션테스트 함수
func TestInitDatabaseAndMigration(t *testing.T) {
	err := InitDatabase()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	db := GetDB()
	if db == nil {
		t.Fatal("Expected initialized DB, got nil")
	}

	err = Migrate()
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// 간단한 CRUD 테스트 (테이블 생성 확인 목적)
	comment := models.Comment{
		Content: "test comment",
	}
	result := db.Create(&comment)
	if result.Error != nil {
		t.Errorf("Failed to insert comment: %v", result.Error)
	}

	var count int64
	db.Model(&models.Comment{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 comment, got %d", count)
	}
}
