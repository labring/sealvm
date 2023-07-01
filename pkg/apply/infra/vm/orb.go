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
	"github.com/labring/sealvm/pkg/configs"
	"github.com/labring/sealvm/pkg/utils/exec"
	"github.com/labring/sealvm/pkg/utils/logger"
	"github.com/labring/sealvm/pkg/utils/strings"
	v1 "github.com/labring/sealvm/types/api/v1"
	errors2 "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
	"path"
	strings2 "strings"
)

func NewOrb() Interface {
	return &orb{}
}

type orb struct {
}

func (r *orb) CreateVM(infra *v1.VirtualMachine, host *v1.Host, index int) error {
	cfg := GetCloudInitYamlByRole(infra.Name, host.Role)
	cloudCfg := cloudInit(cfg)
	scriptPath := path.Join(configs.GetEtcDir(infra.Name), fmt.Sprintf("%s.sh", host.Role))
	_ = os.WriteFile(scriptPath, []byte(cloudCfg.toScript()), 0755)
	logger.Info("cloud init to bash success")
	vmID := strings.GetID(infra.Name, host.Role, index)
	if _, err := r.GetById(vmID); err != nil {
		//orb create %[1]s %[2]s && orb -m %[2]s -u root %[3]s
		cmd := fmt.Sprintf("orb create %[1]s %[2]s && orb -m %[2]s -u root %[3]s", host.Image, strings.GetID(infra.Name, host.Role, index), scriptPath)
		logger.Info("executing... %s \n", cmd)
		return exec.Cmd("bash", "-c", cmd)
	}
	return nil
}

func (r *orb) DeleteVM(infra *v1.VirtualMachine, host *v1.VirtualMachineHostStatus) error {
	if _, err := r.GetById(host.ID); err == nil {
		cmd := fmt.Sprintf("orbctl delete -f %s", host.ID)
		return exec.Cmd("bash", "-c", cmd)
	}
	return nil
}

func (r *orb) Get(name, role string, index int) (string, error) {
	cmd := fmt.Sprintf("orb info %s --format json", strings.GetID(name, role, index))
	out, _ := exec.RunBashCmd(cmd)
	if out == "" {
		return "", errors.New("not found instance")
	}
	return out, nil
}

type Image struct {
	Distro  string `json:"distro"`
	Version string `json:"version"`
	Arch    string `json:"arch"`
	Variant string `json:"variant"`
}

type InspectData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Image    Image  `json:"image"`
	Isolated bool   `json:"isolated"`
	Builtin  bool   `json:"builtin"`
	State    string `json:"state"`
}

var imgName = func(i *Image) string {
	return fmt.Sprintf("%s:%s %s", i.Distro, i.Version, i.Arch)
}

func (r *orb) getIPs(name string) ([]string, error) {
	ipv4 := fmt.Sprintf(`orb run -m %s ip -4 addr show | grep global | awk '{print $2}'`, name)
	ipv6 := fmt.Sprintf(`orb run -m %s ip -6 addr show | grep global | awk '{print $2}'`, name)
	ipv4Out, _ := exec.RunBashCmd(ipv4)
	ipv6Out, _ := exec.RunBashCmd(ipv6)
	if ipv4Out == "" && ipv6Out == "" {
		return nil, errors.New("not found ip")
	}
	ips := make([]string, 0)
	if ipv4Out != "" {
		ipv4OutArr := strings.SplitRemoveEmpty(ipv4Out, "/")
		if len(ipv4OutArr) > 0 {
			ipv4Out = ipv4OutArr[0]
		}
		ips = append(ips, strings.TrimSpaceWS(ipv4Out))
	}
	if ipv6Out != "" {
		ipv6OutArr := strings.SplitRemoveEmpty(ipv6Out, "/")
		if len(ipv6OutArr) > 0 {
			ipv6Out = ipv6OutArr[0]
		}
		ips = append(ips, strings.TrimSpaceWS(ipv6Out))
	}
	return ips, nil
}

