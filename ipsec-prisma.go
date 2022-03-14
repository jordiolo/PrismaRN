package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/xhoms/gopanosapi"
	"io"
	"log"
	"os"
)

func pushPanorama(xpath string, apiconn gopanosapi.ApiConnector, out string, taskname string) {

	log.Println("Configuring ", taskname)
	_, err := apiconn.Config(2, xpath, out)
	if err != nil {
		log.Println("Error configuring ", taskname)
		os.Exit(1)
	}
	if apiconn.LastStatus == gopanosapi.STATUS_ERROR {
		log.Fatal(apiconn.LastResponseMessage)
		os.Exit(1)
	}
}

func createipsec(data Data) string {

	ikeGwXml := ""
	ikeGWName := data.GetIkeGWXML(&ikeGwXml)
	xpath := "/config/devices/entry[@name='localhost.localdomain']/template/entry[@name='" + data.config.Panorama.Template_name + "']/config/devices/entry[@name='localhost.localdomain']/network/ike/gateway"
	//fmt.Println("simulate push ", xpath, ikeGwXml, ikeGWName)
	pushPanorama(xpath, data.apiconn, ikeGwXml, "IPGW :"+ikeGWName)

	tunnelIpsecXml := ""
	tunnelIpsecName := data.GetIPSecTunnelXML(&tunnelIpsecXml)
	xpath = "/config/devices/entry[@name='localhost.localdomain']/template/entry[@name='" + data.config.Panorama.Template_name + "']/config/devices/entry[@name='localhost.localdomain']/network/tunnel/ipsec"
	//fmt.Println("simulate push ", xpath, tunnelIpsecXml, tunnelIpsecName)
	pushPanorama(xpath, data.apiconn, tunnelIpsecXml, "IpSec Tunnel: "+tunnelIpsecName)
	return tunnelIpsecName

}

func createremotenetwork(data Data) {

	var tunnelnumbers []int

	stunnels := data.getdata(data.config.Data.Remote_networks.Number_ipsec_tunnels)
	if err := json.Unmarshal([]byte(stunnels), &tunnelnumbers); err != nil {
		panic(err)
	}
	//fill the tunnel values for the Remote Network
	for _, tunelnumber := range tunnelnumbers {
		data.ChangeTunnelNumber(tunelnumber)
		tunnelname := createipsec(data)
		data.AddRNEntry(tunnelname)
	}
	rnXML := ""
	rnName := data.GetRNXML(&rnXML)

	var xpath string
	if data.config.Panorama.Tenant_name == "" {
		xpath = "/config/devices/entry[@name='localhost.localdomain']/plugins/cloud_services/remote-networks/onboarding"
	} else {
		xpath = "/config/devices/entry[@name='localhost.localdomain']/plugins/cloud_services/" + "multi-tenant/tenants/entry[@name='" + data.config.Panorama.Tenant_name + "']/remote-networks/onboarding"
	}
	pushPanorama(xpath, data.apiconn, rnXML, "IpSec Tunnel: "+rnName)
}

func main() {
	var configfile = flag.String("f", "IPSEC-prisma_config.yaml", "Config File Name")
	flag.Parse()
	var data Data
	data.Init(*configfile)
	for {
		err := data.Readnextrow()
		if err == io.EOF {
			break
		}
		if (data.config.Stopfirstone && (data.counter == 1)) || (data.config.Stopevery == data.partialcounter) {
			data.partialcounter = 0
			fmt.Println("############################################################################## Time to check in Panorama and/or Commit and Push . Press any key to continue")
			fmt.Scanln()
		}
		data.counter++
		data.partialcounter++
		log.Println("##############################################################################  Start Config Remote Network, row: ", data.counter)
		createremotenetwork(data)
	}
	log.Println("##############################################################################   Process finished, Total RN configured", data.counter)

}
