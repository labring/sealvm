/*
Copyright 2023 cuisongliu@qq.com.

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

package vm

import (
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/labring/sealvm/pkg/utils/exec"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/labring/sealvm/pkg/utils/strings"
	v1 "github.com/labring/sealvm/types/api/v1"
	errors2 "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"strconv"
	strings2 "strings"
)

func NewMultipass() Interface {
	return &multipass{}
}

type multipass struct {
}

func (r *multipass) DeleteVM(infra *v1.VirtualMachine, host *v1.VirtualMachineHostStatus) error {
	if _, err := r.GetById(host.ID); err == nil {
		cmd := fmt.Sprintf("multipass stop %s && multipass delete -p   %s ", host.ID, host.ID)
		return exec.Cmd("bash", "-c", cmd)
	}
	return nil
}

func (r *multipass) Get(name, role string, index int) (string, error) {
	cmd := fmt.Sprintf("multipass info %s --format=json", strings.GetID(name, role, index))
	out, _ := exec.RunBashCmd(cmd)
	if out == "" {
		return "", errors.New("not found instance")
	}
	return out, nil
}

func (r *multipass) List() (string, error) {
	cmd := fmt.Sprintf("multipass list --format json")
	out, _ := exec.RunBashCmd(cmd)
	if out == "" {
		return "", errors.New("not found list instances")
	}
	return out, nil
}

func (r *multipass) GetById(name string) (string, error) {
	cmd := fmt.Sprintf("multipass info %s --format=json", name)
	out, _ := exec.RunBashCmd(cmd)
	if out == "" || strings2.Contains(out, "does not exist") {
		return "", errors.New("not found instance")
	}
	return out, nil
}
func (r *multipass) CreateVM(infra *v1.VirtualMachine, host *v1.Host, index int) error {
	cfg := GetCloudInitYamlByRole(infra.Name, host.Role)
	debugFlag := ""
	if logger.IsDebugMode() {
		debugFlag = "-vvv"
	}
	vmID := strings.GetID(infra.Name, host.Role, index)
	if _, err := r.GetById(vmID); err != nil {
		cmd := fmt.Sprintf("multipass launch --name %s --cpus %s --mem %sG --disk %sG --cloud-init %s %s %s ", strings.GetID(infra.Name, host.Role, index), host.Resources[v1.CPUKey], host.Resources[v1.MEMKey], host.Resources[v1.DISKKey], cfg, debugFlag, host.Image)
		logger.Info("executing... %s \n", cmd)
		return exec.Cmd("bash", "-c", cmd)
	}
	return nil
}

func (r *multipass) Inspect(name string, role v1.Host, index int) (*v1.VirtualMachineHostStatus, error) {
	info, err := r.Get(name, role.Role, index)
	if err != nil {
		return nil, err
	}
	var outStruct map[string]interface{}
	err = json.Unmarshal([]byte(info), &outStruct)
	if err != nil {
		return nil, errors2.Wrap(err, "decode out json from multipass info failed")
	}
	hostStatus := &v1.VirtualMachineHostStatus{
		State:     "",
		Role:      role.Role,
		ID:        strings.GetID(name, role.Role, index),
		IPs:       nil,
		ImageID:   "",
		ImageName: "",
		Capacity:  nil,
		Used:      map[string]string{},
		Mounts:    map[string]string{},
	}

	memUsed, _, _ := unstructured.NestedInt64(outStruct, "info", hostStatus.ID, "memory", "used")
	diskUsed, _, _ := unstructured.NestedString(outStruct, "info", hostStatus.ID, "disks", "sda1", "used")
	cpuUsed, _, _ := unstructured.NestedSlice(outStruct, "info", hostStatus.ID, "load")
	logger.Debug("memUsed:", memUsed, "diskUsed:", diskUsed, "cpuUsed:", cpuUsed)
	hostStatus.Used[v1.MEMKey] = humanize.Bytes(uint64(memUsed))
	diskUsedInt, _ := strconv.Atoi(diskUsed)
	hostStatus.Used[v1.DISKKey] = humanize.Bytes(uint64(diskUsedInt))
	hostStatus.Used[v1.CPUKey] = fmt.Sprintf("%v", cpuUsed)
	hostStatus.Capacity = role.Resources
	hostStatus.State, _, _ = unstructured.NestedString(outStruct, "info", hostStatus.ID, "state")
	hostStatus.ImageID, _, _ = unstructured.NestedString(outStruct, "info", hostStatus.ID, "image_hash")
	hostStatus.ImageName, _, _ = unstructured.NestedString(outStruct, "info", hostStatus.ID, "release")
	hostStatus.IPs, _, _ = unstructured.NestedStringSlice(outStruct, "info", hostStatus.ID, "ipv4")
	newIPs := make([]string, 0)
	if len(hostStatus.IPs) > 0 {
		for _, ip := range hostStatus.IPs {
			if strings2.HasPrefix(ip, "172.17") || strings2.HasPrefix(ip, "10.96") {
				continue
			} else {
				newIPs = append(newIPs, ip)
			}
		}
	}
	hostStatus.IPs = newIPs
	hostStatus.Index = index
	mounts, _, _ := unstructured.NestedMap(outStruct, "info", hostStatus.ID, "mounts")
	for k := range mounts {
		hostMount, _, _ := unstructured.NestedString(outStruct, "info", hostStatus.ID, "mounts", k, "source_path")
		hostStatus.Mounts[hostMount] = k
	}

	return hostStatus, nil
}
