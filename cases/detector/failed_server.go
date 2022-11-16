package detector

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/kakao/detek/cases/collector"
	"github.com/kakao/detek/pkg/detek"
)

var _ detek.Detector = &FailedServer{}

type FailedServer struct{}

func (*FailedServer) GetMeta() detek.DetectorInfo {
	return detek.DetectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "failed_server",
			Description: "check if there is a server with a 'Failed' status",
			Labels:      []string{"openstack", "server"},
		},
		Level: detek.Error,
		IfHappened: detek.Description{
			Explanation: `some of servers are in a "Failed" status`,
			Solution:    `check why servers are failed`,
		},
		Required: detek.DependencyMeta{
			collector.KeyOpenStackServerList: {Type: detek.TypeOf([]servers.Server{})},
		},
	}
}

func (i *FailedServer) Do(ctx detek.DetekContext) (*detek.ReportSpec, error) {
	serverList, err := detek.Typing[[]servers.Server](
		ctx.Get(collector.KeyOpenStackServerList, nil),
	)

	if err != nil {
		return nil, err
	}

	type Problem struct {
		Id, Name, Reason string
	}
	problems := []Problem{}

	for _, s := range serverList {
		if s.Status == "ERROR" {
			problems = append(problems, Problem{
				Id:     s.ID,
				Name:   s.Name,
				Reason: "server status error",
			})
		}
	}

	report := &detek.ReportSpec{
		HasPassed:  len(problems) == 0,
		Attachment: []detek.JSONableData{{Description: "# of Servers", Data: len(serverList)}},
	}
	if len(problems) != 0 {
		report.HasPassed = false
		report.Problem = detek.JSONableData{
			Description: "Failed server list",
			Data:        problems,
		}
	}

	return report, nil
}
