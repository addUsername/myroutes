package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type Route struct {
	Name    string   `json:"name"`
	Clients []Client `json:"clients"`
}

// warnings
type Client struct {
	Name       string     `json:"name"`
	Dir        string     `json:"dir"`
	Route      string     `json:"routename"`
	YWarn      int        `json:"y_warn"`
	RWarn      int        `json:"r_warn"`
	Activities []Activity `json:"activities"`
}

// TODO:
type Activity struct {
	Name     string `json:"name"`
	EndDate  string `json:"enddate"`
	MadeDate string `json:"madedate"`
}

func parseDate(date string) time.Time {
	fmt.Println("PARSE DATE FROM REPO")
	fmt.Println(date)
	d := strings.Split(date, "/")
	if len(d) < 2 {
		d = strings.Split(date, "-")
	}
	if len(d[2]) == 2 {
		d[2] = "20" + d[2]
	}
	year, _ := strconv.Atoi(d[2])
	month, _ := strconv.Atoi(d[0])
	day, _ := strconv.Atoi(d[1])

	fmt.Println("-------------------")
	fmt.Println(day)
	fmt.Println(time.Month(month))
	fmt.Println(year)
	fmt.Println("-------------------")

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func selectDays(routes string, enddate string) string {

	fmt.Println("SELECT DAAYS FROM REPO")
	dateToCompare := parseDate(enddate)

	rout := parseJson(routes)

	var toReturn []Route

	for _, route := range rout {
		var ro []Client

		for _, client := range route.Clients {
			var ac []Activity

			for _, activity := range client.Activities {

				d := parseDate(activity.EndDate)

				if d.Before(dateToCompare) {
					ac = append(ac, activity)
				}
			}
			if len(ac) > 0 {
				client.Activities = ac
				ro = append(ro, client)
			}
		}
		if len(ro) > 0 {
			route.Clients = ro
			toReturn = append(toReturn, route)
		}

	}
	e, err := json.Marshal(toReturn)
	//e, err := json.MarshalIndent(clients, "", "  ")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(e)

}

func parseJson(clientsJSON string) []Route {

	fmt.Println("PARSE JSON()")
	var rou []Route
	err := json.Unmarshal([]byte(clientsJSON), &rou)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return rou
}
func filterRoute(clientsJSON string, routename string) string {

	fmt.Println("FILTER ROUTES")
	rou := parseJson(clientsJSON)

	var toReturn []Route
	for _, route := range rou {
		if route.Name == routename {
			toReturn = append(toReturn, route)
			break
		}
	}
	e, err := json.Marshal(toReturn)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(e)
}
func loadFile(path string) string {

	fmt.Println("LOAD FILEE")

	//Open spreadsheet
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var all []Route
	var rou Route
	var client Client
	end := ""

	// Get all sheets names
	routes := f.GetSheetList()

	for _, route := range routes {
		rou = Route{}
		rou.Name = route

		var clients []Client
		// Get all the rows in the Sheet.
		rows, err := f.GetRows(route)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		// If there is blank sheets
		if len(rows) == 0 {
			break
		}
		for _, row := range rows[1:] {

			if len(row) == 1 {

				if len(client.Name) != 0 {
					clients = append(clients, client)
				}

				client = Client{}
				client.Route = route
				// TODO check regex posibilities
				split := strings.Split(row[0], "TOLEDO")
				client.Name = split[0]
				client.Dir = split[1]

			} else {
				// madedate
				end = ""
				if len(row) > 3 {
					end = row[3]
				}
				activity := Activity{row[1], row[2], end}
				client.Activities = append(client.Activities, activity)
			}
		}
		rou.Clients = clients
		all = append(all, rou)

	}

	e, err := json.Marshal(all)
	//e, err := json.MarshalIndent(clients, "", "  ")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(e)
}
