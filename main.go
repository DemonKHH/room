package main

import (
	"room/router"
)

func main() {
	// r := gin.Default()
	// r.Use(middleware.Cors())
	// router.LoadRoutes(r)
	// r.Run(":8090")
	router.InitRoutes()
}
