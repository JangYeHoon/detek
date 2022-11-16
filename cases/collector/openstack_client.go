package collector

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/kakao/detek/pkg/detek"
)

const (
	KeyOpenStackClient = "openstack_client"
)

var _ detek.Collector = &OpenStackCollector{}

type OpenStackCollector struct {
	OpenStackConfigPath string
}

func (*OpenStackCollector) GetMeta() detek.CollectorInfo {
	return detek.CollectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "openstack_client",
			Description: "generate openstack client",
			Labels:      []string{"openstack", "client"},
		},
		Required: detek.DependencyMeta{},
		Producing: detek.DependencyMeta{
			KeyOpenStackClient: {Type: detek.TypeOf(&gophercloud.ProviderClient{})},
		},
	}
}

func (c *OpenStackCollector) Do(ctx detek.DetekContext) error {
	opt := gophercloud.AuthOptions{
		IdentityEndpoint: "http://{openstack_ip}:5000/v3",
		Username:         "{user_name}",
		Password:         "{user_password}",
		DomainID:         "default}",
	}

	provider, err := openstack.AuthenticatedClient(opt)

	if err != nil {
		return fmt.Errorf("fail to get openstack client : %w", err)
	}

	if err := ctx.Set(KeyOpenStackClient, provider); err != nil {
		return err
	}

	return nil
}
