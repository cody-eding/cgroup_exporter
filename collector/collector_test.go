// Copyright 2026 Grand Valley State University
// Copyright 2020 Trey Dockendorf
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"os"
	"log/slog"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

)

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	fixture := filepath.Join(dir, "..", "fixtures")
	CgroupRoot = &fixture
	procFixture := filepath.Join(fixture, "proc")
	ProcRoot = &procFixture
	varTrue := true
	collectProc = &varTrue

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestParseCpuSet(t *testing.T) {
	expected := []string{"0", "1", "2"}
	if cpus, err := parseCpuSet("0-2"); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	} else if !reflect.DeepEqual(cpus, expected) {
		t.Errorf("Unexpected cpus, expected %v got %v", expected, cpus)
	}
	expected = []string{"0", "1", "4", "5", "8", "9"}
	if cpus, err := parseCpuSet("0-1,4-5,8-9"); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	} else if !reflect.DeepEqual(cpus, expected) {
		t.Errorf("Unexpected cpus, expected %v got %v", expected, cpus)
	}
	expected = []string{"1", "3", "5", "7"}
	if cpus, err := parseCpuSet("1,3,5,7"); err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	} else if !reflect.DeepEqual(cpus, expected) {
		t.Errorf("Unexpected cpus, expected %v got %v", expected, cpus)
	}
}
