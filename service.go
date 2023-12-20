/*
Copyright 2023 Rick.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package serivce

import (
	"bytes"
	_ "embed"
	"errors"
	"text/template"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
)

const (
	SystemCtl = "systemctl"
)

// Service is the interface of service
type Service interface {
	Start() (string, error)     // start the service
	Stop() (string, error)      // stop the service gracefully
	Restart() (string, error)   // restart the service gracefully
	Status() (string, error)    // status of the service
	Install() (string, error)   // install the service
	Uninstall() (string, error) // uninstall the service
	Available() bool
}

type CommonService struct {
	ID          string
	Name        string
	Description string
	Command     string
	Args        []string

	Execer     fakeruntime.Execer
	scriptPath string
	script     string
}

type ServiceMode string

const (
	ServiceModeOS        ServiceMode = "os"
	ServiceModeContainer ServiceMode = "container"
	ServiceModePodman    ServiceMode = "podman"
	ServiceModeDocker    ServiceMode = "docker"
)

func (s ServiceMode) All() []ServiceMode {
	return []ServiceMode{ServiceModeOS, ServiceModeContainer,
		ServiceModePodman, ServiceModeDocker}
}

func (s ServiceMode) String() string {
	return string(s)
}

type ServerFeatureOption struct {
	ID      string
	Name    string
	Command string
	Args    []string
}

func GetAvailableService(mode ServiceMode,
	containerOption ContainerOption, service CommonService) (svc Service, err error) {
	svc, err = NewService(service)
	dockerService := NewContainerService(service, "docker", containerOption)
	podmanService := NewContainerService(service, "podman", containerOption)

	if mode == "" && svc.Available() {
		mode = ServiceModeOS
	}

	switch mode {
	case ServiceModeOS:
		// using the default value
	case ServiceModeDocker, ServiceModeContainer:
		svc = dockerService
	case ServiceModePodman:
		svc = podmanService
	case "":
		if dockerService.Available() {
			svc = dockerService
		} else if podmanService.Available() {
			svc = podmanService
		}
	}

	if svc != nil && !svc.Available() {
		err = errors.New("not found available service")
	}
	return
}

func render(tplText string, svc CommonService) (result string, err error) {
	var tpl *template.Template
	if tpl, err = template.New(svc.Name).Parse(tplText); err != nil {
		return
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, svc); err == nil {
		result = buf.String()
	}
	return
}

//go:embed data/macos_service.xml
var macOSServiceScript string

//go:embed data/linux_service.txt
var linuxServiceScript string
