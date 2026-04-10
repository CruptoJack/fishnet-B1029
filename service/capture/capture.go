package capture

import (
	"context"
	"fmt"

	"github.com/google/gopacket/pcap"
)

func StartCapture(ctx context.Context, device string, outChanel chan<- []byte) error {
	handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("ошибка открытия устройства %s: %w", device, err)
	}
	defer handle.Close()

	for {
		data, capInfo, err := handle.ReadPacketData()
		if err != nil {
			return fmt.Errorf("ошибка чтения пакета: %w", err)
		}

		packetCopy := append([]byte(nil), data...)
		select {
		case outChanel <- packetCopy:
			fmt.Printf("Информация о пакете %v", capInfo)
		case <-ctx.Done():
			fmt.Println("Cancelled:", ctx.Err())
			return ctx.Err()
		}
	}
}
