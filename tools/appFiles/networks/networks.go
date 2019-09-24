// Package networks provides information on the computers network and ip scanning tools.
package networks

import (
	"errors"
	"math"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"bytes"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
)

const (
	NETWORK_PATH = "/etc/network"
)

type Network struct {
	Address            string  `json:"Address"`
	CIDR               string  `json:"CIDR"`
	AvailableAddresses float64 `json:"AvailableAddresses"`
	Name               string  `json:"Name"`
	MAC                string  `json:"MAC"`
	Status             string  `json:"Status"`
	StatusInt          int     `json:"StatusInt"`
	Subnet             string  `json:"Subnet"`
}

type LinuxNetworkFile struct {
	Name           string        `json:"Name"`
	DHCP           bool          `json:"DHCP"`
	Address        string        `json:"Address"`
	Subnet         string        `json:"Subnet"`
	Gateway        string        `json:"Gateway"`
	HWAddress      string        `json:"HWAddress"`
	DNSNameserver1 string        `json:"DNSNameserver1"`
	DNSNameserver2 string        `json:"DNSNameserver2"`
	StaticRoutes   []StaticRoute `json:"StaticRoutes"`
	Errors         struct {
		Address        string `json:"Address"`
		Subnet         string `json:"Subnet"`
		Gateway        string `json:"Gateway"`
		DNSNameserver1 string `json:"DNSNameserver1"`
		DNSNameserver2 string `json:"DNSNameserver2"`
	} `json:"Errors"`
	WPASupplicant WPASupplicant `json:"WPASupplicant"`
}

type StaticRoute struct {
	Address   string `json:"Address"`
	Netmask   string `json:"Netmask"`
	Gateway   string `json:"Gateway"`
	Interface string `json:"Interface"`
	Errors    struct {
		Address string `json:"Address"`
		Netmask string `json:"Netmask"`
		Gateway string `json:"Gateway"`
	} `json:"Errors"`
}

type Host struct {
	Address string `json:"Address"`
	Domain  string `json:"Domain"`
	Errors  struct {
		Address string `json:"Address"`
		Domain  string `json:"Domain"`
	} `json:"Errors"`
}

type WPASupplicant struct {
	Enabled bool `json:"Enabled"`
	/* SSID              string `json:"SSID"`
	ScanSSID          string `json:"ScanSSID"` */
	KeyMgmt string `json:"KeyMgmt"`
	/* Pairwise          string `json:"Pairwise"`
	Group             string `json:"Group"`
	PSK               string `json:"PSK"` */
	EAP        string `json:"EAP"`
	Identity   string `json:"Identity"`
	Password   string `json:"Password"`
	CACert     string `json:"CACert"`
	CACertOld  string `json:"CACertOld"`
	CACertFile string `json:"CACertFile"`
	Phase2     string `json:"Phase2"`
	/*ClientCert        string `json:"ClientCert"`
	PrivateKey        string `json:"PrivateKey"`
	PrivateKeyPasswd  string `json:"PrivateKeyPasswd"`
	Phase1            string `json:"Phase1"`
	CACert2           string `json:"CACert2"`
	ClientCert2       string `json:"ClientCert2"`
	PrivateKey2       string `json:"PrivateKey2"`
	PrivateKey2Passwd string `json:"PrivateKey2Passwd"` */
	Errors struct {
		/* SSID              string `json:"SSID"`
		ScanSSID          string `json:"ScanSSID"` */
		KeyMgmt string `json:"KeyMgmt"`
		/* Pairwise          string `json:"Pairwise"`
		Group             string `json:"Group"`
		PSK               string `json:"PSK"` */
		EAP        string `json:"EAP"`
		Identity   string `json:"Identity"`
		Password   string `json:"Password"`
		CACert     string `json:"CACert"`
		CACertFile string `json:"CACertFile"`
		Phase2     string `json:"Phase2"`
		/*ClientCert        string `json:"ClientCert"`
		PrivateKey        string `json:"PrivateKey"`
		PrivateKeyPasswd  string `json:"PrivateKeyPasswd"`
		Phase1            string `json:"Phase1"`
		CACert2           string `json:"CACert2"`
		ClientCert2       string `json:"ClientCert2"`
		PrivateKey2       string `json:"PrivateKey2"`
		PrivateKey2Passwd string `json:"PrivateKey2Passwd"` */
	} `json:"Errors"`
}

