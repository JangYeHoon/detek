package collector

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/kakao/detek/pkg/detek"
)

const (
	KeyOpenStackServerList       = "openstack_serverlist"
	KeyOpenStackLoadBalancerList = "openstack_loadbalancerlist"
)

var _ detek.Collector = &OpenStackCoreCollector{}

type OpenStackCoreCollector struct{}

func (*OpenStackCoreCollector) GetMeta() detek.CollectorInfo {
	return detek.CollectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "openstack_core",
			Description: "collect core resources from openstack",
			Labels:      []string{"openstack", "core", "manifest"},
		},
		Required: detek.DependencyMeta{
			KeyOpenStackClient: {Type: detek.TypeOf(&gophercloud.ProviderClient{})},
		},
		Producing: detek.DependencyMeta{
			KeyOpenStackServerList:       {Type: detek.TypeOf([]servers.Server{})},
			KeyOpenStackLoadBalancerList: {Type: detek.TypeOf([]servers.Server{})},
		},
	}
}

func (*OpenStackCoreCollector) Do(dctx detek.DetekContext) error {
	c, err := detek.Typing[*gophercloud.ProviderClient](
		dctx.Get(KeyOpenStackClient, nil),
	)
	if err != nil {
		return fmt.Errorf("fail to get openstack client: %w", err)
	}

	client, err := openstack.NewComputeV2(c, gophercloud.EndpointOpts{
		Region:       "RegionOne",
		Name:         "nova",
		Type:         "compute",
		Availability: gophercloud.AvailabilityInternal,
	})
	if err != nil {
		return fmt.Errorf("fail to get openstack compute client: %w", err)
	}
	var errs = &multierror.Error{}

	serverPager := servers.List(client, servers.ListOpts{})
	err = serverPager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)
		errs = multierror.Append(errs, err)
		errs = multierror.Append(errs,
			dctx.Set(KeyOpenStackServerList, serverList),
		)
		return true, nil
	})

	if err != nil {
		return fmt.Errorf("fail get server list : %w", err)
	}

	return errs.ErrorOrNil()
}
