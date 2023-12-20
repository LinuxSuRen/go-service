//go:build linux
// +build linux

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
	"os"
	"strings"
)

func NewService(svc CommonService) (service Service, err error) {
	svc.scriptPath = EmptyThenDefault(svc.scriptPath, "/lib/systemd/system/%s.service", svc.ID)
	svc.script, err = render(linuxServiceScript, svc)
	service = &linuxService{svc}
	return
}

type linuxService struct {
	CommonService
}

func (s *linuxService) Start() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "start", s.Name)
	return
}

func (s *linuxService) Stop() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "stop", s.Name)
	return
}

func (s *linuxService) Restart() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "restart", s.Name)
	return
}

func (s *linuxService) Status() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "status", s.Name)
	if err != nil && err.Error() == "exit status 3" {
		// this is normal case
		err = nil
	}
	return
}

func (s *linuxService) Install() (output string, err error) {
	if err = os.WriteFile(s.scriptPath, []byte(s.script), os.ModeAppend); err == nil {
		output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "enable", s.Name)
	}
	return
}

func (s *linuxService) Uninstall() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "disable", s.Name)
	return
}

func (s *linuxService) Available() bool {
	output, _ := s.Execer.RunCommandAndReturn(SystemCtl, "", "is-system-running")
	output = strings.TrimSpace(output)
	return output != "offline"
}