type DHCPDConf struct {
	Enabled          bool   `json:"Enabled"`
	DefaultLeaseTime string `json:"DefaultLeaseTime"`
	MaxLeaseTime     string `json:"MaxLeaseTime"`
	Subnet           string `json:"Subnet"`
	Netmask          string `json:"Netmask"`
	RangeFrom        string `json:"RangeFrom"`
	RangeTo          string `json:"RangeTo"`
	OptionRouter     string `json:"OptionRouter"`
	OptionSubnetMask string `json:"OptionSubnetMask"`
	OptionDNS1       string `json:"OptionDNS1"`
	OptionDNS2       string `json:"OptionDNS2"`
	OptionNTPServer  string `json:"OptionNTPServer"`
	Errors           struct {
		Subnet           string `json:"Subnet"`
		Netmask          string `json:"Netmask"`
		RangeFrom        string `json:"RangeFrom"`
		RangeTo          string `json:"RangeTo"`
		OptionRouter     string `json:"OptionRouter"`
		OptionSubnetMask string `json:"OptionSubnetMask"`
		OptionDNS1       string `json:"OptionDNS1"`
		OptionDNS2       string `json:"OptionDNS2"`
		OptionNTPServer  string `json:"OptionNTPServer"`
	} `json:"Errors"`
}

type ScanResult struct {
	sync.RWMutex
	Devices []string
}

type CurrentLinuxNetworksSync struct {
	sync.RWMutex
	Networks []LinuxNetworkFile
}

type CurrentLinuxNICDHCPStatusSync struct {
	sync.RWMutex
	CurrentStatus map[string]bool
}

type ScanResultsCallback func(sr *ScanResult)

type DNSRecordAdd func(domain string, address string)

var AddDNSRecord DNSRecordAdd
var CurrentLinuxNetworks CurrentLinuxNetworksSync
var CurrentLinuxNICDHCPStatus CurrentLinuxNICDHCPStatusSync

func init() {
	CurrentLinuxNetworks.Lock()
	CurrentLinuxNetworks.Networks, _ = GetLinuxNetworks()

	CurrentLinuxNICDHCPStatus.Lock()
	CurrentLinuxNICDHCPStatus.CurrentStatus = make(map[string]bool)

	for i := range CurrentLinuxNetworks.Networks {
		network := CurrentLinuxNetworks.Networks[i]
		if network.DHCP { //Initialize the NIC to true assuming it is up
			CurrentLinuxNICDHCPStatus.CurrentStatus[network.Name] = true
		}
	}

	CurrentLinuxNICDHCPStatus.Unlock()
	CurrentLinuxNetworks.Unlock()
}

// GetNetworks - Get Networks for gateway
func GetNetworks() (networks []Network) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, inf := range interfaces {
			interfac := inf
			addresses, err := interfac.Addrs()
			if err == nil {
				for _, addr := range addresses {

					if addr.String() == "127.0.0.1/8" || len(interfac.HardwareAddr.String()) > 18 {
						continue
					}
					if addr.String()[:4] == "169." { //DHCP that has not been assigned.
						continue
					}

					ip, _, err := net.ParseCIDR(addr.String())
					if err == nil {

						ones, _ := ip.DefaultMask().Size()

						if ip.To4() == nil {
							continue
						}

						var n Network
						n.Address = ip.String()
						n.Subnet = ip.DefaultMask().String()
						n.CIDR = addr.String()
						n.Name = interfac.Name
						n.MAC = interfac.HardwareAddr.String()

						n.StatusInt = int(interfac.Flags)
						n.Status = interfac.Flags.String()
						n.AvailableAddresses = math.Pow(2, float64(32-ones))

						networks = append(networks, n)
					}
				}
			}
		}
	}
	// session_functions.Log("Networks", fmt.Sprintf("%+v", networks))
	return
}

// GetCIDR - Get CIDR address
func GetCIDR(cidr string, callback ScanResultsCallback) {

	sr := new(ScanResult)
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err == nil {
		_, bits := ip.DefaultMask().Size()

		if bits >= 16 {

			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
				sr.Lock()
				sr.Devices = append(sr.Devices, ip.String())
				sr.Unlock()
			}

		}
		callback(sr)
	}
}

// GetAll - Get all networks
func GetAll(callback ScanResultsCallback) {
	sr := new(ScanResult)
	var wg sync.WaitGroup
	for _, network := range GetNetworks() {
		wg.Add(1)
		n := network
		go func() {
			GetCIDR(n.CIDR, func(srr *ScanResult) {
				sr.Lock()
				sr.Devices = append(sr.Devices, srr.Devices...)
				sr.Unlock()
				wg.Done()
			})
		}()
	}
	wg.Wait()
	callback(sr)
}

