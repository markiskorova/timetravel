package apiv1

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
