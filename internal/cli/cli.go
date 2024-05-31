package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/murtaza-u/ebpf-lab/internal/drop"

	"github.com/urfave/cli/v2"
)

// DefaultDropPort is the default port on which the eBPF program will drop
// packets, if no port is specified.
const DefaultDropPort = 4040

// Cmd provides a command-line interface to run the eBPF program from
// userspace.
var Cmd = &cli.App{
	Name:                 "drop",
	Usage:                "Drop packets on a given port",
	UsageText:            "--interface [INTERFACE] [--port]",
	EnableBashCompletion: true,
	Copyright:            "GPL",
	Authors: []*cli.Author{
		{Name: "Murtaza Udaipurwala", Email: "murtaza@murtazau.xyz"},
	},
	Flags: []cli.Flag{
		&cli.UintFlag{
			Name:  "port",
			Value: DefaultDropPort,
		},
		&cli.StringFlag{
			Name:     "interface",
			Required: true,
			Usage:    "name of the network interface. Hint: ip link show",
		},
	},
	Action: func(ctx *cli.Context) error {
		cc, cancel := context.WithCancel(ctx.Context)
		defer cancel()

		go handleInterrupt(cc, cancel)

		port := ctx.Uint("port")
		if port > 65535 {
			return fmt.Errorf("invalid port %d", port)
		}
		iface := ctx.String("interface")

		log.Printf("attempting to drop traffic on %q on port %d", iface, port)
		return drop.Run(cc, iface, uint8(port))
	},
}

func handleInterrupt(ctx context.Context, cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		return
	case <-sig:
		cancel()
	}
}
