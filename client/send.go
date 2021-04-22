package client

import "github.com/M4THYOU/some_mqtt_broker/mqtt"

func (client *Client) SendPacket(packetCode uint8) error {
	packet, err := client.BuildPacket(mqtt.ConnackCode)
	if err != nil {
		return err
	}
	return mqtt.SendPacket(client.Conn, packet)
}

func (client *Client) BuildPacket(packetCode uint8) ([]byte, error) {
	packet := []byte{}
	var reqType byte
	if packetCode == mqtt.PublishCode {
		// TODO HANDLE PUBLISH FLAGS!!
		reqType = mqtt.SetRequestType(packetCode, false, false, 0)
	} else {
		reqType = mqtt.SetRequestType(packetCode, false, false, 0)
	}
	packet = append(packet, reqType)
	return packet, nil
}
