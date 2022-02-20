 

This script will help to configure Prisma Remote Networks on Panorama 

The scrip will configure :
 * IPsec GW
 * Ipsec Tunnels 
 * Remote GW  ( with or without ECMP ) 
 
 
 Data soure is  provided by a CSV file ( called: datafile)
 
IPSEc-Prisma script will read row by row the data file and will be configure  the tunnels & remote network addording the template defined on the config file

## Usage

$ ./IPSEC-Prisma -f Config_File_Name

Config File Name (default "IPSEC-prisma_config.yaml") 

## IPSEC-Prisma Config file 
 ```yaml
  

panorama:
  name:             ""  
  apikey:           ""  
  template_name:    "Remote_Network_Template"  
  tenant_name:      ""                             # nothing for non-multitenat enviroment

debugenabled: false                                 # boolean: true or false
datafile:   data.csv                                #dat file csv type

stopfirstone:   true                                # boolean: true or false
stopevery:      25                                  # interger



# syntax
# {{index . n }}                return the n colum from the datafile
# {{tunnelnumber}}              return the tunnel number
# {{ip_split (index . n)}}      reurn the IP octect located in the position n   1.2.3.4
# conditional :
#               Return xxxxxx if colum 3 is XXX and return yyyyyy if colum 3 is YYYY
#               {{ if  (eq (index . 3) \"XXX\") }} xxxxxx
#               {{else if (eq (index . 3) \"YYY\") }} yyyyyy
#               {{end}}
#
# number_ipsec_tunnels entry must have a slice with the number of tunnel you want to use , ex: for tunel 1 & 3 : [ 1 3 ]



data:

  ipsec_gw:
      #general Parameters
      name: "IKEGW_Branch_{{index . 0}}_tunnel_{{tunnelnumber}}"  #madatory
      ikeversion: "ikev2"                                   # options: ikev1 , ikev2 , ikev2-ikev2-preferred
      peeraddres: "dynamic"                                 # dynamic or "IP" ; default value = dynamic
      psk: "lfjlñadjflñasjfñalsfjñsal"
      # not used in Prisma
      localadress:                                          #  IP
      localinterface:                                       # vlan for prisma ; default value = vlan

      # Identification

      localid:
      localidtype:                                       # options: none, fqdn , ipaddr , keyid, ufqdn , default = none

      peerid:                   "ipsec_gw_branch_{{index . 0}}_tunnel{{tunnelnumber}}_{{ index . 15}}"
      peeridtype:               "fqdn"                    # options: none, fqdn , ipaddr , keyid, ufqdn , default = none

      #Advanced Config
      passivemode:              "no"                        # options: yes or no
      nattraversal:             "yes"                       # options: yes or no ; default value = no

            # IKE v1 Config
      v1deadpeerdetection:                                  # options: yes or no ; default value = yes
      v1exchangemode:                                       # auto, main , aggressive ; blank if ikev1 is no  selected
      v1ikecryptoprofile:                                   # CryptoProfile's name ; blank if ikev1 is no  selected
      fragmentation:            "no"                        # options: yes or no ; default value = no
            # IKE v2 Config
      v2deadpeerdetection:      "yes"                       # options: yes or no ; default value = yes
      v2cookievalidation:                                   # options: yes or no ; default value = no
      v2ikecryptoprofile:       "IKE_GW_Branches"           # CryptoProfile's name ; blank if ikev2 is no  selected

      comment:                  "Tunnel {{tunnelnumber}} to branch {{ index . 0 }} on {{ index . 15}}"



# IPSec Tunnel config

  ipsec_tunnel:

      name:                     "ipsec_tunnel_branch_{{index . 0}}_tunnel{{tunnelnumber}}_{{ index . 15}}"
      ikegateway:               "IKEGW_Branch_{{index . 0}}_tunnel_{{tunnelnumber}}"
      ipseccryptoprofile:       "IPSEC_Crypto"
      #flags
      antireplay:                                             # value : yes or no , default value : yes
      antireplaywindow:                                       # value : 64,128,256,512,1024,2048,4096 ; default value : 1024
      enablegreencap:                                         # value : yes or no , default value : no
      copytos:                                                # value : yes or no , default value : no

      tunelinterface:                                         # not used in Prisma , default value: tunnel
      #tunnel monitor


      tunnelmonitor:                                          # value : yes or no , default value : no
      proxyid:                                                # value : IP or variable , default value : none
      destinationip:                                          # value : string  , default value : none

      comment:                  "IPSEC tunnel {{tunnelnumber}} to branch {{ index . 0 }} on  {{ index . 15}}"


  remote_networks:

      name:                 "Branch_{{index . 0}}"
      ecmploadbalancing:    "enabled-with-symmetric-return"                                 # This version only support ECMP
      licensetype:          "FWAAS-AGGREGATE"                                               #this parameter is needed
      location:             "spain-central"                                                #  region
      node_termination:     "{{ if  (eq (index . 15) \"NC1\") }}europe-central-A{{     
                             else if (eq (index . 15) \"NC2\") }}europe-central-B{{
                             else if (eq (index . 15) \"NC3\") }}europe-central-C{{
                             else if (eq (index . 15) \"NC4\") }}europe-central-D{{ 
                             else if (eq (index . 15) \"NC5\") }}europe-central-E{{
                             else if (eq (index . 15) \"NC6\") }}europe-central-F{{
                             else if (eq (index . 15) \"NC7\") }}europe-central-G{{
                             else if (eq (index . 15) \"NC8\") }}europe-central-H{{
                             else if (eq (index . 15) \"NC9\") }}europe-central-I{{
                             else if (eq (index . 15) \"NC10\") }}europe-central-J{{
                             else if (eq (index . 15) \"NC11\") }}europe-central-K{{
                             else if (eq (index . 15) \"NC12\") }}europe-central-L{{
                             else if (eq (index . 15) \"NC13\") }}europe-central-M{{
                             else if (eq (index . 15) \"NC14\") }}europe-central-{{
                             else if (eq (index . 15) \"NC15\") }}europe-central-{{
                             else if (eq (index . 15) \"NC16\") }}europe-central-{{
                             else if (eq (index . 15) \"NC17\") }}europe-central-{{
                             else if (eq (index . 15) \"NC18\") }}europe-central-{{
                             else if (eq (index . 15) \"NC19\") }}europe-central-{{
                             else if (eq (index . 15) \"NC20\") }}europe-central-{{
                             end}}"                                                               # spn-name
      #ipsectunnels config
      number_ipsec_tunnels: "{{ if  (eq (index . 2) \"ST\") }}[ 1, 2, 3 ]{{else}}[ 1 , 3 ]{{end}}"
      advertisedefaultroute:                                # value : yes or no , default value : no
      sumarizemuroutes:                                     # value : yes or no , default value : no
      donotexportroutes:                                    # value : yes or no , default value : no
      peer_as:              "65000"
      peer_ip_address:      "{{if (eq (tunnelnumber) \"1\") }}10.1.{{(ip_split (index . 13) 2 )}}.{{(ip_split (index . 13) 3 )}}{{
                            else if (eq (tunnelnumber) \"2\") }}10.2.{{(ip_split (index . 13) 2 )}}.{{(ip_split (index . 13) 3 )}}{{
                            else if (eq (tunnelnumber) \"3\") }}10.3.{{(ip_split (index . 13) 2 )}}.{{(ip_split (index . 13) 3 )}}{{
                            end}}"
      local_ip_address:     "{{if (eq (tunnelnumber) \"1\") }}10.1.255.{{ slice (index . 15) 2}}{{
                            else if (eq (tunnelnumber) \"2\") }}10.2.255.{{ slice (index . 15) 2}}{{
                            else if (eq (tunnelnumber) \"3\") }}10.3.255.{{ slice (index . 15) 2}}{{
                            end}}"
