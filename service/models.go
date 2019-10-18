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
	var memsdata rosco.MemsData

	// Grab id
	id := p.ByName("id")

	if id == "" {
		memsdata = rosco.MemsRead(mems)

		// Marshal provided interface into JSON structure
		mdj, _ := json.Marshal(memsdata)

		// Write content-type, statuscode, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", mdj)
	} else {
		cmd, _ := hex.DecodeString(id)
		rosco.MemsWriteSerial(mems, cmd)
		response := rosco.MemsReadSerial(mems)

		// Write content-type, statuscode, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%x", response)
	}
}
