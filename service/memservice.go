package service

import (
	"andrewj.com/readmems/rosco"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// StartService starts an http listener
func StartService(mems *rosco.Mems, config *rosco.ReadmemsConfig) {
	// Instantiate a new router
	r := httprouter.New()

	// Get a MemsDataController instance
	c := NewMemsDataController(mems)

	// Get mems data
	r.GET("/mems", c.GetMemsData)
	r.GET("/mems/:id", c.GetMemsData)
	r.GET("/exit", c.Exit)

	port := config.WebPort
	// Fire up the server
	http.ListenAndServe("localhost:"+port, r)
}
