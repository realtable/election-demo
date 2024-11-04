package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/realtable/template/model"
)

func GetVotes(w http.ResponseWriter, r *http.Request) {
	votes, err := model.GetVotes()
	if err != nil {
		slog.Error("could not get votes", "error", err)
		http.Error(w, "could not get votes", http.StatusInternalServerError)
		return
	}
	for _, vote := range votes {
		w.Write([]byte("<tr><td>" + vote.VoteKey + "</td><td>" + vote.Candidate + "</td></tr>"))
	}
}

func AddVote(w http.ResponseWriter, r *http.Request) {
	err := model.AddVote(r.FormValue("voteKey"), r.FormValue("candidate"))
	if errors.Is(err, model.ErrVoterNotFound) || errors.Is(err, model.ErrAlreadyVoted) {
		http.Error(w, "<div class=\"alert alert-danger\" role=\"alert\">Vote failed to submit: "+err.Error()+"</div>", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "<div class=\"alert alert-danger\" role=\"alert\">Vote failed to submit. Please try again later.</div>", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("<div class=\"alert alert-success\" role=\"alert\">Vote successfully submitted!</div>"))
}

func ClearVotes(w http.ResponseWriter, r *http.Request) {
	model.ClearVotes()
}
