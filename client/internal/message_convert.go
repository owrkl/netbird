package internal

import (
	"errors"
	"fmt"
	"net"

	firewallManager "github.com/netbirdio/netbird/client/firewall/manager"
	mgmProto "github.com/netbirdio/netbird/management/proto"
)

func convertToFirewallProtocol(protocol mgmProto.RuleProtocol) (firewallManager.Protocol, error) {
	switch protocol {
	case mgmProto.RuleProtocol_TCP:
		return firewallManager.ProtocolTCP, nil
	case mgmProto.RuleProtocol_UDP:
		return firewallManager.ProtocolUDP, nil
	case mgmProto.RuleProtocol_ICMP:
		return firewallManager.ProtocolICMP, nil
	case mgmProto.RuleProtocol_ALL:
		return firewallManager.ProtocolALL, nil
	default:
		return firewallManager.ProtocolALL, fmt.Errorf("invalid protocol type: %s", protocol.String())
	}
}

// convertPortInfo todo: write validation for portInfo
func convertPortInfo(portInfo *mgmProto.PortInfo) *firewallManager.Port {
	if portInfo == nil {
		return nil
	}

	if portInfo.GetPort() != 0 {
		return &firewallManager.Port{
			Values: []int{int(portInfo.GetPort())},
		}
	}

	if portInfo.GetRange() != nil {
		return &firewallManager.Port{
			IsRange: true,
			Values:  []int{int(portInfo.GetRange().Start), int(portInfo.GetRange().End)},
		}
	}

	return nil
}

func convertToIP(rawIP []byte) (net.IP, error) {
	if rawIP == nil {
		return nil, errors.New("input bytes cannot be nil")
	}

	if len(rawIP) != net.IPv4len && len(rawIP) != net.IPv6len {
		return nil, fmt.Errorf("invalid IP length: %d", len(rawIP))
	}

	return rawIP, nil
}
