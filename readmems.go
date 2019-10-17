package main

import (
	"andrewj.com/readmems/config"
	"andrewj.com/readmems/rosco"
	"andrewj.com/readmems/service"
	"os"
)

// fileExists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func main() {
	// use if the readmems config is supplied
	var readmemsConfig = config.ReadConfig()

	// if argument is supplied then use that as the port id
	if len(os.Args) > 1 {
		readmemsConfig.Port = os.Args[1]
	}

	mems := rosco.New()

	rosco.MemsConnect(mems, readmemsConfig)
	rosco.MemsInitialise(mems)

	// start http service
	service.StartService(mems)
}
