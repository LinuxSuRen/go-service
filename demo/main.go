package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	svc "github.com/linuxsuren/go-service"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-svc [service|server]")
		os.Exit(1)
	}

	var err error
	switch cmd := os.Args[1]; cmd {
	case "server":
		server()
	case "service":
		if len(os.Args) < 3 {
			err = errors.New("Usage: go-svc service [start|stop|install|uninstall]")
		} else {
			err = service(os.Args[2])
		}
	default:
		fmt.Println("unknown command:", cmd)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func server() {
	for {
		fmt.Println(".")
		time.Sleep(3 * time.Second)
	}
}

func service(action string) (err error) {
	var ss svc.Service
	ss, err = svc.GetAvailableService(svc.ServiceModeOS,
		svc.ContainerOption{},
		svc.CommonService{
			ID:      "go-svc",
			Name:    "go-svc",
			Command: "go-svc",
			Args:    []string{"server"},
			Execer:  fakeruntime.NewDefaultExecer(),
		})
	if err != nil {
		return
	}

	var output string
	switch action {
	case "start":
		output, err = ss.Start()
	case "install":
		output, err = ss.Install()
	case "stop":
		output, err = ss.Stop()
	case "uninstall":
		output, err = ss.Uninstall()
	}

	if output != "" {
		fmt.Println(output)
	}
	return
}
