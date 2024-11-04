package model

import (
	"errors"

	"gorm.io/gorm"
)

var ErrVoterNotFound = errors.New("voter not found")
var ErrAlreadyVoted = errors.New("already voted")

type Vote struct {
	gorm.Model
	VoteKey   string `gorm:"uniqueIndex"`
	Candidate string
}

func GetVotes() ([]Vote, error) {
	var votes []Vote
	result := db.Find(&votes)
	return votes, result.Error
}

func AddVote(voteKey string, candidate string) error {
	var eligibleVoters int64
	if db.Model(&Voter{}).Where("vote_key = ?", voteKey).Count(&eligibleVoters); eligibleVoters == 0 {
		return ErrVoterNotFound
	}

	var existingVotes int64
	if db.Model(&Vote{}).Where("vote_key = ?", voteKey).Count(&existingVotes); existingVotes > 0 {
		return ErrAlreadyVoted
	}

	db.Create(&Vote{VoteKey: voteKey, Candidate: candidate})
	return nil
}

func ClearVotes() {
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Vote{})
}
