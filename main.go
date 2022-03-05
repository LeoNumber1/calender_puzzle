package main

import (
	"flag"
	"log"
	"sign/puzzle/server"
)

var (
	mode    string
	showPic string
)

func init() {
	flag.StringVar(&mode, "mode", "server", "running mode, server or local")
	flag.StringVar(&showPic, "show", "true", "show pic in log when use server")
}

func main() {
	flag.Parse()
	if mode != "server" {
		log.Println("running in local")
		server.RunLocal()
		return
	}

	log.Println("running as a server ...")
	server.Init(showPic)
	server.Run()
}
