package model

type Power struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		DeveloperPower   string `json:"developerPower"`
		SpPower          string `json:"spPower"`
		ClientPower      string `json:"clientPower"`
		TokenHolderPower string `json:"tokenHolderPower"`
	} `json:"data"`
}
