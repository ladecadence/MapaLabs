package models

type Lab struct {
	Id                  int     `json:"id"`
	Name                string  `json:"name"`
	City                string  `json:"city"`
	Country             string  `json:"country"`
	Description         string  `json:"description"`
	Date                string  `json:"date"`
	Works               string  `json:"works"`
	Motivations         string  `json:"motivations"`
	Networks            string  `json:"networks"`
	Web                 string  `json:"web"`
	Mastodon            string  `json:"mastodon"`
	Instagram           string  `json:"instagram"`
	Facebook            string  `json:"facebook"`
	Twitter             string  `json:"twitter"`
	Spotify             string  `json:"spotify"`
	Linkedin            string  `json:"linkedin"`
	TikTok              string  `json:"tiktok"`
	Flickr              string  `json:"flickr"`
	Twitch              string  `json:"twitch"`
	Youtube             string  `json:"youtube"`
	Delegate            string  `json:"delegate"`
	DelegatePosition    string  `json:"delegate_position"`
	DelegateDescription string  `json:"delegate_description"`
	Image               string  `json:"image"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
}
