package main

type Route struct {
	Id      int      `json:"routeid"`
	Name    string   `json:"name"`
	Clients []Client `json:"clients"`
}

type Client struct {
	Id         int        `json:"clientid"`
	Name       string     `json:"name"`
	Dir        string     `json:"dir"`
	Route      string     `json:"routename"`
	Comments   string     `json:"comments"`
	Activities []Activity `json:"activities"`
}

type Activity struct {
	Id        int    `json:"activityid"`
	Name      string `json:"name"`
	EndDate   string `json:"enddate"`
	MadeDate  string `json:"madedate"`
	ColorWarn byte   `json:"colorwarning"`
}

type Update struct {
	RouteId    int    `json:"routeid"`
	ClientId   int    `json:"clientid"`
	ActivityId int    `json:"activityid"`
	Madedate   string `json:"madedate"`
	Comments   string `json:"comments"`
}
