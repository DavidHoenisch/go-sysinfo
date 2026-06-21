package gosysinfo

import (
	"fmt"
	"os"
)

const (
	net_base string = "/sys/class/net"
)

type Networks struct {
	NInterfaces []string
}

type NetworkConnectionInfo struct {
	IP string
}

type NetworkStatistics struct {
	Addr_assign_type     string
	Address              string
	Addr_len             string
	Broadcast            string
	Carrier              string
	Carrier_changes      string
	Carrier_down_count   string
	Carrier_up_count     string
	Device               string
	Dev_id               string
	Dev_port             string
	Dormant              string
	Duplex               string
	Flags                string
	Gro_flush_timeout    string
	Ifalias              string
	Ifindex              string
	Iflink               string
	Link_mode            string
	Mtu                  string
	Name_assign_type     string
	Napi_defer_hard_irqs string
	Netdev_group         string
	Operstate            string
	Phy80211             string
	Power                string
	Proto_down           string
	Queues               string
	Speed                string
	Statistics           string
	Subsystem            string
	Testing              string
	Threaded             string
	Tx_queue_len         string
	Type                 string
	Uevent               string
	Wireless             string
}

type NetworkInfo interface {
	GetNetworkStatistics(ifname string) (*NetworkStatistics, error)
	GetNetowrkConnectionInfo(ifname string) (NetworkInfo, error)
}

func GetNetowrkConnectionInfo(ifname string) (*NetworkConnectionInfo, error) {

	path, err := checkNetworkPathByName(ifname)

	if err != nil {
		return nil, err
	}

	return &NetworkConnectionInfo{
		// TODO: this does not make a distinction between IPv6
		// IPv4
		IP: getFileContent(fmt.Sprintf("%s/%s", path, "address")),
	}, nil

}

func GetNetworkStatistics(ifname string) (*NetworkStatistics, error) {
	path, err := checkNetworkPathByName(ifname)

	if err != nil {
		return nil, err
	}

	return &NetworkStatistics{
		Addr_assign_type:     getFileContent(path),
		Address:              getFileContent(path),
		Addr_len:             getFileContent(path),
		Broadcast:            getFileContent(path),
		Carrier:              getFileContent(path),
		Carrier_changes:      getFileContent(path),
		Carrier_down_count:   getFileContent(path),
		Carrier_up_count:     getFileContent(path),
		Device:               getFileContent(path),
		Dev_id:               getFileContent(path),
		Dev_port:             getFileContent(path),
		Dormant:              getFileContent(path),
		Duplex:               getFileContent(path),
		Flags:                getFileContent(path),
		Gro_flush_timeout:    getFileContent(path),
		Ifalias:              getFileContent(path),
		Ifindex:              getFileContent(path),
		Iflink:               getFileContent(path),
		Link_mode:            getFileContent(path),
		Mtu:                  getFileContent(path),
		Name_assign_type:     getFileContent(path),
		Napi_defer_hard_irqs: getFileContent(path),
		Netdev_group:         getFileContent(path),
		Operstate:            getFileContent(path),
		Phy80211:             getFileContent(path),
		Power:                getFileContent(path),
		Proto_down:           getFileContent(path),
		Queues:               getFileContent(path),
		Speed:                getFileContent(path),
		Statistics:           getFileContent(path),
		Subsystem:            getFileContent(path),
		Testing:              getFileContent(path),
		Threaded:             getFileContent(path),
		Tx_queue_len:         getFileContent(path),
		Type:                 getFileContent(path),
		Uevent:               getFileContent(path),
		Wireless:             getFileContent(path),
	}, nil
}

// getNetworkPathByName() checks that the requested network
// exists. If not, returns ErrNetworkNotFound error.
func checkNetworkPathByName(ifname string) (string, error) {

	path := fmt.Sprintf("%s/%s", net_base, ifname)
	_, err := os.Stat(path)

	if err != nil {
		// TODO: add in custom error
		return "", ErrNetworkNotFound
	}

	return path, nil

}