// GetLinuxDHCPServer - Returns Linux DHCP server Config
func GetLinuxDHCPServer() (dhcp DHCPDConf, err error) {
	if extensions.DoesFileExist("/etc/dhcp/dhcpd.conf") == false {
		return
	}
	data, err := extensions.ReadFile("/etc/dhcp/dhcpd.conf")
	if err != nil {
		session_functions.Log("Error", "Failed to read /etc/dhcp/dhcpd.conf:  "+err.Error())
		return
	}

	if extensions.DoesFileExist("/etc/default/isc-dhcp-server") == false {
		return
	}

	data1, err := extensions.ReadFile("/etc/default/isc-dhcp-server")
	if err != nil {
		session_functions.Log("Error", "Failed to read /etc/default/isc-dhcp-server:  "+err.Error())
		return
	}

	lines := strings.Split(string(data), "\n")
	for i := range lines {
		line := strings.TrimRight(lines[i], ";")
		if len(line) == 0 {
			continue
		}
		if strings.Contains(line, "ddns-update-style") || strings.Contains(line, "authoritative") || strings.Contains(line, "}") {
			continue
		} else if strings.Contains(line, "default-lease-time ") {
			dhcp.DefaultLeaseTime = strings.Replace(line, "default-lease-time ", "", -1)
		} else if strings.Contains(line, "max-lease-time ") {
			dhcp.MaxLeaseTime = strings.Replace(line, "max-lease-time ", "", -1)
		} else if strings.Contains(line, " {") {
			lineDetails := strings.Split(strings.Replace(line, " {", "", -1), " ")
			if len(lineDetails) > 1 {
				dhcp.Subnet = lineDetails[1]
			}
			if len(lineDetails) > 3 {
				dhcp.Netmask = lineDetails[3]
			}
		} else if strings.Contains(line, "range") {
			lineDetails := strings.Split(strings.Replace(strings.Trim(line, " "), "\t", "", -1), " ")
			if len(lineDetails) > 1 {
				dhcp.RangeFrom = lineDetails[1]
			}
			if len(lineDetails) > 2 {
				dhcp.RangeTo = lineDetails[2]
			}
		} else if strings.Contains(line, "option routers") {
			lineDetails := strings.Split(strings.Replace(strings.Trim(line, " "), "\t", "", -1), " ")
			if len(lineDetails) > 2 {
				dhcp.OptionRouter = lineDetails[2]
			}
		} else if strings.Contains(line, "option subnet-mask") {
			lineDetails := strings.Split(strings.Replace(strings.Trim(line, " "), "\t", "", -1), " ")
			if len(lineDetails) > 2 {
				dhcp.OptionSubnetMask = lineDetails[2]
			}
		} else if strings.Contains(line, "option domain-name-servers") {
			lineDetails := strings.Split(strings.Replace(strings.Trim(line, " "), "\t", "", -1), " ")
			if len(lineDetails) > 2 {
				dhcp.OptionDNS1 = strings.TrimRight(lineDetails[2], ",")
			}
			if len(lineDetails) > 3 {
				dhcp.OptionDNS2 = lineDetails[3]
			}
		} else if strings.Contains(line, "option ntp-servers") {
			lineDetails := strings.Split(strings.Replace(strings.Trim(line, " "), "\t", "", -1), " ")
			if len(lineDetails) > 2 {
				dhcp.OptionNTPServer = lineDetails[2]
			}
		}
	}

	lines = strings.Split(string(data1), "\n")
	for i := range lines {
		line := lines[i]
		if len(line) == 0 || string(line[0]) == "#" {
			continue
		}
		if strings.Contains(line, "INTERFACES=") {
			parsed := strings.Replace(strings.Replace(line, "INTERFACES=\"", "", -1), "\"", "", -1)
			if parsed == "" {
				dhcp.Enabled = false
			} else {
				dhcp.Enabled = true
			}
		}
	}

	return
}

