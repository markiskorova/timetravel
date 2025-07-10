package apiv2

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/util"
)

// GET /records/{id}
func (a *API) GetRecords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		util.WriteError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		return
	}

	record, err := a.records.GetRecord(ctx, int(idNumber))
	if err != nil {
		err := util.WriteError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		util.LogError(err)
		return
	}

	err = util.WriteJSON(w, record, http.StatusOK)
	util.LogError(err)
}

// GET /api/v2/records/{id}/versions
func (a *API) ListRecordVersions(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		return
	}

	versions, err := a.records.ListRecordVersions(r.Context(), id)
	if err != nil {
		util.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.WriteJSON(w, versions, http.StatusOK)
}

// GET /api/v2/records/{id}/versions/{version}
func (a *API) GetRecordVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err1 := strconv.Atoi(vars["id"])
	version, err2 := strconv.Atoi(vars["version"])

	if err1 != nil || err2 != nil {
		util.WriteError(w, "invalid id or version", http.StatusBadRequest)
		return
	}

	record, err := a.records.GetRecordVersion(r.Context(), id, version)
	if err != nil {
		util.WriteError(w, err.Error(), http.StatusNotFound)
		return
	}

	util.WriteJSON(w, record, http.StatusOK)
}
