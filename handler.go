package main

import (
	"fmt"
	"os"

	"github.com/wailsapp/wails"
)

type Handler struct {
	r      *wails.Runtime
	store  *wails.Store
	search *wails.Store
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

	fmt.Println("HANDLER!!! ")

	// refactor
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	clients := loadFile(path + "\\PROGRAMACION.xlsx")
	fmt.Println("clients loaded")

	c.r = runtime
	c.store = runtime.Store.New("Clients", clients)
	//c.search = runtime.Store.New("Search", "")
	return nil
}
func (c *Handler) GetAllClients() string {

	//c.SendDialog("hi from goo")
	fmt.Println("HANDLER RESET CLIENTS!!! ")
	//c.r.Events.Emit("debug", "hello from gooo")
	return c.store.Get().(string)
}

func (c *Handler) DoSearch(route string, enddate string, madedate string) string {
	fmt.Println("dousearch() ROUTE HANDLEER")
	fmt.Println(route + " " + enddate + "" + madedate)

	clients := c.store.Get()
	search := filterRoute(clients.(string), route)
	if len(enddate) > 1 {
		search = selectDays(search, enddate)
	}
	//c.search.Set(final_search)
	return search
}
func (c *Handler) GetByRoute(routeName string) string {
	fmt.Println("GET BY ROUTE HANDLEER")
	clients := c.store.Get()

	search := filterRoute(clients.(string), routeName)
	c.search.Set(search)
	return search
}
func (c *Handler) Save(updates string) string {
	fmt.Println("SAVE")

	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	//update clients
	json := updateRoutes(c.store.Get().(string), updates)

	bol := saveJson(path+"\\routes.json", json)
	c.store.Set(json)
	//c.store.Update("")
	//c.store.Set("")
	//c.r.Events.Emit("debug", "SAVEED MTFACKA")
	if bol {
		return "saaveed"
	}
	return "fuck"
}
func (c *Handler) GetByName(word string) {
	fmt.Println("GET BY NAME HANDLEER")
}
