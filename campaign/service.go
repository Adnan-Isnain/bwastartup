package campaign

import (
	"bwastartup/user"
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(UserId int) ([]Campaign, error)
	GetcampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	// Get parameter from URL
	var campaigns []Campaign
	var err error

	if userID == 0 {
		campaigns, err = s.repository.FindAll()
	} else {
		campaigns, err = s.repository.FindByUserID(userID)
	}

	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (s *service) GetcampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByCampaignID(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	slugName := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign := Campaign{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		UserID:           input.User.ID,
		Slug:             slug.Make(slugName),
	}

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	//  Find data on DB
	myCampaign, err := s.repository.FindByCampaignID(inputID.ID)

	if err != nil {
		return myCampaign, err
	}

	// Checking owner of the campaign
	if myCampaign.UserID != inputData.User.ID {
		return myCampaign, errors.New("not an owner of the campaign")
	}

	// After get data, edit, then updated to repository
	myCampaign.Name = inputData.Name
	myCampaign.ShortDescription = inputData.ShortDescription
	myCampaign.Description = inputData.Description
	myCampaign.Perks = inputData.Perks
	myCampaign.GoalAmount = inputData.GoalAmount

	// Setting user to default to prevent error from postgres
	myCampaign.User = user.User{}
	updCampaign, err := s.repository.UpdateCampaign(myCampaign)
	if err != nil {
		return updCampaign, err
	}

	return updCampaign, nil
}
