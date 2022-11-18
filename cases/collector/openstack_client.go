package collector

import (
	"encoding/json"
	"fmt"
	"os"

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

type OpenStackConfig struct {
	KeystoneIp string `json:"keystone_ip"`
	Username   string `json:"user_name"`
	Password   string `json:"password"`
	DomainId   string `json:"domain_id"`
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
	conf, _ := LoadConfigFile("./cases/openstack_config.json")
	opt := gophercloud.AuthOptions{
		IdentityEndpoint: conf.KeystoneIp,
		Username:         conf.Username,
		Password:         conf.Password,
		DomainID:         conf.DomainId,
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

func LoadConfigFile(file_path string) (OpenStackConfig, error) {
	var config OpenStackConfig
	file, _ := os.Open(file_path)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)
	return config, err
}
