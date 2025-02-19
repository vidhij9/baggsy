package models

import (
	"gorm.io/gorm"
)

// DB is the global database connection (initialized in main).
var DB *gorm.DB
