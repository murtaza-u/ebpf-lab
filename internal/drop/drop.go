package drop

import (
	"context"
	"fmt"
	"net"

	"github.com/cilium/ebpf/link"
)

// Run attaches the eBPF program to the XDP probe on the specified network
// interface. It then communicates with the eBPF program through the eBPF map
// to notify it about the port on which to drop packets.
func Run(ctx context.Context, ifaceName string, port uint8) error {
	var objs dropObjects
	if err := loadDropObjects(&objs, nil); err != nil {
		return fmt.Errorf("loading objects: %w", err)
	}
	defer objs.Close()

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return fmt.Errorf("failed to get interface %q: %w", ifaceName, err)
	}

	lnk, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.DropPacketByPort,
		Interface: iface.Index,
	})
	if err != nil {
		return fmt.Errorf("attaching to XDP: %w", err)
	}
	defer lnk.Close()

	if err := objs.PortMap.Put(uint32(0), uint32(port)); err != nil {
		return fmt.Errorf("putting value into eBPF map: %w", err)
	}

	<-ctx.Done()

	return nil
}
