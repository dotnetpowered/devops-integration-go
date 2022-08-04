package main

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func main() {

	var cfg Config

	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Collecting VMWare...")
	vms, _ := vmware_collector(cfg.VmWare)
	printHosts(vms)
	saveHosts(vms, cfg.Database.Dsn)

	fmt.Println("Collecting Octopus...")
	targets := octopus_collector(cfg.OctopusDeploy)
	printHosts(targets)
	saveHosts(targets, cfg.Database.Dsn)

	fmt.Println("Collecting Zabbix...")
	hosts := zabbix_collector(cfg.Zabbix)
	printHosts(hosts)
	saveHosts(hosts, cfg.Database.Dsn)

	fmt.Println("Finished.")
}
