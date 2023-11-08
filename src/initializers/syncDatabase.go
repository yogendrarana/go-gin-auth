package initializers

import (
	"go-gin-auth/src/models"
)

func SyncDatabase() {
	db.AutoMigrate(&models.User{})
}