// SaveLinuxDHCPServer - Save Linux DHCP server config
func SaveLinuxDHCPServer(dhcp DHCPDConf) (err error) {

	file := "ddns-update-style none;\n"
	path := "/etc/dhcp/dhcpd.conf"

	if extensions.DoesFileExist(path) == false {
		return
	}

	file += "default-lease-time " + dhcp.DefaultLeaseTime + ";\n"
	file += "max-lease-time " + dhcp.MaxLeaseTime + ";\n"
	file += "authoritative;\n\n"

	file += "subnet " + dhcp.Subnet + " netmask " + dhcp.Netmask + " {\n"
	file += "    range " + dhcp.RangeFrom + " " + dhcp.RangeTo + ";\n"
	if dhcp.OptionRouter != "" {
		file += "    option routers " + dhcp.OptionRouter + ";\n"
	}
	if dhcp.OptionSubnetMask != "" {
		file += "    option subnet-mask " + dhcp.OptionSubnetMask + ";\n"
	}
	if dhcp.OptionDNS1 != "" {
		file += "    option domain-name-servers " + dhcp.OptionDNS1
	}
	if dhcp.OptionDNS2 != "" {
		file += ", " + dhcp.OptionDNS2
	}
	if dhcp.OptionDNS1 != "" {
		file += ";\n"
	}
	if dhcp.OptionNTPServer != "" {
		file += "    option ntp-servers " + dhcp.OptionNTPServer + ";\n"
	}
	file += "}"

	info, _ := os.Stat(path)
	mode := info.Mode()

	err = extensions.WriteToFile(file, path, mode)
	if err != nil {
		return
	}

	file = ""
	path = "/etc/default/isc-dhcp-server"

	if extensions.DoesFileExist(path) == false {
		return
	}

	data, err := extensions.ReadFile(path)
	if err != nil {
		session_functions.Log("Error", "Failed to read /etc/default/isc-dhcp-server:  "+err.Error())
		return
	}

	lines := strings.Split(string(data), "\n")
	for i := range lines {
		line := lines[i]
		if len(line) > 0 && string(line[0]) == "#" {
			file += line + "\n"
		}
	}

	file += "INTERFACES=\""
	if dhcp.Enabled == true {
		file += "eth1"
	}

	file += "\"\n"

	info, _ = os.Stat(path)
	mode = info.Mode()

	err = extensions.WriteToFile(file, path, mode)
	if err != nil {
		return
	}

	//Run the commands to either Start or Stop the DHCP Service
	if runtime.GOOS == "linux" {
		if dhcp.Enabled {
			err = exec.Command("/usr/bin/sudo", "/bin/systemctl", "enable", "isc-dhcp-server").Run()
			if err != nil {
				session_functions.Log("Error->networks->networks.go->SaveLinuxDHCPServer-> Enabling isc-dhcp-server:  ", err.Error())
				return
			}
			err = exec.Command("/usr/bin/sudo", "/bin/systemctl", "start", "isc-dhcp-server").Run()
			if err != nil {
				session_functions.Log("Error->networks->networks.go->SaveLinuxDHCPServer-> Enabling isc-dhcp-server:  ", err.Error())
				return
			}
		} else {
			err = exec.Command("/usr/bin/sudo", "/bin/systemctl", "disable", "isc-dhcp-server").Run()
			if err != nil {
				session_functions.Log("Error->networks->networks.go->SaveLinuxDHCPServer-> Enabling isc-dhcp-server:  ", err.Error())
				return
			}
			err = exec.Command("/usr/bin/sudo", "/bin/systemctl", "stop", "isc-dhcp-server").Run()
			if err != nil {
				session_functions.Log("Error->networks->networks.go->SaveLinuxDHCPServer-> Enabling isc-dhcp-server:  ", err.Error())
				return
			}
		}
	}

	return
}

// ValidateDHCPServer - Validate DHCP server config
func ValidateDHCPServer(dhcp *DHCPDConf) (err error) {
	if dhcp.Enabled {
		if dhcp.Subnet == "" {
			err = errors.New("Validation Error")
			dhcp.Errors.Subnet = "required"
			return
		}
		if dhcp.Netmask == "" {
			err = errors.New("Validation Error")
			dhcp.Errors.Netmask = "required"
			return
		}
		if dhcp.RangeTo == "" {
			err = errors.New("Validation Error")
			dhcp.Errors.RangeTo = "required"
			return
		}
		if dhcp.RangeFrom == "" {
			err = errors.New("Validation Error")
			dhcp.Errors.RangeFrom = "required"
			return
		}
		if dhcp.Subnet != "" && net.ParseIP(dhcp.Subnet) == nil {
			err = errors.New("Validation Error")
			dhcp.Errors.Subnet = "invalidValue"
			return
		}
		if dhcp.Netmask != "" && net.ParseIP(dhcp.Netmask) == nil {
			err = errors.New("Validation Error")
			dhcp.Errors.Netmask = "invalidValue"
			return
		}
		if dhcp.RangeTo != "" && net.ParseIP(dhcp.RangeTo) == nil {
			err = errors.New("Validation Error")
			dhcp.Errors.RangeTo = "invalidValue"
			return
		}
		if dhcp.RangeFrom != "" && net.ParseIP(dhcp.RangeFrom) == nil {
			err = errors.New("Validation Error")
			dhcp.Errors.RangeFrom = "invalidValue"
			return
		}
		if dhcp.OptionRouter != "" && net.ParseIP(dhcp.OptionRouter) == nil {
			err = errors.New("Validation Error")
			dhcp.Errors.OptionRouter = "invalidValue"
			return
		}
		if dhcp.OptionSubnetMask != "" && net.ParseIP(dhcp.OptionSubnetMask) == nil {
			err = errors.New("Validation Error")
			dhcp.Errors.OptionSubnetMask = "invalidValue"
			return
		}
		if dhcp.OptionDNS1 != "" && net.ParseIP(dhcp.OptionDNS1) == nil {
			err = errors.New("Validation Error")
			dhcp.Errors.OptionDNS1 = "invalidValue"
			return
		}
		if dhcp.OptionDNS2 != "" && net.ParseIP(dhcp.OptionDNS2) == nil {
			err = errors.New("Validation Error")
			dhcp.Errors.OptionDNS2 = "invalidValue"
			return
		}
		if dhcp.OptionNTPServer != "" && net.ParseIP(dhcp.OptionNTPServer) == nil {
			err = errors.New("Validation Error")
			dhcp.Errors.OptionNTPServer = "invalidValue"
			return
		}
	}
	return
}

