package main

import (
	"fmt"
	"strings"

	"github.com/rday/zabbix"
)

// Find and return a single host object by name
func ZabbixGetHosts(api *zabbix.API) ([]zabbix.ZabbixHost, error) {
	params := make(map[string]interface{}, 0)
	filter := make(map[string]string, 0)
	inventory := []string{"os"}
	output := []string{"name", "description", "host", "status", "lastaccess"}
	template := []string{"templateid", "name"}
	params["filter"] = filter
	params["output"] = output
	params["selectInventory"] = inventory
	params["selectTags"] = "extend"
	params["selectInterfaces"] = "extend"
	params["selectGroups"] = "extend"
	params["selectParentTemplates"] = template
	params["inheritedTags"] = "true"
	//	params["templated_hosts"] = 1
	ret, err := api.Host("get", params)

	// This happens if there was an RPC error
	if err != nil {
		return nil, err
	}

	// If our call was successful
	if len(ret) > 0 {
		return ret, err
	}

	// This will be the case if the RPC call was successful, but
	// Zabbix had an issue with the data we passed.
	return nil, &zabbix.ZabbixError{Code: 0, Message: "", Data: "Error getting hosts"}
}

func zabbix_collector(cfg ZabbixConfig) []Host {
	api, err := zabbix.NewAPI(cfg.Url, cfg.User, cfg.Password)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	versionresult, err := api.Version()
	if err != nil {
		fmt.Println(err)
	}

	_, err = api.Login()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(" Connected to API using ", versionresult)

	zabbix_hosts, _ := ZabbixGetHosts(api)

	var hosts []Host
	for i := 0; i < len(zabbix_hosts); i++ {
		var h Host
		zh := zabbix_hosts[i]
		dnsName, ipAddr := getIpAddr(zh["host"].(string))
		h.Name = dnsName
		h.Source = "Zabbix"
		h.Ip = ipAddr
		h.Description = zh["description"].(string)
		h.Status = zh["status"].(string)

		// get template name
		parentTemplate := zh["parentTemplates"].([]interface{})
		if len(parentTemplate) > 0 {
			switch v := parentTemplate[0].(type) {
			case (map[string]interface{}):
				h.Template = v["name"].(string)
			}
		}
		// get OS string from inventory
		inventory := zh["inventory"]
		switch v := inventory.(type) {
		case (map[string]interface{}):
			// remove machine name from OS string
			os := strings.Split(v["os"].(string), " ")
			if len(os) > 2 {
				os = remove(os, 1)
			}
			h.OS = strings.Join(os, " ")
		}
		groups := zh["groups"].([]interface{})
		for _, g := range groups {
			grp := g.(map[string]interface{})
			h.Tags = append(h.Tags, grp["name"].(string))
		}
		hosts = append(hosts, h)
	}
	return hosts

}
