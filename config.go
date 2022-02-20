package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/xhoms/gopanosapi"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type ConfigFile struct {
	Panorama struct {
		Name          string `required:"true"`
		Apikey        string `required:"true"`
		Template_name string `required:"true"`
		Tenant_name   string `default:""`
	}
	Debugenabled bool
	Datafile     string `required:"true"`

	Stopfirstone bool
	Stopevery    int

	Data struct {
		Ipsec_gw struct {
			Name string `required:"true"`
			Psk  string `required:"true"`

			V1deadpeerdetection string `default:"yes"`
			V1exchangemode      string
			V1ikecryptoprofile  string

			V2deadpeerdetection string `default:"yes"`
			V2cookievalidation  string
			V2Ikecryptoprofile  string

			Ikeversion     string
			Nattraversal   string `default:"no"`
			Fragmentation  string `default:"no"`
			Passivemode    string
			Localadress    string
			Localinterface string `default:"vlan"`
			Peeraddres     string `default:"dynamic"`
			Localid        string
			Localidtype    string
			Peerid         string
			Peeridtype     string
			Comment        string
		}
		Ipsec_tunnel struct {
			Name               string `required:"true"`
			Ikegateway         string `required:"true"`
			Ipseccryptoprofile string `required:"true"`

			Antireplay       string
			Antireplaywindow string
			Enablegreencap   string
			Copytos          string

			Tunnelinterface string `default:"tunnel"`

			Tunnelmonitor string `default:"no"`
			Proxyid       string
			destinationip string

			Comment string
		}

		Remote_networks struct {
			Name                  string `required:"true"`
			Ecmploadbalancing     string `required:"true"`
			Licensetype           string `required:"true"`
			Location              string `required:"true"`
			Node_termination      string `required:"true"`
			Number_ipsec_tunnels  string `required:"true"`
			Advertisedefaultroute string `default:"no"`
			Sumarizemuroutes      string `default:"no"`
			Donotexportroutes     string `default:"no"`
			Peer_as               string `required:"true"`
			Peer_ip_address       string `required:"true"`
			Local_ip_address      string `required:"true"`
		}
	}
}

type Data struct {
	data_row       []string
	config         ConfigFile
	counter        int
	partialcounter int
	apiconn        gopanosapi.ApiConnector
	datafile       *csv.Reader
	rnEntries      []entries
	tunnelnumber   int
}

// Init
//
// Init :  Load  config file , Open  Datafile and Init the API to Panorama

func (data *Data) Init(configfile string) {
	// Open file confi
	err := configor.Load(&data.config, configfile)
	if err != nil {
		log.Fatal("config: ", err.Error())
	}
	// Open Config file
	f, err := os.Open(data.config.Datafile)

	if err != nil {
		log.Println("No data file is present")
		log.Fatal(err)
	}
	//defer f.Close()
	data.datafile = csv.NewReader(f)

	// open Api to Panorma
	data.apiconn.Init(data.config.Panorama.Name)
	err = data.apiconn.SetKey(data.config.Panorama.Apikey)
	if err != nil {
		log.Fatal(err.Error())
	}
	//setup debgu level
	data.apiconn.Debug(data.config.Debugenabled)

	//init counters which control de the loop
	data.counter = 0
	data.partialcounter = 0
}

func (data *Data) Readnextrow() (err error) {
	data.data_row, err = data.datafile.Read()
	return err

}

func (data *Data) ChangeTunnelNumber(tunnelnumber int) {
	data.tunnelnumber = tunnelnumber
}

func (data *Data) gettunnelnumbers() string {
	return strconv.Itoa(data.tunnelnumber)
}

func (data *Data) ipsplit(ip string, octet int) string {

	s := strings.Split(ip, ".")
	if _, err := strconv.Atoi(s[octet-1]); err == nil {
		return s[octet-1]
	} else {
		os.Exit(-1)
	}
	return ""
}

func (data *Data) getdata(cmd string) string {

	funcMap := template.FuncMap{
		"tunnelnumber": data.gettunnelnumbers,
		"ip_split":     data.ipsplit,
	}
	tmpl, err := template.New("").Funcs(funcMap).Parse(cmd)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	// Run the template to verify the output.
	out := new(strings.Builder)
	if err = tmpl.Execute(out, data.data_row); err != nil {
		log.Fatalf("execution: %s", err)
	}
	return out.String()

}

