package main

import(
	"github.com/coderudu/cotach/Routes"
	
)



func main(){
port := os.getEnv("PORT") 

if port == " "{
	port = 8000
}

router := gin.new()
router.Use(gin.logger())

routes.AuthRoutes(router)
routes.UserRoutes(router)

router.GET("/api-1", func(c *gin.Context)){
	c.JSON(200, gin.H{"success": "you can now access api-1"})
}

router.GET("/api-2", func(c *gin.Context)){
	c.JSON(200, gin.H{"success": "you can now access api-2"})
}

router.Run(":" + port)

}