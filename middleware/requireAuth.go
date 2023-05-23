package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
     
	
	"github.com/MLCavalcante/go-jwt/initializers"
	"github.com/MLCavalcante/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// pegar o cookie do request
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// decodificar/ validar
    
	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(os.Getenv("SECRET")), nil
    })

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	    // checar exp

        if float64(time.Now().Unix()) > claims["exp"].(float64){ // n é muito limpo dessa maneira...
	        c.AbortWithStatus(http.StatusUnauthorized)

		}
	    // Achar o user com o token sub
		var user models.User
	    initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
            c.AbortWithStatus(http.StatusUnauthorized)

		}
	    
		//anexar a req
		c.Set("user", user)
	    
		//continuar
	    c.Next()
    } else {
	c.AbortWithStatus(http.StatusUnauthorized)
    }

	
}