func (data *Data) GetIkeGWXML(IkeGWXML *string) (ikeGwName string) {

	IkeGWData := data.config.Data.Ipsec_gw

	IkeGW := IKEGateway{
		Name:                data.getdata(IkeGWData.Name),
		PSK:                 data.getdata(IkeGWData.Psk),
		V1DeadPeerDetection: data.getdata(IkeGWData.V1deadpeerdetection),
		V1ExchangeMode:      data.getdata(IkeGWData.V1exchangemode),
		V1IkeCryptoProfile:  data.getdata(IkeGWData.V1ikecryptoprofile),
		V2DeadPeerDetection: data.getdata(IkeGWData.V2deadpeerdetection),
		V2CookieValidation:  data.getdata(IkeGWData.V2cookievalidation),
		V2IkeCryptoProfile:  data.getdata(IkeGWData.V2Ikecryptoprofile),
		Version:             data.getdata(IkeGWData.Ikeversion),
		LocalAddress:        data.getdata(IkeGWData.Localadress),
		LocalInterface:      data.getdata(IkeGWData.Localinterface),
		PeerAddress:         data.getdata(IkeGWData.Peeraddres),
		LocalID:             data.getdata(IkeGWData.Localid),
		LocalIDType:         data.getdata(IkeGWData.Localidtype),
		PeerID:              data.getdata(IkeGWData.Peerid),
		PeerIDType:          data.getdata(IkeGWData.Peeridtype),
		NATTraversal:        data.getdata(IkeGWData.Nattraversal),
		Fragmentation:       data.getdata(IkeGWData.Fragmentation),
		PassiveMode:         data.getdata(IkeGWData.Passivemode),
		Comment:             data.getdata(IkeGWData.Comment),
	}
	output, err := xml.MarshalIndent(IkeGW, "  ", "    ")
	if err != nil {
		log.Fatal("error: %v\n", err)
	}
	//clean XML to avoid errors on Panorama

	*IkeGWXML = strings.Replace(string(output), "<ip>dynamic</ip>", "<dynamic/>", 1)
	*IkeGWXML = strings.Replace(*IkeGWXML, "<local-id></local-id>", "", 1)
	*IkeGWXML = strings.Replace(*IkeGWXML, "<peer-id></peer-id>", "", 1)
	return IkeGW.Name

}

func (data *Data) GetIPSecTunnelXML(ipsecTunnelXML *string) (ipsecTunnelName string) {

	IPsecTunnelData := data.config.Data.Ipsec_tunnel
	IPsecTunnel := IPsecTunnel{
		Name: data.getdata(IPsecTunnelData.Name),
		AutoKey: AutoKey{
			Entry: struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			}{Text: "", Name: data.getdata(data.config.Data.Ipsec_tunnel.Ikegateway)},
		},
		IpsecCryptoProfile:     data.getdata(IPsecTunnelData.Ipseccryptoprofile),
		AntiReplay:             data.getdata(IPsecTunnelData.Antireplay),
		AntiReplayWindow:       data.getdata(IPsecTunnelData.Antireplaywindow),
		EnableGreEncapsulation: data.getdata(IPsecTunnelData.Enablegreencap),
		CopyTos:                data.getdata(IPsecTunnelData.Copytos),
		TunnelInterface:        data.getdata(IPsecTunnelData.Tunnelinterface),
		EnableTunnelMonitor:    data.getdata(IPsecTunnelData.Tunnelmonitor),
		ProxyId:                data.getdata(IPsecTunnelData.Proxyid),
		Destination_ip:         data.getdata(IPsecTunnelData.destinationip),
		Comment:                data.getdata(IPsecTunnelData.Comment),
	}
	output, err := xml.MarshalIndent(IPsecTunnel, "  ", "    ")
	if err != nil {
		log.Fatal("error: %v\n", err)
	}
	*ipsecTunnelXML = string(output)
	return IPsecTunnel.Name
}

func (data *Data) AddRNEntry(tunnelname string) {

	rnEntry := data.config.Data.Remote_networks
	entry := entries{
		Name:                      tunnelname,
		Enable:                    "yes",
		OriginateDefaultRoute:     data.getdata(rnEntry.Advertisedefaultroute),
		SummarizeMobileUserRoutes: data.getdata(rnEntry.Sumarizemuroutes),
		DoNotExportRoutes:         data.getdata(rnEntry.Advertisedefaultroute),
		PeerAs:                    data.getdata(rnEntry.Peer_as),
		PeerIpAddress:             data.getdata(rnEntry.Peer_ip_address),
		LocalIpAddress:            data.getdata(rnEntry.Local_ip_address),
		IpsecTunnel:               tunnelname,
	}
	data.rnEntries = append(data.rnEntries, entry)
}

func (data *Data) GetRNXML(rnXML *string) (rnName string) {

	rnEntry := data.config.Data.Remote_networks
	RemoteNetworkXML := RemoteNetwork{
		Name:                data.getdata(rnEntry.Name),
		Bgp:                 "no",
		Entry:               data.rnEntries,
		Region:              data.getdata(rnEntry.Location),
		LicenseType:         data.getdata(rnEntry.Licensetype),
		SecondaryWanEnabled: "no",
		SpnName:             data.getdata(rnEntry.Node_termination),
		EcmpLoadBalancing:   data.getdata(rnEntry.Ecmploadbalancing),
	}

	output, err := xml.MarshalIndent(RemoteNetworkXML, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	data.rnEntries = nil
	*rnXML = string(output)
	return RemoteNetworkXML.Name
}

func (config *ConfigFile) Readconfigfile(filename string) {

	if err := configor.Load(config, "IPSEC-Prisma_config.yaml"); err != nil {
		fmt.Printf("error!!!!!!!!!!!!!!!!!!!!!")
		fmt.Printf("config: %#v", err.Error())
		os.Exit(2)
	}
}
