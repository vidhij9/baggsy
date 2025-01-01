package database_test

// import (
// 	"baggsy/backend/database"
// 	"baggsy/backend/models"
// 	"errors"
// 	"testing"
// )

// func TestRunMigrations(t *testing.T) {
// 	// Mock DB.AutoMigrate to avoid actual database operations
// 	originalAutoMigrate := DB.AutoMigrate
// 	defer func() { DB.AutoMigrate = originalAutoMigrate }()
// 	DB.AutoMigrate = func(dst ...interface{}) error {
// 		for _, model := range dst {
// 			switch model.(type) {
// 			case *models.Bag, *models.BagMap, *models.Link:
// 				// Simulate successful migration
// 			default:
// 				return errors.New("unexpected model")
// 			}
// 		}
// 		return nil
// 	}

// 	err := database.RunMigrations()
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}
// }

// func TestRunMigrations_Error(t *testing.T) {
// 	// Mock DB.AutoMigrate to simulate an error
// 	originalAutoMigrate := DB.AutoMigrate
// 	defer func() { DB.AutoMigrate = originalAutoMigrate }()
// 	DB.AutoMigrate = func(dst ...interface{}) error {
// 		return errors.New("migration error")
// 	}

// 	err := database.RunMigrations()
// 	if err == nil {
// 		t.Errorf("Expected error, got nil")
// 	}
// }
