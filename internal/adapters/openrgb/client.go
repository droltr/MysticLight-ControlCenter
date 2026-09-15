package openrgb

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
	"github.com/droltr/MysticLight-ControlCenter/internal/provider"
)

const (
	defaultAddress     = "127.0.0.1:6742"
	protocolVersionID  = uint32(40)
	controllerCountID  = uint32(0)
	packetHeaderSize   = 16
	maxProtocolVersion = uint32(6)
)

var (
	magic          = [4]byte{'O', 'R', 'G', 'B'}
	ErrBadResponse = errors.New("invalid OpenRGB SDK response")
)

type Client struct {
	Address string
	Timeout time.Duration
}

type Observation struct {
	ProtocolVersion uint32
	ControllerCount uint32
}

var _ provider.Adapter = Client{}

func New(address string, timeout time.Duration) Client {
	if address == "" {
		address = defaultAddress
	}
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return Client{Address: address, Timeout: timeout}
}

func (c Client) Name() string {
	return "openrgb"
}

func (c Client) Capabilities() []domain.Capability {
	return []domain.Capability{{
		Resource:   domain.ResourceRGB,
		Operations: []domain.Operation{domain.OperationObserve},
	}}
}

func (c Client) Health(context.Context) domain.ProviderHealth {
	return domain.ProviderHealth{Provider: c.Name()}
}

func (c Client) Observe(ctx context.Context) (any, error) {
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	dialer := net.Dialer{Timeout: timeout}
	connection, err := dialer.DialContext(ctx, "tcp", c.Address)
	if err != nil {
		return nil, fmt.Errorf("connect to OpenRGB SDK: %w", err)
	}
	defer connection.Close()
	if err := connection.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, fmt.Errorf("set OpenRGB SDK deadline: %w", err)
	}

	if err := writePacket(connection, 0, protocolVersionID, nil); err != nil {
		return nil, err
	}
	protocolPayload, err := readPacket(connection, protocolVersionID)
	if err != nil {
		return nil, err
	}
	if len(protocolPayload) < 4 {
		return nil, ErrBadResponse
	}
	protocolVersion := binary.LittleEndian.Uint32(protocolPayload[:4])
	if protocolVersion > maxProtocolVersion {
		protocolVersion = maxProtocolVersion
	}

	if err := writePacket(connection, 0, controllerCountID, nil); err != nil {
		return nil, err
	}
	countPayload, err := readPacket(connection, controllerCountID)
	if err != nil {
		return nil, err
	}
	if len(countPayload) < 4 {
		return nil, ErrBadResponse
	}
	return Observation{
		ProtocolVersion: protocolVersion,
		ControllerCount: binary.LittleEndian.Uint32(countPayload[:4]),
	}, nil
}

func writePacket(writer io.Writer, deviceID, packetID uint32, payload []byte) error {
	header := make([]byte, packetHeaderSize)
	copy(header[:4], magic[:])
	binary.LittleEndian.PutUint32(header[4:8], deviceID)
	binary.LittleEndian.PutUint32(header[8:12], packetID)
	binary.LittleEndian.PutUint32(header[12:16], uint32(len(payload)))
	if _, err := writer.Write(header); err != nil {
		return fmt.Errorf("write OpenRGB SDK header: %w", err)
	}
	if len(payload) > 0 {
		if _, err := writer.Write(payload); err != nil {
			return fmt.Errorf("write OpenRGB SDK payload: %w", err)
		}
	}
	return nil
}

func readPacket(reader io.Reader, expectedID uint32) ([]byte, error) {
	header := make([]byte, packetHeaderSize)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, fmt.Errorf("read OpenRGB SDK header: %w", err)
	}
	if string(header[:4]) != string(magic[:]) || binary.LittleEndian.Uint32(header[8:12]) != expectedID {
		return nil, ErrBadResponse
	}
	size := binary.LittleEndian.Uint32(header[12:16])
	if size > 1024*1024 {
		return nil, ErrBadResponse
	}
	payload := make([]byte, size)
	if _, err := io.ReadFull(reader, payload); err != nil {
		return nil, fmt.Errorf("read OpenRGB SDK payload: %w", err)
	}
	return payload, nil
}
