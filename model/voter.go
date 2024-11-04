package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Voter struct {
	gorm.Model
	VoteKey string `gorm:"uniqueIndex"`
	Name    string
	Address string
}

func GetVoters() ([]Voter, error) {
	var voters []Voter
	result := db.Find(&voters)
	return voters, result.Error
}

func AddVoter(name string, address string) error {
	voteKey, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	result := db.Create(&Voter{VoteKey: voteKey.String(), Name: name, Address: address})
	return result.Error
}

func ClearVoters() {
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Voter{})
}
