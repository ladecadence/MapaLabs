package models

type Typology struct {
	Location bool `json:"location"`
	Nomad    bool `json:"nomad"`
}

type Governance struct {
	Public  bool `json:"public"`
	Private bool `json:"private"`
}

type Theme struct {
	Environment    bool `json:"environment"`
	DigitalCulture bool `json:"digital_culture"`
	Arts           bool `json:"arts"`
	Territory      bool `json:"territory"`
	CitizenScience bool `json:"citizen_science"`
	Memory         bool `json:"memory"`
	Gender         bool `json:"gender"`
}

type Lab struct {
	Id                  int        `json:"id"`
	Name                string     `json:"name"`
	City                string     `json:"city"`
	Country             string     `json:"country"`
	Description         string     `json:"description"`
	Date                string     `json:"date"`
	Works               string     `json:"works"`
	Motivations         string     `json:"motivations"`
	Networks            string     `json:"networks"`
	Typology            Typology   `json:"typology" gorm:"embedded"`
	Governance          Governance `json:"governance" gorm:"embedded"`
	Themes              Theme      `json:"themes" gorm:"embedded"`
	Web                 string     `json:"web"`
	Mastodon            string     `json:"mastodon"`
	Instagram           string     `json:"instagram"`
	Facebook            string     `json:"facebook"`
	Twitter             string     `json:"twitter"`
	Spotify             string     `json:"spotify"`
	Linkedin            string     `json:"linkedin"`
	TikTok              string     `json:"tiktok"`
	Flickr              string     `json:"flickr"`
	Twitch              string     `json:"twitch"`
	Youtube             string     `json:"youtube"`
	Delegate            string     `json:"delegate"`
	DelegatePosition    string     `json:"delegate_position"`
	DelegateDescription string     `json:"delegate_description"`
	Image               string     `json:"image"`
	Latitude            float64    `json:"latitude"`
	Longitude           float64    `json:"longitude"`
}
