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
	_ "embed"
	"testing"
)

func TestRender(t *testing.T) {
	service := CommonService{
		ID:          "id",
		Name:        "name",
		Description: "description",
		Command:     "command",
		Args:        []string{"arg-1", "arg-2"},
	}

	t.Run("macOS service", func(t *testing.T) {
		macOSService, err := render(macOSServiceScript, service)
		if err != nil {
			t.Fatal(err)
		}
		if macOSService != expectmacOSService {
			t.Fatal("macOS service does not meet the expect")
		}
	})

	t.Run("linux service", func(t *testing.T) {
		linuxService, err := render(linuxServiceScript, service)
		if err != nil {
			t.Fatal(err)
		}
		if linuxService != expectLinuxService {
			t.Log(linuxService)
			t.Fatal("linux service does not meet the expect")
		}
	})
}

//go:embed testdata/expect_macos_service.xml
var expectmacOSService string

//go:embed testdata/expect_linux_service.txt
var expectLinuxService string
