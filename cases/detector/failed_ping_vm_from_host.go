package detector

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/kakao/detek/cases/collector"
	"github.com/kakao/detek/pkg/detek"
)

var _ detek.Detector = &FailedPingVMFromHost{}

type FailedPingVMFromHost struct{}

func (*FailedPingVMFromHost) GetMeta() detek.DetectorInfo {
	return detek.DetectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "faild_ping_vm_from_host",
			Description: "Finding VMs that failed pings from the server",
			Labels:      []string{"openstack", "vm", "ping"},
		},
		Required: detek.DependencyMeta{
			collector.KeyOpenStackServerList: {Type: detek.TypeOf([]servers.Server{})},
		},
		Level: detek.Error,
		IfHappened: detek.Description{
			Explanation: "some of VMs failed pings",
			Solution:    "check security gorup or flow rule",
		},
	}
}

func (*FailedPingVMFromHost) Do(ctx detek.DetekContext) (*detek.ReportSpec, error) {
	serverList, err := detek.Typing[[]servers.Server](
		ctx.Get(collector.KeyOpenStackServerList, nil),
	)

	if err != nil {
		return nil, err
	}

	type Problem struct {
		Id, Name, Ip, Reason string
	}
	problems := []Problem{}

	for _, s := range serverList {
		if s.Status == "ACTIVE" {
			var ipList []string
			for _, val := range s.Addresses {
				str := fmt.Sprintf("%v", val)
				first_num := strings.LastIndex(str, "addr:") + 5
				last_num := strings.LastIndex(str, "version") - 1
				ipList = append(ipList, str[first_num:last_num])
			}

			for _, value := range ipList {
				ping_result, _ := exec.Command("ping", value, "-c 1").Output()
				if strings.Contains(string(ping_result), "Destination Host Unreachable") ||
					strings.Contains(string(ping_result), "Request Timed Out") {
					problems = append(problems, Problem{
						Id:     s.ID,
						Name:   s.Name,
						Ip:     value,
						Reason: "ping fail",
					})
				}
			}
		}
	}

	return &detek.ReportSpec{
		HasPassed: len(problems) == 0,
		Problem: detek.JSONableData{
			Description: "Failed ping",
			Data:        problems,
		},
	}, nil
}
