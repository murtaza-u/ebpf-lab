// Code generated by bpf2go; DO NOT EDIT.
//go:build mips || mips64 || ppc64 || s390x

package drop

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

// loadDrop returns the embedded CollectionSpec for drop.
func loadDrop() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_DropBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load drop: %w", err)
	}

	return spec, err
}

// loadDropObjects loads drop and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*dropObjects
//	*dropPrograms
//	*dropMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadDropObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadDrop()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// dropSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type dropSpecs struct {
	dropProgramSpecs
	dropMapSpecs
}

// dropSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type dropProgramSpecs struct {
	DropPacketByPort *ebpf.ProgramSpec `ebpf:"drop_packet_by_port"`
}

// dropMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type dropMapSpecs struct {
	PortMap *ebpf.MapSpec `ebpf:"port_map"`
}

// dropObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadDropObjects or ebpf.CollectionSpec.LoadAndAssign.
type dropObjects struct {
	dropPrograms
	dropMaps
}

func (o *dropObjects) Close() error {
	return _DropClose(
		&o.dropPrograms,
		&o.dropMaps,
	)
}

// dropMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadDropObjects or ebpf.CollectionSpec.LoadAndAssign.
type dropMaps struct {
	PortMap *ebpf.Map `ebpf:"port_map"`
}

func (m *dropMaps) Close() error {
	return _DropClose(
		m.PortMap,
	)
}

// dropPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadDropObjects or ebpf.CollectionSpec.LoadAndAssign.
type dropPrograms struct {
	DropPacketByPort *ebpf.Program `ebpf:"drop_packet_by_port"`
}

func (p *dropPrograms) Close() error {
	return _DropClose(
		p.DropPacketByPort,
	)
}

func _DropClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed drop_bpfeb.o
var _DropBytes []byte