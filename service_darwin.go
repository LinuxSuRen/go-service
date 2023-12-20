//go:build darwin
// +build darwin

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
)

type macOSService struct {
	CommonService
	cli string
}

func NewService(svc CommonService) (service Service, err error) {
	svc.scriptPath = EmptyThenDefault(svc.scriptPath, "/Library/LaunchDaemons/%s.plist", svc.ID)
	svc.script, err = render(macOSServiceScript, svc)
	service = &macOSService{
		CommonService: svc,
	}
	return
}

func (s *macOSService) Start() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("sudo", "", s.cli, "start", s.ID)
	return
}

func (s *macOSService) Stop() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("sudo", "", s.cli, "stop", s.ID)
	return
}

func (s *macOSService) Restart() (output string, err error) {
	if output, err = s.Stop(); err == nil {
		output, err = s.Start()
	}
	return
}

func (s *macOSService) Status() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("sudo", "", s.cli, "runstats", s.ID)
	return
}

func (s *macOSService) Install() (output string, err error) {
	if err = os.WriteFile(s.scriptPath, []byte(s.script), os.ModeAppend); err == nil {
		output, err = s.Execer.RunCommandAndReturn("sudo", "", s.cli, "enable", s.ID)
	}
	return
}

func (s *macOSService) Uninstall() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("sudo", "", s.cli, "disable", s.ID)
	return
}

func (s *macOSService) Available() bool {
	// TODO need a way to determine if it's available
	return true
}
