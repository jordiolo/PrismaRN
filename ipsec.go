package main

import "encoding/xml"

type IKEGateway struct {
	XMLName             xml.Name `xml:"entry"`
	Name                string   `xml:"name,attr"`
	PSK                 string   `xml:"authentication>pre-shared-key>key"`
	V1DeadPeerDetection string   `xml:"protocol>ikev1>dpd>enable"`
	V1ExchangeMode      string   `xml:"protocol>ikev1>exchange-mode,omitempty"`
	V1IkeCryptoProfile  string   `xml:"protocol>ikev1>ike-crypto-profile,omitempty"`
	V2DeadPeerDetection string   `xml:"protocol>ikev2>dpd>enable"`
	V2CookieValidation  string   `xml:"protocol>ikev2>require-cookie,omitempty"`
	V2IkeCryptoProfile  string   `xml:"protocol>ikev2>ike-crypto-profile,omitempty"`
	Version             string   `xml:"protocol>version"`
	LocalAddress        string   `xml:"local-address>ip,omitempty"`
	LocalInterface      string   `xml:"local-address>interface,omitempty"`
	PeerDynamicAddress  string   `xml:"peer-address>dinamic,omitempty"`
	PeerAddress         string   `xml:"peer-address>ip,omitempty"`
	LocalID             string   `xml:"local-id>id,omitempty"`
	LocalIDType         string   `xml:"local-id>type,omitempty"`
	PeerID              string   `xml:"peer-id>id,omitempty"`
	PeerIDType          string   `xml:"peer-id>type,omitempty"`
	NATTraversal        string   `xml:"protocol-common>nat-traversal>enable"`
	Fragmentation       string   `xml:"protocol-common>fragmentation>enable"`
	PassiveMode         string   `xml:"protocol-common>passive-mode"`
	Comment             string   `xml:"comment,omitempty"`
}

type AutoKey struct {
	Text  string `xml:",chardata"`
	Entry struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
	} `xml:"entry"`
}

type IPsecTunnel struct {
	XMLName                xml.Name `xml:"entry"`
	Text                   string   `xml:",chardata"`
	Name                   string   `xml:"name,attr"`
	AutoKey                AutoKey  `xml:"auto-key>ike-gateway"`
	IpsecCryptoProfile     string   `xml:"auto-key>ipsec-crypto-profile"`
	EnableTunnelMonitor    string   `xml:"tunnel-monitor>enable"`
	ProxyId                string   `xml:"tunnel-monitor>proxy-id,omitempty"`
	Destination_ip         string   `xml:"tunnel-monitor>destination-ip,omitempty"`
	TunnelInterface        string   `xml:"tunnel-interface"`
	EnableGreEncapsulation string   `xml:"enable-gre-encapsulation,omitempty"`
	CopyTos                string   `xml:"copy-tos,omitempty"`
	AntiReplay             string   `xml:"anti-replay,omitempty"`
	AntiReplayWindow       string   `xml:"anti-replay-window,omitempty"`
	Comment                string   `xml:"comment,omitempty"`
}

type entries struct {
	Text                      string `xml:",chardata"`
	Name                      string `xml:"name,attr"`
	Enable                    string `xml:"protocol>bgp>enable"`
	OriginateDefaultRoute     string `xml:"protocol>bgp>originate-default-route"`
	SummarizeMobileUserRoutes string `xml:"protocol>bgp>summarize-mobile-user-routes"`
	DoNotExportRoutes         string `xml:"protocol>bgp>do-not-export-routes"`
	PeerAs                    string `xml:"protocol>bgp>peer-as"`
	PeerIpAddress             string `xml:"protocol>bgp>peer-ip-address"`
	LocalIpAddress            string `xml:"protocol>bgp>local-ip-address"`
	/*Protocol    Protocol `xml:"protocol"`*/
	IpsecTunnel string `xml:"ipsec-tunnel"`
}

type RemoteNetwork struct {
	XMLName             xml.Name  `xml:"entry"`
	Text                string    `xml:",chardata"`
	Name                string    `xml:"name,attr"`
	Bgp                 string    `xml:"protocol>bgp>enable"`
	Entry               []entries `xml:"link>entry"`
	Region              string    `xml:"region"`
	LicenseType         string    `xml:"license-type"`
	SecondaryWanEnabled string    `xml:"secondary-wan-enabled"`
	SpnName             string    `xml:"spn-name"`
	EcmpLoadBalancing   string    `xml:"ecmp-load-balancing"`
}
