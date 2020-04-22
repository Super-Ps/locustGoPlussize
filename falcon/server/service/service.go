package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type ServerEntry struct {
	Engine 			*gin.Engine
	Ip 				string
	Port 			int
	RouteData		[]*RouteData
}

type RouteData struct {
	Name 			string
	Method 			string
	Route 			string
	Func 			gin.HandlerFunc
}

type JsonRoute struct {
	Route			 Route	`json:"route"`
}

type Route struct {
	Start 			string	`json:"start"`
	Stop 			string	`json:"stop"`
	Restart 		string	`json:"restart"`
	Exit 			string	`json:"exit"`
}

var indexName = "index.html"
var indexPath = "./" + indexName
var staticName = "static"

func init() {
	LoadArgv()
	_ = RestoreAssets("./", indexName)
	_ = RestoreAssets(fmt.Sprintf("./%s", strings.Trim(Param.Route.Root, "/")), staticName)
}

func New() *ServerEntry {
	var routeData []*RouteData
	var engine *gin.Engine
	gin.SetMode(Param.Mode)
	if Param.NoOut {
		engine = gin.New()
	} else {
		engine = gin.Default()
	}
	return &ServerEntry{
		Engine: engine,
		Ip: Param.Ip,
		Port: Param.Port,
		RouteData: routeData,
	}
}

func (s *ServerEntry) Start() {
	go UpdateStart()
	go s.gc_monitor()
	s.loadResource()
	s.addDefaultFunc()
	for _, v := range s.RouteData {
		s.Engine.Handle(v.Method, v.Route, v.Func)
		Printf("%s(%s) route is %s\n", v.Name, v.Method, v.Route)
	}
	time.Sleep(time.Duration(3)*time.Second)
	Printf("Server name is %s\n", Param.Name)
	Printf("Server mode is %s\n", Param.Mode)
	Printf("GC Time is %ds\n", Param.GcTime)
	Printf("Pause Time is %ds\n", Param.PauseTime)
	Printf("Interval Time is %ds\n", Param.IntervalTime)
	Printf("Pid path is %s\n", Param.PidPath)
	Printf("Http server bind ip is %s\n", s.Ip)
	Printf("Http server listen port at %d\n", s.Port)
	go s.exit()
	if err := s.Engine.Run(fmt.Sprintf("%s:%d", s.Ip, s.Port)); err != nil {
		Printf("Http server listen fail at %s\n", err.Error())
		os.Exit(1)
	}
	Printf("Http server stop\n")
	os.Exit(0)
}

func (s *ServerEntry) AddRoute(name string, method string, route string, funcHandler gin.HandlerFunc) {
	s.RouteData = append(s.RouteData, &RouteData{
		Name: name,
		Method: method,
		Route: route,
		Func: funcHandler,
	})
}

func (s *ServerEntry) addDefaultFunc() {
	s.AddRoute("Root", "GET", Param.Route.Root, s.WebPage)
	s.AddRoute("Monitor-Map", "GET", fmt.Sprintf("%s/%s", Param.Route.Monitor, "map"), s.HttpMonitorMap)
	s.AddRoute("Monitor-List", "GET", fmt.Sprintf("%s/%s", Param.Route.Monitor, "list"), s.HttpMonitorList)
	s.AddRoute("Stop", "POST", Param.Route.Stop, s.HttpStop)
	s.AddRoute("Restart", "POST", Param.Route.Restart, s.HttpRestart)
	s.AddRoute("Quit", "POST", Param.Route.Quit, s.HttpQuit)
	s.AddRoute("Address", "GET", Param.Route.Address, s.HttpAddress)
	s.AddRoute("GetDemo", "GET", Param.Route.GetDemo, s.HttpGetDemo)
	s.AddRoute("PostDemo", "POST", Param.Route.PostDemo, s.HttpPostDemo)

}

func (s *ServerEntry) loadResource() {
	s.Engine.LoadHTMLFiles(indexPath)
	s.Engine.Static(fmt.Sprintf("./%s", strings.Trim(Param.Route.Static, "/")), strings.Trim(Param.Route.Static, "/"))
}

func (s *ServerEntry) gc_monitor(){
	for {
		time.Sleep(time.Duration(Param.GcTime)*time.Second)
		runtime.GC()
	}
}

func (s *ServerEntry) exit() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	go s.quit()
}

func (s *ServerEntry) quit() {
	go s.exit()
	time.Sleep(time.Duration(1)*time.Second)
	Printf("Http server stop\n")
	os.Exit(0)
}

func deletePidFile() {
	_, err := os.Stat(Param.PidPath)
	if err == nil {
		_ = os.Remove(Param.PidPath)
	}
}

