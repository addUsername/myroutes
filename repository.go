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
	Id      int      `json:"routeid"`
	Name    string   `json:"name"`
	Clients []Client `json:"clients"`
}

type Client struct {
	Id         int        `json:"clientid"`
	Name       string     `json:"name"`
	Dir        string     `json:"dir"`
	Route      string     `json:"routename"`
	Activities []Activity `json:"activities"`
}

// TODO:
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
}

var DAYS int = 30

func updateRoutes(oldjson string, update string) string {
	fmt.Println("update routes")
	fmt.Println(update)

	up := parseUpdateJson(update)
	routes := parseJson(oldjson)

	fmt.Println("from parsing jsons")
	// lel, we could improve this by correlating each id with its index (route id is ok)
	for _, u := range up {
	next:
		//for _, route := range routes {
		for i := 0; i < len(routes); i++ {

			if routes[i].Id == u.RouteId {

				for j := 0; j < len(routes[i].Clients); j++ {

					if routes[i].Clients[j].Id == u.ClientId {

						for k := 0; k < len(routes[i].Clients[j].Activities); k++ {
							activity := routes[i].Clients[j].Activities[k]

							if activity.Id == u.ActivityId {
								fmt.Println("hello, creating activity")

								a := CreateActivity([]string{activity.Name, activity.EndDate, u.Madedate})
								a.Id = activity.Id
								routes[i].Clients[j].Activities[k] = a
								// cool
								break next
							}
						}
					}

				}
			}
		}
	}

	e, err := json.Marshal(routes)
	//e, err := json.MarshalIndent(clients, "", "  ")
	if err != nil {
		fmt.Println("--error on updateRoutes--")
		fmt.Println(err)
		return ""
	}

	return string(e)

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

// TODO.. filter by routeid !!!, modify selector value on main.js too
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
	fmt.Println(path)
	//Open spreadsheet
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var all []Route
	var rou Route
	var client Client
	//end := ""

	// Get all sheets names
	routes := f.GetSheetList()

	fmt.Print("LOAD FILE Routes:")
	fmt.Println(len(routes))
	for index, route := range routes {
		rou = Route{}
		rou.Name = route
		rou.Id = index

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
		for idx, row := range rows[1:] {

			if len(row) == 1 {
				// if row has 1 column then create new client
				if len(client.Name) != 0 {
					clients = append(clients, client)
				}

				client = Client{}
				client.Route = route
				client.Id = idx
				// TODO check regex posibilities
				split := strings.Split(row[0], "TOLEDO")
				client.Name = split[0]
				client.Dir = split[1]

			} else {
				// if not, that row is an activity
				// TODO: no hardcode the days
				activity := CreateActivity(row[1:])
				activity.Id = idx
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
func CreateActivity(row []string) Activity {

	toReturn := Activity{}

	toReturn.Name = row[0]
	toReturn.EndDate = row[1]

	// madedate
	end := ""
	if len(row) == 3 {
		fmt.Print("len roow 3")
		fmt.Println(row[2])
		end = row[2]
	}
	toReturn.MadeDate = end

	// check warnings
	// r -not completed on time / y - before timeline / g - done / b - pendint but good
	// dates to compare, var now should come as argument
	now := time.Now()
	timeline := now.AddDate(0, 0, +DAYS)
	endate := parseDate(toReturn.EndDate)

	if endate.Before(now) {
		if toReturn.MadeDate == "" {
			toReturn.ColorWarn = 'r'
		} else {
			toReturn.ColorWarn = 'g'
		}
	} else if endate.Before(timeline) {
		toReturn.ColorWarn = 'y'
	} else {
		toReturn.ColorWarn = 'b'
	}
	return toReturn
}

func parseDate(date string) time.Time {
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

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
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

func parseUpdateJson(up string) []Update {

	fmt.Println("parse update json()")
	var u []Update
	err := json.Unmarshal([]byte(up), &u)
	if err != nil {
		fmt.Println("--Error--")
		fmt.Println(err)
		return nil
	}
	return u

}
