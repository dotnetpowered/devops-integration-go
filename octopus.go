package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func getOctopusClient(cfg OctopusConfig) *octopusdeploy.Client {
	octopusURL := cfg.Url //os.Getenv("OCTOPUS_URL")
	apiKey := cfg.ApiKey  //os.Getenv("OCTOPUS_APIKEY")

	if isEmpty(octopusURL) || isEmpty(apiKey) {
		log.Fatal("Please make sure to set the env variables 'OCTOPUS_URL' and 'OCTOPUS_APIKEY' before running this test")
	}

	apiURL, err := url.Parse(octopusURL)
	if err != nil {
		_ = fmt.Errorf("error parsing URL for Octopus API: %v", err)
		return nil
	}

	// NOTE: You can direct traffic through a proxy trace like Fiddler
	// Everywhere by preconfiguring the client to route traffic through a
	// proxy.

	// proxyStr := "http://127.0.0.1:5555"
	// proxyURL, err := url.Parse(proxyStr)
	// if err != nil {
	// 	log.Println(err)
	// }

	// tr := &http.Transport{
	// 	Proxy: http.ProxyURL(proxyURL),
	// }
	// httpClient := http.Client{Transport: tr}

	octopusClient, err := octopusdeploy.NewClient(nil, apiURL, apiKey, emptyString)
	if err != nil {
		log.Fatal(err)
	}

	return octopusClient
}

func getMachines(client *octopusdeploy.Client) *octopusdeploy.DeploymentTargets {
	query := octopusdeploy.MachinesQuery{
		//CommunicationStyles: []string{"Kubernetes"},
	}
	resources, err := client.Machines.Get(query)
	_ = err
	return resources
}

func octopus_collector(cfg OctopusConfig) []Host {
	client := getOctopusClient(cfg)
	machines := getMachines(client)

	var hosts []Host
	for i := 0; i < len(machines.Items); i++ {
		var h Host
		machine := machines.Items[i]
		uri, err := url.Parse(machine.URI)
		_ = err
		host_parts := strings.Split(uri.Host, ":")
		hostName := host_parts[0]

		dnsName, ipAddr := getIpAddr(hostName)
		h.Name = dnsName
		h.Ip = ipAddr
		h.OS = machine.OperatingSystem
		h.Tags = machine.Roles
		h.Status = machine.Status
		h.Template = machine.ShellName
		h.Source = "OctopusDeploy"
		//h.template

		hosts = append(hosts, h)

	}

	p, _ := client.Projects.GetAll()

	fmt.Println(p[0].Description)

	query := octopusdeploy.DeploymentQuery{Skip: 0, Take: 10}
	d, _ := client.Deployments.GetDeployments(nil, &query)
	fmt.Println(d.Items[0].Name)

	return hosts

}
