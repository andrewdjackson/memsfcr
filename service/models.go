package service

import (
	"andrewj.com/readmems/rosco"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var mems *rosco.Mems

// MemsDataController struct
type (
	MemsDataController struct{}
)

// NewMemsDataController creates a new mems data controller
func NewMemsDataController(m *rosco.Mems) *MemsDataController {
	mems = m
	return &MemsDataController{}
}

// GetMemsData retrieves the mems data
func (mdc MemsDataController) GetMemsData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var d1, d2 []byte

	// Grab id
	id := p.ByName("id")

	if id == "" {
		d1, d2 = rosco.MemsReadRaw(mems)
	} else {
		cmd, _ := hex.DecodeString(id)
		rosco.MemsWriteSerial(mems, cmd)
		d1 = rosco.MemsReadSerial(mems)
	}

	// Stub an example
	md := rosco.MemsData{
		Id:          id,
		EngineRPM:   0,
		DataFrame80: hex.EncodeToString(d1),
		DataFrame7d: hex.EncodeToString(d2),
	}

	// Marshal provided interface into JSON structure
	mdj, _ := json.Marshal(md)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", mdj)
}
