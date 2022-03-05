package main

import (
	"flag"
	"log"
	"sign/puzzle/server"
)

var (
	mode    string
	showPic string
	port    string
)

func init() {
	flag.StringVar(&mode, "mode", "server", "running mode, server or local")
	flag.StringVar(&showPic, "show", "true", "show pic in log when use server")
	flag.StringVar(&port, "port", "8888", "port when use server")
}

func main() {
	flag.Parse()
	if mode != "server" {
		log.Println("running in local")
		server.RunLocal()
		return
	}

	log.Println("running as a server ...")
	server.Init(showPic, port)
	server.Run()
}