// GetLinuxHosts - Returns linux hosts
func GetLinuxHosts() (hosts []Host, err error) {
	if extensions.DoesFileExist("/etc/hosts") == false {
		return
	}

	data, err := extensions.ReadFile("/etc/hosts")
	if err != nil {
		session_functions.Log("Error", "Failed to read /etc/hosts:  "+err.Error())
		return
	}

	lines := strings.Split(string(data), "\n")
	for i := range lines {
		line := lines[i]
		if len(line) == 0 || string(line[0]) == "#" {
			continue
		}
		hostItems := strings.Split(strings.TrimRight(line, "\t"), "\t")

		if len(hostItems) > 1 && len(hostItems[len(hostItems)-1]) > 0 {
			var h Host
			h.Address = hostItems[0]
			h.Domain = hostItems[len(hostItems)-1]
			hosts = append(hosts, h)
		}

		hostItems = strings.Split(strings.TrimRight(line, " "), " ")

		if len(hostItems) > 1 && len(hostItems[len(hostItems)-1]) > 0 {
			var h Host
			h.Address = hostItems[0]
			h.Domain = hostItems[len(hostItems)-1]
			hosts = append(hosts, h)
		}

	}

	return
}

// SaveLinuxHosts - Saves and validates Linux hosts. Returns hosts with errors
func SaveLinuxHosts(hosts []Host) (hostsReturned []Host, err error) {

	file := ""
	path := "/etc/hosts"

	data, err := extensions.ReadFile("/etc/hosts")
	if err != nil {
		session_functions.Log("Error", "Failed to read /etc/hosts:  "+err.Error())
		return
	}
	lines := strings.Split(string(data), "\n")
	for i := range lines {
		line := lines[i]
		if len(line) > 0 && string(line[0]) == "#" {
			file += line + "\n"
		}
	}

	for i := range hosts {
		host := &hosts[i]
		err = ValidateLinuxHost(host)
		if err != nil {
			hostsReturned = hosts
			return
		}
		file += host.Address + "\t" + host.Domain + "\n"
		if strings.Contains(host.Address, ":") == false { //Don't add IPV6 Records  This is for LAN only networks.  For now.
			if AddDNSRecord != nil {
				AddDNSRecord(host.Domain, host.Address)
			}
		}
	}

	info, _ := os.Stat(path)
	mode := info.Mode()

	err = extensions.WriteToFile(file, path, mode)
	if err != nil {
		return
	}
	return
}

// ValidateLinuxHosts - Validate linux hosts
func ValidateLinuxHosts(hosts []Host) (hostsReturned []Host, err error) {
	for i := range hosts {
		host := &hosts[i]
		err = ValidateLinuxHost(host)
		if err != nil {
			hostsReturned = hosts
			return
		}
	}
	return
}

// ValidateLinuxHost - Validate Linux host
func ValidateLinuxHost(host *Host) (err error) {
	if host.Address == "" {
		err = errors.New("Validation Error")
		host.Errors.Address = "required"
		return
	}
	if net.ParseIP(host.Address) == nil {
		err = errors.New("Validation Error")
		host.Errors.Address = "invalidValue"
		return
	}
	return
}

// AddHost - Add Host to Linux network
func AddHost(address string, domain string) (err error) {
	hosts, err := GetLinuxHosts()
	if err != nil {
		return
	}
	//Find if the host exists
	doesHostExist := false
	for i := range hosts {
		host := hosts[i]
		if host.Address == address {
			doesHostExist = true
			host.Domain = domain
		}
	}

	if doesHostExist == false {
		var h Host
		h.Address = address
		h.Domain = domain
		hosts = append(hosts, h)
	}

	_, err = SaveLinuxHosts(hosts)
	return
}

// GetLinuxNetworks - Returns linux network files
func GetLinuxNetworks() (networks []LinuxNetworkFile, err error) {
	if extensions.DoesFileExist(NETWORK_PATH+"/interfaces.d/eth0") == false {
		extensions.WriteToFile("auto eth0\niface eth0 inet static\naddress 192.168.1.254\nnetmask 255.255.255.0\ngateway 192.168.1.1", NETWORK_PATH+"/interfaces.d/eth0", 777)
	}

	if extensions.DoesFileExist(NETWORK_PATH+"/interfaces.d/eth1") == false {
		extensions.WriteToFile("auto eth1\niface eth1 inet dhcp", NETWORK_PATH+"/interfaces.d/eth1", 777)
	}

	data0, err := extensions.ReadFile(NETWORK_PATH + "/interfaces.d/eth0")
	if err != nil {
		session_functions.Log("Error", "Failed to read /etc/network/interfaces.d/eth0:  "+err.Error())
		return
	}

	data1, err := extensions.ReadFile(NETWORK_PATH + "/interfaces.d/eth1")
	if err != nil {
		session_functions.Log("Error", "Failed to read /etc/network/interfaces.d/eth1:  "+err.Error())
		return
	}

	network0, err := ParseNetworkFileWPA(data0)
	if err != nil {
		session_functions.Log("Error", "Failed to parse network 0:  "+err.Error())
		return
	}

	network1, err := ParseNetworkFileWPA(data1)
	if err != nil {
		session_functions.Log("Error", "Failed to parse network 1:  "+err.Error())
		return
	}

	networks = append(networks, network0)
	networks = append(networks, network1)
	return
}

