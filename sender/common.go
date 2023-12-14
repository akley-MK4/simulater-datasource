package sender

import (
	"errors"
	"fmt"
	"github.com/akley-MK4/go-tools-box/netaic"
	"io"
	"net"
)

func FindInterfaceInfoByName(name string) (retInfo net.Interface, retErr error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		retErr = err
		return
	}

	for _, info := range interfaces {
		if info.Name == name {
			retInfo = info
			return
		}
	}

	retErr = fmt.Errorf("interface with name %s not exists", name)
	return
}

func CreateIOWriter(network string, srcInterfaceName string, dstIp string, dstPort uint16) (retIoWriter io.Writer, retErr error) {
	dstNIP := net.ParseIP(dstIp)
	if dstNIP == nil {
		return nil, errors.New("invalid dst ip")
	}

	srcIP, getSrcIpErr := GetInterfaceIp(srcInterfaceName)
	if getSrcIpErr != nil {
		return nil, getSrcIpErr
	}

	var lAddr *net.UDPAddr
	if srcIP != "" {
		lAddr = &net.UDPAddr{
			IP: net.ParseIP(srcIP),
		}
	}

	if lAddr != nil {
		retIoWriter, retErr = net.DialUDP(network, lAddr, &net.UDPAddr{
			IP: dstNIP,
		})
		if retErr != nil {
			retErr = fmt.Errorf("DialUDP Failed, LAddr: %v, Err: %v", srcIP, retErr)
		}
		return
	}

	retIoWriter, retErr = net.Dial(network, fmt.Sprintf("%v:%d", dstIp, dstPort))
	return
}

func GetInterfaceIp(interfaceName string) (retIp string, retErr error) {
	if interfaceName == "" {
		return
	}

	return netaic.GetNetInterfaceAddr(interfaceName)
}
