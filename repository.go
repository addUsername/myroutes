package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

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

						if u.ActivityId == -1 {
							fmt.Println(" ADDing comment")
							routes[i].Clients[j].Comments = u.Comments
							break next

						} else {
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

func selectDays(routes string, startdate string, enddate string) string {

	fmt.Println("SELECT DAAYS FROM REPO")
	fmt.Println(startdate)
	fmt.Println(enddate)

	start := false
	end := false

	var startToCompare time.Time
	var endToCompare time.Time

	if startdate != "" {
		startToCompare = parseDate(startdate)
		start = true
	}
	if enddate != "" {
		endToCompare = parseDate(enddate)
		end = true
	}

	rout := parseJson(routes)

	var toReturn []Route

	for _, route := range rout {
		var ro []Client

		for _, client := range route.Clients {
			var ac []Activity

			for _, activity := range client.Activities {

				d := parseDate(activity.EndDate)

				if start && end {
					if d.After(startToCompare) && d.Before(endToCompare) {
						ac = append(ac, activity)
					}
				} else if start {
					if d.After(startToCompare) {
						ac = append(ac, activity)
					}
				} else {
					if d.Before(endToCompare) {
						ac = append(ac, activity)
					}
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

//TODO Update path from hander to point ./files/.
func JsonToExcel(pathNew string, pathOld string, routes string) bool {
	fmt.Println("json to excel")
	fmt.Println(pathOld)
	fmt.Println(pathNew)
	f, err := excelize.OpenFile(pathOld)
	if err != nil {
		fmt.Println(err)
		return false
	}
	rou := parseJson(routes)

	for i := 0; i < len(rou); i++ {
		clientIndex := 0
		activityIndex := 0
		client := Client{}
		rows, err := f.GetRows(rou[i].Name)
		if err != nil {
			fmt.Println(err)
			return false
		}

	next:
		for j := 1; j < len(rows); j++ {

			if len(rows[j]) == 1 {
				if len(rou[i].Clients) <= clientIndex {
					break next
				}
				fmt.Println("creating CLIENT")
				//search client name and update its dates
				client = rou[i].Clients[clientIndex]
				clientIndex = clientIndex + 1
				activityIndex = 0
			} else if len(client.Activities) > activityIndex && client.Activities[activityIndex].MadeDate != "" {
				/*if (client.Activities[activityIndex].MadeDate != ""){

				}*/
				f.SetCellValue(rou[i].Name, "D"+strconv.Itoa(j+1), client.Activities[activityIndex].MadeDate)
				style, err := f.NewStyle(`{"font":{"family":"Verdana","size":8,"color":"#000000"}}`)
				if err != nil {
					fmt.Println(err)
				}
				err = f.SetCellStyle(rou[i].Name, "D"+strconv.Itoa(j+1), "D"+strconv.Itoa(j+1), style)

				activityIndex = activityIndex + 1
			} else {
				activityIndex = activityIndex + 1
			}

		}
	}
	fmt.Println("saviing")
	// Save spreadsheet by the given path.
	if err := f.SaveAs(pathNew); err != nil {
		fmt.Println(err)
		return false
	}
	return true

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