// SaveLinuxNetworks - Saves Linux Network files
func SaveLinuxNetworks(networks []LinuxNetworkFile) (networksReturned []LinuxNetworkFile, err error) {

	for i := range networks {
		network := &networks[i]
		err = ValidateNetwork(network)
		if err != nil {
			networksReturned = networks
			return
		}

		path := NETWORK_PATH + "/interfaces.d/" + network.Name

		file := "auto " + network.Name + "\n"
		file += "iface " + network.Name + " inet "
		if network.DHCP {
			file += "dhcp\n"
			if network.HWAddress != "" {
				file += "hwaddress ether " + network.HWAddress + "\n"
			}
		} else {
			file += "static\n"
			file += "address " + network.Address + "\n"
			file += "netmask " + network.Subnet + "\n"
			if network.Gateway != "" {
				file += "gateway " + network.Gateway + "\n"
			}
			if network.HWAddress != "" {
				file += "hwaddress ether " + network.HWAddress + "\n"
			}
			if network.DNSNameserver1 != "" || network.DNSNameserver2 != "" {
				file += "dns-nameservers " + network.DNSNameserver1 + " " + network.DNSNameserver2 + "\n"
			}
		}
		file += "\n\n\n\n\n"
		if len(network.StaticRoutes) > 0 {
			file += "#### Static Routes #####\n"
			for j := range network.StaticRoutes {
				route := network.StaticRoutes[j]
				file += "up route add -net " + route.Address + " netmask " + route.Netmask + " gw " + route.Gateway + " dev " + network.Name + "\n"
			}
		}

		err = extensions.WriteToFile(file, path, 777)
		if err != nil {
			return
		}
	}

	return
}

// ValidateNetworks - Returns all networks with invalid configurations
func ValidateNetworks(networks []LinuxNetworkFile) (networksReturned []LinuxNetworkFile, err error) {
	for i := range networks {
		network := &networks[i]
		err = ValidateNetwork(network)
		if err != nil {
			networksReturned = networks
			return
		}
	}
	return
}

// ValidateNetwork - Validates if the network is a valid configuration
func ValidateNetwork(network *LinuxNetworkFile) (err error) {
	if network.DHCP == false {
		if network.Address == "" {
			err = errors.New("Validation Error")
			network.Errors.Address = "required"
			return
		}
		if network.Subnet == "" {
			err = errors.New("Validation Error")
			network.Errors.Subnet = "required"
			return
		}
		if network.Address != "" && net.ParseIP(network.Address) == nil {
			err = errors.New("Validation Error")
			network.Errors.Address = "invalidValue"
			return
		}
		if network.Subnet != "" && net.ParseIP(network.Subnet) == nil {
			err = errors.New("Validation Error")
			network.Errors.Subnet = "invalidValue"
			return
		}
		if network.Gateway != "" && net.ParseIP(network.Gateway) == nil {
			err = errors.New("Validation Error")
			network.Errors.Gateway = "invalidValue"
			return
		}
		if network.DNSNameserver1 != "" && net.ParseIP(network.DNSNameserver1) == nil {
			err = errors.New("Validation Error")
			network.Errors.DNSNameserver1 = "invalidValue"
			return
		}
		if network.DNSNameserver2 != "" && net.ParseIP(network.DNSNameserver2) == nil {
			err = errors.New("Validation Error")
			network.Errors.DNSNameserver2 = "invalidValue"
			return
		}
	}

	for i := range network.StaticRoutes {
		route := &network.StaticRoutes[i]
		if route.Address == "" {
			err = errors.New("Validation Error")
			route.Errors.Address = "required"
			return
		}
		if route.Netmask == "" {
			err = errors.New("Validation Error")
			route.Errors.Netmask = "required"
			return
		}
		if route.Gateway == "" {
			err = errors.New("Validation Error")
			route.Errors.Gateway = "required"
			return
		}
		if net.ParseIP(route.Address) == nil {
			err = errors.New("Validation Error")
			route.Errors.Address = "invalidValue"
			return
		}
		if net.ParseIP(route.Netmask) == nil {
			err = errors.New("Validation Error")
			route.Errors.Netmask = "invalidValue"
			return
		}
		if net.ParseIP(route.Gateway) == nil {
			err = errors.New("Validation Error")
			route.Errors.Gateway = "invalidValue"
			return
		}
	}
	return
}

