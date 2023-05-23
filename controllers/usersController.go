package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/MLCavalcante/go-jwt/initializers"
	"github.com/MLCavalcante/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
   //Aqui nós vamos pegar o email e password 
   
   var body struct {   //criar uma var que guarda os dados que vão entrar 
       Email    string
	   Password string 
    } 
     //para popular essa var com esses dados
    
	 if c.Bind(&body) != nil{    
	   c.JSON(http.StatusBadRequest, gin.H{
		    "error":"Failed to read body",
	   })
	   return    
	   //A função c.Bind(&body) retorna um erro (não nulo) se não for  possível popular a estrutura body com os dados do request.
	}

   //Hash da senha
    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
		    "error":"Failed to hash Password",
	   })
	   return
	}

   //Criar user 
    user := models.User{Email: body.Email, Password: string(hash)}
    result := initializers.DB.Create(&user) 

	if result.Error != nil {
       c.JSON(http.StatusBadRequest, gin.H{
		  "error":"Failed to create user",
	   })
	   return
	}

   //Responder
   c.JSON(http.StatusOK, gin.H{})

}

func Login(c *gin.Context) {
	//Pegar o email e o password da req
	 var body struct { 
	   Email    string 
	   Password string
    } 
    
    
	 if c.Bind(&body) != nil{    
	   c.JSON(http.StatusBadRequest, gin.H{
		    "error":"Failed to read body",
	   })
	   return    
	   
	}

	//Procurar o user da req
    var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
		    "error":"Invalid email or password",
	    })
	    return 
	}

	//Comparar a senha enviada com a senha salva do hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
		    "error":"Invalid email or password",
	    })
		return
	}

	//Gerar o token jtw 
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"sub": user.ID,
	"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
    })

   
    tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

    if err != nil {
         c.JSON(http.StatusBadRequest, gin.H{
		    "error":"Failed to create token",
	    })
		return
	}
	
	
	//Resposta 
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "", "", false, true)  // criar cookie
	c.JSON(http.StatusOK, gin.H{
       //"token": tokenString,  poderiamos mandar assim mas cookies são melhores 
	})

}

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":"Logei",
	})
}