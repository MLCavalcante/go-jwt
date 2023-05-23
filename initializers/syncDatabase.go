package initializers

import "github.com/MLCavalcante/go-jwt/models"


func SyncDatabase(){
   DB.AutoMigrate(&models.User{})
}