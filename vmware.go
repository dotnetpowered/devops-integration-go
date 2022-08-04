package main

import (
	"context"
	"net/url"

	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
)

// NewClient creates a vim25.Client for use
func NewClient(ctx context.Context, cfg VmWareConfig) (*vim25.Client, error) {
	// Parse URL from string
	u, err := soap.ParseURL(cfg.Url)
	if err != nil {
		return nil, err
	}

	u.User = url.UserPassword(cfg.User, cfg.Password)

	// Share govc's session cache
	s := &cache.Session{
		URL:      u,
		Insecure: true,
	}

	c := new(vim25.Client)
	err = s.Login(ctx, c, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func vmware_collector(cfg VmWareConfig) ([]Host, error) {

	vms, err := getVMs(cfg)

	if err != nil {
		return nil, err
	}

	var hosts []Host
	for _, vm := range vms {
		host := Host{
			Name:        vm.Summary.Guest.HostName,
			Source:      "VMWare",
			OS:          vm.Summary.Config.GuestFullName,
			Ip:          vm.Summary.Guest.IpAddress,
			Status:      string(vm.Summary.Runtime.PowerState),
			Description: vm.Summary.Config.Name,
			NumCpu:      int(vm.Summary.Config.NumCpu),
			MemSize:     int(vm.Summary.Config.MemorySizeMB),
			Notes:       vm.Summary.Config.Annotation,
			Uptime:      int(vm.Summary.QuickStats.UptimeSeconds),
		}
		hosts = append(hosts, host)
		// Print summary per vm (see also: govc/vm/info.go)
		//fmt.Printf("%s: %s %s\n", vm.Summary.Config.Name, vm.Summary.Config.GuestFullName, vm.Summary.Guest.HostName)
	}
	return hosts, nil
}

func getVMs(cfg VmWareConfig) ([]mo.VirtualMachine, error) {
	var c *vim25.Client

	var err error

	ctx := context.Background()
	c, err = NewClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	// Create view of VirtualMachine objects
	m := view.NewManager(c)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var vms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		return nil, err
	} else {
		return vms, nil
	}

}
