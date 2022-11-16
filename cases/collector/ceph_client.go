package collector

import (
	"fmt"

	"github.com/ceph/go-ceph/rados"
	"github.com/kakao/detek/pkg/detek"
)

const (
	KeyCephClient = "ceph_client"
)

var _ detek.Collector = &CephClientCollector{}

type CephClientCollector struct {
	CephConfigPath string
}

func (*CephClientCollector) GetMeta() detek.CollectorInfo {
	return detek.CollectorInfo{
		MetaInfo: detek.MetaInfo{
			ID:          "ceph_client",
			Description: "generate ceph client",
			Labels:      []string{"ceph", "client"},
		},
		Required: detek.DependencyMeta{},
		Producing: detek.DependencyMeta{
			KeyCephClient: {Type: detek.TypeOf(&rados.Conn{})},
		},
	}
}

func (c *CephClientCollector) Do(ctx detek.DetekContext) error {
	conn, err := rados.NewConn()
	if err != nil {
		return fmt.Errorf("fail to set new connect ceph client: %w", err)
	}

	err = conn.ReadConfigFile("")
	if err != nil {
		return fmt.Errorf("fail to get config file: %w", err)
	}

	err = conn.Connect()
	if err != nil {
		return fmt.Errorf("fail to get ceph client: %w", err)
	}

	if err := ctx.Set(KeyCephClient, conn); err != nil {
		return err
	}

	return nil
}
