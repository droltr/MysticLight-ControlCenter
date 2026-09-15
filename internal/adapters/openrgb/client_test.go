package openrgb

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"
)

func TestClientObservesControllerCountWithoutWritePackets(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	serverDone := make(chan error, 1)
	go func() {
		connection, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverDone <- acceptErr
			return
		}
		defer connection.Close()
		for index, response := range [][]byte{{6, 0, 0, 0}, {3, 0, 0, 0}} {
			header := make([]byte, packetHeaderSize)
			if _, err := io.ReadFull(connection, header); err != nil {
				serverDone <- err
				return
			}
			if string(header[:4]) != string(magic[:]) {
				serverDone <- ErrBadResponse
				return
			}
			packetID := binary.LittleEndian.Uint32(header[8:12])
			if index == 0 && packetID != protocolVersionID || index == 1 && packetID != controllerCountID {
				serverDone <- ErrBadResponse
				return
			}
			responseHeader := make([]byte, packetHeaderSize)
			copy(responseHeader[:4], magic[:])
			binary.LittleEndian.PutUint32(responseHeader[8:12], packetID)
			binary.LittleEndian.PutUint32(responseHeader[12:16], uint32(len(response)))
			if _, err := connection.Write(responseHeader); err != nil {
				serverDone <- err
				return
			}
			if _, err := connection.Write(response); err != nil {
				serverDone <- err
				return
			}
		}
		serverDone <- nil
	}()

	observation, err := New(listener.Addr().String(), time.Second).Observe(context.Background())
	if err != nil {
		t.Fatalf("observe: %v", err)
	}
	result := observation.(Observation)
	if result.ProtocolVersion != 6 || result.ControllerCount != 3 {
		t.Fatalf("unexpected observation: %#v", result)
	}
	if err := <-serverDone; err != nil {
		t.Fatalf("mock server: %v", err)
	}
}
