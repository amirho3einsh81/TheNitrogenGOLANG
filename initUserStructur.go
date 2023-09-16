package main

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"userName"`
	Pass      string `json:"pass"`
	PhoneNum  string `json:"phoneNum"`
	Title     string `json:"title"`
	Biography string `json:"biography"`
}
type UserInformation struct {
	Username  string `json:"userName"`
	Pass      string `json:"pass"`
	Title     string `json:"title"`
	Family    string `json:"family"`
	Profile   string `json:"profile"`
	Biography string `json:"biography"`
	Gender    bool   `json:"gender"`
	PhoneNum  string `json:"phoneNum"`
	Byte      []byte `json:"bytes"`
}
type Register struct {
	Username  string `json:"userName"`
	Title     string `json:"title"`
	Profile   string `json:"profile"`
	Biography string `json:"biography"`
	Pass      string `json:"pass"`
	PhoneNum  string `json:"phoneNum"`
	UserIP    string `json:"userIP"`
	Brand     string `json:"brand"`
	Device    string `json:"device"`
	Model     string `json:"model"`
}
type Location struct {
	Username    string `json:"username"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Region      string `json:"region"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Timezone    string `json:"timezone"`
	Isp         string `json:"isp"`
	Org         string `json:"org"`
	As          string `json:"as"`
	Query       string `json:"query"`
}
type Login struct {
	Username  string `json:"userName"`
	Title     string `json:"title"`
	Profile   string `json:"profile"`
	Biography string `json:"biography"`
	Pass      string `json:"pass"`
	PhoneNum  string `json:"phoneNum"`
	Theme     bool   `json:"theme"`
}
type Challenge struct {
	Id             int          `json:"id"`
	ChallengeID    string       `json:"challengeID"`
	Title          string       `json:"title"`
	SubTitle       string       `json:"subTitle"`
	Description    string       `json:"description"`
	SubDescription string       `json:"subDescription"`
	Image          string       `json:"image"`
	DayCount       int          `json:"dayCount"`
	BgColor        string       `json:"bgColor"`
	TxtColor       string       `json:"txtColor"`
	Tags           []Tag        `json:"tags"`
	Conclusion     []Conclusion `json:"conclusions"`
	Influence      []Influence  `json:"influence"`
	Category       []Category   `json:"category"`
	Prerequisite   []Challenge  `json:"prerequisite"`
	Tool           Tool         `json:"tool"`
	Suggestion     []Challenge  `json:"suggestion"`
	Started        Started      `json:"started"`
	Price          Price        `json:"price"`
}
type Challenge2 struct {
	Id             int        `json:"id"`
	ChallengeID    string     `json:"challengeID"`
	Title          string     `json:"title"`
	SubTitle       string     `json:"subTitle"`
	Description    string     `json:"description"`
	SubDescription string     `json:"subDescription"`
	Image          string     `json:"image"`
	DayCount       int        `json:"dayCount"`
	BgColor        string     `json:"bgColor"`
	TxtColor       string     `json:"txtColor"`
	Category       []Category `json:"category"`
}
type ImageEditor struct {
	Image string `json:"image"`
	Name  string `json:"name"`
	//Username string `json:"username"`
}
type Home struct {
	Newchallenges   []Challenge `json:"newchallenges"`
	Success         []Challenge `json:"success"`
	Psychology      []Challenge `json:"psychology"`
	Selfdevelopment []Challenge `json:"selfdevelopment"`
	Suggestion5     []Challenge `json:"suggestion5"`
	Banners         []Challenge `json:"banners"`
}
type Started struct {
	Id          int        `json:"id"`
	Username    string     `json:"username"`
	ChallengeID string     `json:"challengeID"`
	Day         int        `json:"day"`
	Month       int        `json:"month"`
	Year        int        `json:"year"`
	DayOfYear   int        `json:"dayOfYear"`
	State       int        `json:"state"`
	DayCount    int        `json:"dayCount"`
	PriceType   string     `json:"priceType"`
	Progress    int        `json:"progress"`
	Challenge   Challenge2 `json:"challenge"`
}
type Tool struct {
	ChallengeID string  `json:"challengeID"`
	Title       string  `json:"title"`
	ToolID      string  `json:"toolID"`
	Description string  `json:"description"`
	Tools       []Tools `json:"tools"`
}
type Tools struct {
	ToolID string `json:"toolID"`
	Name   string `json:"name"`
}
type Tag struct {
	ChallengeID string `json:"challengeID"`
	Title       string `json:"title"`
}
type Price struct {
	PriceID     string `json:"priceID"`
	ChallengeID string `json:"challengeID"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Price       string `json:"price"`
	Sku         string `json:"sku"`
}
type Purchase struct {
	PriceID     string `json:"priceID"`
	ChallengeID string `json:"challengeID"`
	Date        int    `json:"date"`
	DayOfYear   int    `json:"dayOfYear"`
}
type Conclusion struct {
	//Id          int    `json:"id"`
	ChallengeID string `json:"challengeID"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
type Influence struct {
	//Id          int    `json:"id"`
	ChallengeID string `json:"challengeID"`
	Title       string `json:"title"`
}
type Category struct {
	//Id          int    `json:"id"`
	ChallengeID string `json:"challengeID"`
	Title       string `json:"title"`
}
type Prerequisite struct {
	ChallengeID    string `json:"challengeID"`
	PrerequisiteID string `json:"prerequisiteID"`
}
type Banner struct {
	Id          string `json:"id"`
	ChallengeID string `json:"challengeID"`
}
type goal struct {
	//Id     int    `json:"id"`
	GoalID      string `json:"goalID"`
	UserName    string `json:"userName"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Color       string `json:"color"`
	StartDay    string `json:"startDay"`
	DayCount    int    `json:"dayCount"`
	Life        int    `json:"life"`
	Update      string `json:"update"`
}
