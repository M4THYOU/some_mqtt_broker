package client

import (
	"errors"
	"fmt"

	"github.com/M4THYOU/some_mqtt_broker/pkg/mqtt"
)

func (client *Client) SendPacket(packetCode uint8) error {
	packet, err := client.buildPacket(mqtt.ConnackCode)
	if err != nil {
		return err
	}
	return mqtt.SendPacket(client.Conn, packet)
}

func (client *Client) buildPacket(packetCode uint8) ([]byte, error) {
	packet := []byte{}
	var reqType byte
	if packetCode == mqtt.PublishCode {
		// TODO HANDLE PUBLISH FLAGS!!
		reqType = mqtt.SetRequestType(packetCode, false, false, 0)
	} else {
		reqType = mqtt.SetRequestType(packetCode, false, false, 0)
	}
	packet = append(packet, reqType)

	// BUILD THE VARIABLE HEADER!
	varHeader, err := client.buildVarHeader(packetCode)
	if err != nil {
		return nil, err
	}
	fmt.Println(varHeader)

	return packet, nil
}

func (client *Client) buildVarHeader(packetCode uint8) (header []byte, err error) {
	switch packetCode {
	case mqtt.ConnectCode:
		// err = client.handleConnect()
	case mqtt.ConnackCode:
		// err = client.handleConnack()
	case mqtt.PublishCode:
		// err = client.handlePublish()
	case mqtt.PubackCode:
		// err = client.handlePuback()
	case mqtt.PubrecCode:
		// err = client.handlePubrec()
	case mqtt.PubrelCode:
		// err = client.handlePubrel()
	case mqtt.PubcompCode:
		// err = client.handlePubcomp()
	case mqtt.SubscribeCode:
		// err = client.handleSubscribe()
	case mqtt.SubackCode:
		// err = client.handleSuback()
	case mqtt.UnsubscribeCode:
		// err = client.handleUnsubscribe()
	case mqtt.UnsubackCode:
		// err = client.handleUnsuback()
	case mqtt.PingreqCode:
		// err = client.handlePingreq()
	case mqtt.PingrespCode:
		// err = client.handlePingresp()
	case mqtt.DisconnectCode:
		// err = client.handleDisconnect()
	case mqtt.AuthCode:
		// err = client.handleAuth()
	default:
		msg := fmt.Sprintf("server cannot send packet of type: %d", packetCode)
		return nil, errors.New(msg)
	}
	return header, err
}