func (r *orb) InspectByList(name string, role v1.Host, index int) (*v1.VirtualMachineHostStatus, error) {
	type InspectByList []InspectData
	cmd := fmt.Sprintf("orb list --format json")
	out, _ := exec.RunBashCmd(cmd)
	if out == "" {
		return nil, errors.New("not found list instances")
	}
	var outStruct InspectByList
	err := json.Unmarshal([]byte(out), &outStruct)
	if err != nil {
		return nil, errors2.Wrap(err, "decode out json from local vm info failed")
	}

	for _, l := range outStruct {
		if l.Name == strings.GetID(name, role.Role, index) {
			newIPs := make([]string, 0)
			newIPs = append(newIPs, fmt.Sprintf("%s@orb", l.Name))
			ips, _ := r.getIPs(l.Name)
			if len(ips) > 0 {
				newIPs = append(newIPs, ips...)
			}
			return &v1.VirtualMachineHostStatus{
				State:     l.State,
				Role:      role.Role,
				ID:        strings.GetID(name, role.Role, index),
				IPs:       newIPs,
				ImageID:   "",
				ImageName: imgName(&l.Image),
				Capacity:  nil,
				Used:      map[string]string{},
				Mounts:    map[string]string{},
				Index:     index,
			}, nil
		}
	}
	return nil, errors.New("not found this instance")
}

func (r *orb) GetById(name string) (string, error) {
	cmd := fmt.Sprintf("orb info %s --format json", name)
	out, _ := exec.RunBashCmd(cmd)
	if out == "" || strings2.Contains(out, "machine not found") {
		return "", errors.New("not found instance")
	}
	return out, nil
}

func (r *orb) Inspect(name string, role v1.Host, index int) (*v1.VirtualMachineHostStatus, error) {
	//	{
	//	 "id": "01H48FEKCCNFHMB5NRKBFAKZ3R",
	//	 "name": "new-ubuntu",
	//	 "image": {
	//	   "distro": "ubuntu",
	//	   "version": "jammy",
	//	   "arch": "arm64",
	//	   "variant": "default"
	//	 },
	//	 "isolated": false,
	//	 "builtin": false,
	//	 "state": "running"
	//	}
	info, err := r.Get(name, role.Role, index)
	if err != nil {
		return nil, err
	}
	var outStruct map[string]interface{}
	err = json.Unmarshal([]byte(info), &outStruct)
	if err != nil {
		return nil, errors2.Wrap(err, "decode out json from orb info failed")
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

	hostStatus.Capacity = role.Resources
	hostStatus.State, _, _ = unstructured.NestedString(outStruct, "state")
	imageName, _, _ := unstructured.NestedString(outStruct, "image", "distro")
	imageVersion, _, _ := unstructured.NestedString(outStruct, "image", "version")
	imageArch, _, _ := unstructured.NestedString(outStruct, "image", "arch")
	hostStatus.ImageName = fmt.Sprintf("%s:%s %s", imageName, imageVersion, imageArch)

	hostStatus.IPs, _, _ = unstructured.NestedStringSlice(outStruct, "info", hostStatus.ID, "ipv4")
	newIPs := make([]string, 0)
	newIPs = append(newIPs, fmt.Sprintf("%s@orb", hostStatus.ID))
	ips, _ := r.getIPs(hostStatus.ID)
	if len(ips) > 0 {
		newIPs = append(newIPs, ips...)
	}
	hostStatus.IPs = newIPs
	hostStatus.Index = index
	return hostStatus, nil
}

func (r *orb) PingVmsForHosts(infra *v1.VirtualMachine, hosts []v1.VirtualMachineHostStatus) error {
	for _, host := range hosts {
		cmd := fmt.Sprintf(`orb run -m %s ip addr`, host.ID)
		out, _ := exec.RunBashCmd(cmd)
		if out == "" {
			return fmt.Errorf("vm %s is not ready", host.ID)
		}
	}
	return nil
}
