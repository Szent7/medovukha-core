package main

import (
	"medovukha/api/rest/middlewares"
	restapi "medovukha/api/rest/v1"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// wsHub := websockets.NewHub()
	// go wsHub.Run()

	router.Static("/_app/immutable/", "./build/_app/immutable/")
	router.NoRoute(func(c *gin.Context) {
		c.File("./build/index.html")
	})

	//! dev headers
	router.Use(middlewares.CORSMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// ws := router.Group("/ws")
	// {
	// 	ws.GET("/containerEvents", websockets.WsHandler(wsHub))
	// }

	rest := router.Group("/rest")
	{
		v1 := rest.Group("/v1")
		{
			//Containers
			v1.POST("/createTest", restapi.CreateTestContainer)
			v1.GET("/getContainerList", restapi.GetContainerList)
			v1.POST("/pauseContainerByid", restapi.PauseContainerByID)
			v1.POST("/unpauseContainerById", restapi.UnpauseContainerByID)
			v1.POST("/killContainerById", restapi.KillContainerByID)
			v1.POST("/startContainerById", restapi.StartContainerByID)
			v1.POST("/stopContainerById", restapi.StopContainerByID)
			v1.POST("/restartContainerById", restapi.RestartContainerByID)
			v1.POST("/removeContainerById", restapi.RemoveContainerByID)
			//Images
			v1.GET("/getImageList", restapi.GetImageList)
			//v1.POST("/buildImageByRepo", restapi.BuildImageByRepo)
			//Networks
			v1.GET("/getNetworkList", restapi.GetNetworkList)
			//Volumes
			v1.GET("/getVolumeList", restapi.GetVolumeList)
			//Deploy
			v1.POST("/createFromGit", restapi.CreateFromGit)
		}
	}

	router.Run("0.0.0.0:10015")
}
