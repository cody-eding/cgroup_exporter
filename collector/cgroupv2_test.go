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
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestGetStatv2(t *testing.T) {
	_, err := getStatv2("swapcached", "/dne")
	if err == nil {
		t.Errorf("Expected error with /dne but none given")
	}
	path := filepath.Join(*CgroupRoot, "system.slice")
	_, err = getStatv2("swapcached", path)
	if err == nil {
		t.Errorf("Expected error with /dne but none given")
	}
	path = filepath.Join(*CgroupRoot, "user.slice/user-20821.slice/memory.max")
	_, err = getStatv2("swapcached", path)
	if err == nil {
		t.Errorf("Expected error with single value file but none given")
	}
	path = filepath.Join(*CgroupRoot, "stat.invalid")
	_, err = getStatv2("nan", path)
	if err == nil {
		t.Errorf("Expected error with stat.invalid but none given")
	}
	path = filepath.Join(*CgroupRoot, "user.slice/user-20821.slice/memory.stat")
	_, err = getStatv2("dne", path)
	if err == nil {
		t.Errorf("Expected error when stat key missing but none given")
	}
	stat, err := getStatv2("swapcached", path)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if stat != 0 {
		t.Errorf("Unexpectd value: %v", stat)
	}
}

func TestCollectv2Error(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	exporter := NewExporter([]string{"/dne"}, logger)
	metrics, err := exporter.collectv2()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
		return
	}
	if val := len(metrics); val != 1 {
		t.Errorf("Unexpected number of metrics, got %d expected 1", val)
		return
	}
	if val := metrics[0].err; val != true {
		t.Errorf("Unexpected value for err, got %v", val)
	}
}

func TestCollectv2SLURM(t *testing.T) {
	varLen := 100
	collectProcMaxExec = &varLen
	PidGroupPath = func(pid int) (string, error) {
		if pid == 49276 {
			return "/system.slice/slurmstepd.scope/job_4/step_0/user/task_0", nil
		}
		if pid == 43310 {
			return "/system.slice/slurmstepd.scope/system", nil
		}
		return "", fmt.Errorf("Could not find cgroup path for %d", pid)
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	exporter := NewExporter([]string{"/system.slice/slurmstepd.scope"}, logger)
	metrics, err := exporter.collectv2()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
		return
	}
	if val := len(metrics); val != 1 {
		t.Errorf("Unexpected number of metrics, got %d expected 1", val)
		return
	}
	var m CgroupMetric
	for _, metric := range metrics {
		if metric.jobid == "4" {
			m = metric
		}
	}
	if m.jobid == "" {
		t.Errorf("Metrics with jobid=4 not found")
		return
	}
	if val := m.name; val != "/system.slice/slurmstepd.scope/job_4" {
		t.Errorf("Unexpected value for name, got %v", val)
	}
	if val := m.cpuUser; val != 0.049043 {
		t.Errorf("Unexpected value for cpuUser, got %v", val)
	}
	if val := m.cpuSystem; val != 0.077642 {
		t.Errorf("Unexpected value for cpuSystem, got %v", val)
	}
	if val := m.cpuTotal; val != 0.126686 {
		t.Errorf("Unexpected value for cpuTotal, got %v", val)
	}
	if val := m.cpus; val != 1 {
		t.Errorf("Unexpected value for cpus, got %v", val)
	}
	if val := m.memoryRSS; val != 2777088 {
		t.Errorf("Unexpected value for memoryRSS, got %v", val)
	}
	if val := m.memoryUsed; val != 5660672 {
		t.Errorf("Unexpected value for memoryUsed, got %v", val)
	}
	if val := m.memoryTotal; val != 1835008000 {
		t.Errorf("Unexpected value for memoryTotal, got %v", val)
	}
	if val := m.memoryFailCount; val != 0 {
		t.Errorf("Unexpected value for memoryFailCount, got %v", val)
	}
	if val := m.memswUsed; val != 0 {
		t.Errorf("Unexpected value for swapUsed, got %v", val)
	}
	if val := m.memswTotal; val != 1835008000 {
		t.Errorf("Unexpected value for swapTotal, got %v", val)
	}
	if val := m.uid; val != "20821" {
		t.Errorf("Unexpected value for uid, got %v", val)
	}
	if val := m.jobid; val != "4" {
		t.Errorf("Unexpected value for jobid, got %v", val)
	}
}
