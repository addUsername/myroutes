package main

import (
	"fmt"
	"os"

	"github.com/wailsapp/wails"
)

type Handler struct {
	r     *wails.Runtime
	store *wails.Store
}

/*
func (c *Handler) SendDialog(text string) {

	fmt.Println("SEND DIALOOG")
	c.r.Events.Emit("debug", text)
}
*/

func (c *Handler) WailsShutdown() {
	// De-Allocate some resources...
	fmt.Println("WAILS SHUTDOWN")
}

func (c *Handler) WailsInit(runtime *wails.Runtime) error {

	fmt.Println("INIT!!! ")

	// get path()
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	// check if exist and create or read config.json file

	// check if routes json exists and load that first
	routes := readRoutesJson(path + "\\routes.json")
	//if not, check if .xlsx generated by this program exist and load that
	if routes == "" {
		routes = loadFile(path + "\\PROGRAMACION-UPDATED.xlsx")
	}
	//if not, check for external .xlsx
	if routes == "" {
		routes = loadFile(path + "\\PROGRAMACION.xlsx")
	}
	fmt.Println("routes loaded")

	c.r = runtime
	c.store = runtime.Store.New("Clients", routes)
	//c.search = runtime.Store.New("Search", "")
	return nil
}

// this should dissapear, i should use the store from js instead
func (c *Handler) GetAllClients() string {

	//c.SendDialog("hi from goo")
	fmt.Println("HANDLER RESET CLIENTS!!! ")
	//c.r.Events.Emit("debug", "hello from gooo")
	return c.store.Get().(string)
}

// DoSearch deberia buscar fechas de endate comprendidas entre startdate y finishdate
func (c *Handler) DoSearch(route string, startdate string, finishdate string) string {
	fmt.Println("dousearch() ROUTE HANDLEER")

	clients := c.store.Get()
	search := filterRoute(clients.(string), route)
	if len(startdate) > 1 {
		search = selectDays(search, startdate)
	}
	//c.search.Set(final_search)
	return search
}

/*
func (c *Handler) GetByRoute(routeName string) string {
	fmt.Println("GET BY ROUTE HANDLEER")
	clients := c.store.Get()

	search := filterRoute(clients.(string), routeName)
	c.search.Set(search)
	return search
}
*/
func (c *Handler) Save(updates string) string {
	fmt.Println("SAVE")

	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	//update clients
	fmt.Println("updating routes")
	json := updateRoutes(c.store.Get().(string), updates)

	bol := saveJson(path+"\\routes.json", json)
	c.store.Set(json)
	//c.store.Update("")
	//c.store.Set("")
	//c.r.Events.Emit("debug", "SAVEED MTFACKA")
	if bol {
		fmt.Println("yeaah, sending json")
		return json
	}
	fmt.Println("oh, noes :(")
	return ""
}
func (c *Handler) GetByName(word string) {
	fmt.Println("GET BY NAME HANDLEER")
}
