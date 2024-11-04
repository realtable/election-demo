package handler

import (
	"log/slog"
	"net/http"

	"github.com/realtable/template/model"
)

func GetVoters(w http.ResponseWriter, r *http.Request) {
	voters, err := model.GetVoters()
	if err != nil {
		slog.Error("could not get voters", "error", err)
		http.Error(w, "could not get voters", http.StatusInternalServerError)
		return
	}
	for _, voter := range voters {
		w.Write([]byte("<tr><td>" + voter.Name + "</td><td>" + voter.Address + "</td><td>" + voter.VoteKey + "</td></tr>"))
	}
}

func AddVoter(w http.ResponseWriter, r *http.Request) {
	err := model.AddVoter(r.PostFormValue("name"), r.PostFormValue("address"))
	if err != nil {
		slog.Error("could not add voter", "error", err)
		http.Error(w, "could not add voter", http.StatusInternalServerError)
	}
}

func ClearVoters(w http.ResponseWriter, r *http.Request) {
	model.ClearVotes()
}