// ParseNetworkFileWPA - Parse the network file WPA? not sure what wpa is
func ParseNetworkFileWPA(data []byte) (network LinuxNetworkFile, err error) {
	ascii := string(data)

	lines := strings.Split(string(data), "\n")
	for i := range lines {
		line := lines[i]
		if strings.Contains(line, "auto") && !strings.Contains(line, "up route") {
			network.Name = strings.Replace(line, "auto ", "", -1)
		} else if strings.Contains(line, "iface") && !strings.Contains(line, "up route") {
			if strings.Contains(ascii, "dhcp") {
				network.DHCP = true
				liveNetworks := GetNetworks()
				for j := range liveNetworks {
					liveNetwork := liveNetworks[j]
					if liveNetwork.Name == network.Name {
						network.Address = liveNetwork.Address
						break
					}
				}
			}
		} else if strings.Contains(line, "address") && !strings.Contains(line, "up route") && !strings.Contains(line, "hwaddress") {
			network.Address = strings.Replace(line, "address ", "", -1)
		} else if strings.Contains(line, "hwaddress") && !strings.Contains(line, "up route") {
			network.HWAddress = strings.Replace(line, "hwaddress ether ", "", -1)
		} else if strings.Contains(line, "netmask") && !strings.Contains(line, "up route") {
			network.Subnet = strings.Replace(line, "netmask ", "", -1)
		} else if strings.Contains(line, "gateway") && !strings.Contains(line, "up route") {
			network.Gateway = strings.Replace(line, "gateway ", "", -1)
		} else if strings.Contains(line, "dns-nameservers") && !strings.Contains(line, "up route") {
			nameserversString := strings.Replace(line, "dns-nameservers ", "", -1)
			nameservers := strings.Split(nameserversString, " ")
			if len(nameservers) >= 1 {
				network.DNSNameserver1 = nameservers[0]
			}
			if len(nameservers) >= 2 {
				network.DNSNameserver2 = nameservers[1]
			}
		} else if strings.Contains(line, "up route") {

			var sr StaticRoute

			words := strings.Split(line, " ")
			sr.Address = words[4]
			sr.Netmask = words[6]
			sr.Gateway = words[8]
			sr.Interface = words[10]
			network.StaticRoutes = append(network.StaticRoutes, sr)

		}
	}

	queriedNetworks := GetNetworks()

	if network.Address == "" {
		for j := range queriedNetworks {
			queriedNetwork := queriedNetworks[j]
			if queriedNetwork.Name == network.Name {
				network.Address = queriedNetwork.Address
				break
			}
		}
	}
	return
}

// ParseNetworkFile - Parse the network file for all of its networking details
func ParseNetworkFile(data []byte) (network LinuxNetworkFile, err error) {
	ascii := string(data)

	lines := strings.Split(string(data), "\n")
	for i := range lines {
		line := lines[i]
		if strings.Contains(line, "auto") && !strings.Contains(line, "up route") {
			network.Name = strings.Replace(line, "auto ", "", -1)
		} else if strings.Contains(line, "iface") && !strings.Contains(line, "up route") {
			if strings.Contains(ascii, "dhcp") {
				network.DHCP = true
				liveNetworks := GetNetworks()
				for j := range liveNetworks {
					liveNetwork := liveNetworks[j]
					if liveNetwork.Name == network.Name {
						network.Address = liveNetwork.Address
					}
				}
			}
		} else if strings.Contains(line, "address") && !strings.Contains(line, "up route") && !strings.Contains(line, "hwaddress") {
			network.Address = strings.Replace(line, "address ", "", -1)
		} else if strings.Contains(line, "hwaddress") && !strings.Contains(line, "up route") {
			network.HWAddress = strings.Replace(line, "hwaddress ether ", "", -1)
		} else if strings.Contains(line, "netmask") && !strings.Contains(line, "up route") {
			network.Subnet = strings.Replace(line, "netmask ", "", -1)
		} else if strings.Contains(line, "gateway") && !strings.Contains(line, "up route") {
			network.Gateway = strings.Replace(line, "gateway ", "", -1)
		} else if strings.Contains(line, "dns-nameservers") && !strings.Contains(line, "up route") {
			nameserversString := strings.Replace(line, "dns-nameservers ", "", -1)
			nameservers := strings.Split(nameserversString, " ")
			if len(nameservers) >= 1 {
				network.DNSNameserver1 = nameservers[0]
			}
			if len(nameservers) >= 2 {
				network.DNSNameserver2 = nameservers[1]
			}
		} else if strings.Contains(line, "up route") {

			var sr StaticRoute

			words := strings.Split(line, " ")
			sr.Address = words[4]
			sr.Netmask = words[6]
			sr.Gateway = words[8]
			sr.Interface = words[10]
			network.StaticRoutes = append(network.StaticRoutes, sr)

		}
	}

	queriedNetworks := GetNetworks()
	for j := range queriedNetworks {
		queriedNetwork := queriedNetworks[j]
		if queriedNetwork.Name == network.Name {
			network.Address = queriedNetwork.Address
		}
	}

	return
}

