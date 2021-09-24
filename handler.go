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

func (s *Handler) WailsShutdown() {
	// De-Allocate some resources...
}

func (c *Handler) WailsInit(runtime *wails.Runtime) error {

	fmt.Println("HANDLER!!! ")

	// refactor
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	clients := loadFile(path + "\\PROGRAMACION.xlsx")

	c.r = runtime
	c.store = runtime.Store.New("Clients", clients)
	c.search = runtime.Store.New("Search", "")
	return nil
}
func (c *Handler) GetAllClients() string {

	fmt.Println("HANDLER RESET CLIENTS!!! ")
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
func (c *Handler) GetByDate(EndDate string, MadeDate string) {
	fmt.Println("GET BY DATE HANDLEER")
	//c.store.Update("")
	//c.store.Set("")
}
func (c *Handler) GetByName(word string) {
	fmt.Println("GET BY NAME HANDLEER")
}
