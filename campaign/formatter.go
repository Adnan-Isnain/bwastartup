package campaign

import (
	"strings"
)

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserId           int    `json:"user_id"`
	Name             string `json:"title"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	Slug             string `json:"slug"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{
		ID:               campaign.ID,
		UserId:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		ImageURL:         "",
		Slug:             campaign.Slug,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
	}

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}
	return campaignFormatter
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	allCampaigns := []CampaignFormatter{}
	for _, campaign := range campaigns {
		allCampaigns = append(allCampaigns, FormatCampaign(campaign))
	}

	return allCampaigns
}

type CampaignDetailFormatter struct {
	ID               int                      `json:"id"`
	Name             string                   `json:"name"`
	ShortDescription string                   `json:"short_description"`
	Description      string                   `json:"description"`
	ImageURL         string                   `json:"image_url"`
	GoalAmount       int                      `json:"goal_amount"`
	CurrentAmount    int                      `json:"current_amount"`
	UserID           int                      `json:"user_id"`
	Slug             string                   `json:"slug"`
	Perks            []string                 `json:"perks"`
	User             CampaignUserFormatter    `json:"user"`
	Images           []CampaignImageFormatter `json:"images"`
}

type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignImageFormatter struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	campaignDetailFormatter := CampaignDetailFormatter{
		ID:               campaign.ID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		ImageURL:         "",
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		UserID:           campaign.UserID,
		Slug:             campaign.Slug,
		Perks:            strings.Split(campaign.Perks, ", "),
		User: CampaignUserFormatter{
			Name:     campaign.User.Name,
			ImageURL: campaign.User.AvatarFileName,
		},
		Images: []CampaignImageFormatter{},
	}

	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	for _, campaignimg := range campaign.CampaignImages {
		isPrimary := false
		if campaignimg.IsPrimary == 1 {
			isPrimary = true
		}

		campaignImage := CampaignImageFormatter{
			ImageURL:  campaignimg.FileName,
			IsPrimary: isPrimary,
		}
		campaignDetailFormatter.Images = append(campaignDetailFormatter.Images, campaignImage)
	}

	return campaignDetailFormatter
}