// CheckDHCPNetworks - Checks DHCP Network if it's still up
func CheckDHCPNetworks() {

	CurrentLinuxNICDHCPStatus.Lock()
	for key, value := range CurrentLinuxNICDHCPStatus.CurrentStatus {
		data, err := extensions.ReadFile("/sys/class/net/" + key + "/operstate")
		if err == nil {
			upDown := strings.TrimRight(string(data), "\n")
			if value == false && strings.Contains(upDown, "up") { //Issue dhclient since we have a change from down to up for dhcp.
				session_functions.Log("DHCP Renewal", "Link changed state from down to up on "+key+".  Releasing and Renewing Lease:  ")
				err = exec.Command("/usr/bin/sudo", "/sbin/dhclient", "-r", key).Run()
				if err != nil {
					session_functions.Log("Error->networks->CheckDHCPNetworks", "Failed to Release DHCP on NIC "+key+":  "+err.Error())
				}
				errRenew := exec.Command("/usr/bin/sudo", "/sbin/dhclient", key).Run()
				if errRenew != nil {
					session_functions.Log("Error->networks->CheckDHCPNetworks", "Failed to Renew Lease DHCP on NIC "+key+":  "+errRenew.Error())
				}
				if err == nil && errRenew == nil {
					session_functions.Log("DHCP Renewal", "Successfully Released and Renewed DHCP Lease "+key)
				}
			}
			if value == true && strings.Contains(upDown, "down") {
				session_functions.Log("DHCP Status", "DHCP Network Down on NIC "+key)
			}
			if strings.Contains(upDown, "up") {
				CurrentLinuxNICDHCPStatus.CurrentStatus[key] = true
			} else {
				CurrentLinuxNICDHCPStatus.CurrentStatus[key] = false
			}
		}
	}
	CurrentLinuxNICDHCPStatus.Unlock()
}

// GetCIDRFromSubnet - Returns CIDR by subnet
func GetCIDRFromSubnet(subnet string) (cidr int) {
	switch subnet {
	case "255.255.255.252":
		cidr = 30
	case "255.255.255.248":
		cidr = 29
	case "255.255.255.240":
		cidr = 28
	case "255.255.255.224":
		cidr = 27
	case "255.255.255.192":
		cidr = 26
	case "255.255.255.128":
		cidr = 25
	case "255.255.255.0":
		cidr = 24
	case "255.255.254.0":
		cidr = 23
	case "255.255.252.0":
		cidr = 22
	case "255.255.248.0":
		cidr = 21
	case "255.255.240.0":
		cidr = 20
	case "255.255.224.0":
		cidr = 19
	case "255.255.192.0":
		cidr = 18
	case "255.255.128.0":
		cidr = 27
	case "255.255.0.0":
		cidr = 16
	}
	return
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func IsWithinRange(ip string, iplow string, iphigh string) bool {
	ip1 := net.ParseIP(iplow)
	if ip1.To4() == nil {
		return false
	}
	ip2 := net.ParseIP(iphigh)
	if ip2.To4() == nil {
		return false
	}
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false
	}
	if bytes.Compare(trial, ip1) >= 0 && bytes.Compare(trial, ip2) <= 0 {
		return true
	}
	return false
}

func IsGreaterThan(ip string, iplow string) bool {
	ip1 := net.ParseIP(iplow)
	if ip1.To4() == nil {
		return false
	}
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false
	}
	if bytes.Compare(trial, ip1) > 0 {
		return true
	}
	return false
}

func IsEqualTo(ip string, iplow string) bool {
	ip1 := net.ParseIP(iplow)
	if ip1.To4() == nil {
		return false
	}
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false
	}
	if bytes.Compare(trial, ip1) == 0 {
		return true
	}
	return false
}

func IsGreaterThanOrEqualTo(ip string, iplow string) bool {
	ip1 := net.ParseIP(iplow)
	if ip1.To4() == nil {
		return false
	}
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false
	}
	if bytes.Compare(trial, ip1) >= 0 {
		return true
	}
	return false
}

func IsLessThan(ip string, iplow string) bool {
	ip1 := net.ParseIP(iplow)
	if ip1.To4() == nil {
		return false
	}
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false
	}
	if bytes.Compare(trial, ip1) < 0 {
		return true
	}
	return false
}

func IsLessThanOrEqualTo(ip string, iplow string) bool {
	ip1 := net.ParseIP(iplow)
	if ip1.To4() == nil {
		return false
	}
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false
	}
	if bytes.Compare(trial, ip1) <= 0 {
		return true
	}
	return false
}
