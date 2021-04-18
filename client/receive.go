package client

import (
	"errors"
	"fmt"
	"log"

	"github.com/M4THYOU/some_mqtt_broker/mqtt"
)

func (client *Client) handleConnect() error {
	// verify the protocol is set to 'MQTT'
	err := mqtt.VerifyProtocol(client.Rdr)
	if err != nil {
		return err
	}

	// Check protocol version. Currently only supports v5.0
	b, err := client.Rdr.ReadByte()
	if err != nil {
		return err
	} else if b != 5 {
		msg := fmt.Sprintf("This broker currently only supports MQTT v5.0. You specified: %d", b)
		return errors.New(msg)
	}

	// Check the connect flags!
	b, err = client.Rdr.ReadByte()
	if err != nil {
		return err
	}
	flags, err := mqtt.GetConnectFlags(b)
	if err != nil {
		return err
	}
	client.connectFlags = flags

	// KeepAlive
	keepAlive, err := mqtt.GetKeepAlive(client.Rdr)
	if err != nil {
		return err
	}
	client.KeepAlive = keepAlive

	// Handle the properties!
	_, propLength, err := client.Rdr.ReadVarByteInt()
	if err != nil {
		return err
	}
	props, userProps, err := mqtt.GetProps(client.Rdr, int(propLength), mqtt.ConnectCode)
	if err != nil {
		return err
	}
	fmt.Printf("Flags: %v\n", client.connectFlags)
	fmt.Printf("Props: %v\n", props)
	fmt.Printf("User Props: %v\n", userProps)
	client.setProperties(mqtt.ConnectCode, props)

	//// Process the payload ////

	// first thing is a utf8 encoded string for the clientId.
	// If zero length, assign one from the server.
	clientId, err := mqtt.GetClientId(client.Rdr)
	if err != nil {
		return err
	}
	fmt.Printf("Client ID: %v\n", clientId)

	// Check for will things in the payload.
	if client.connectFlags.WillFlag {
		_, propLength, err = client.Rdr.ReadVarByteInt()
		if err != nil {
			return err
		}
		willProps, _, err := mqtt.GetProps(client.Rdr, int(propLength), mqtt.WillPropsCode)
		if err != nil {
			return err
		}
		client.setWillProps(willProps)

		// will topic, UTF-8 enc string
		_, willTopic, err := client.Rdr.ReadUtf8Str()
		if err != nil {
			return err
		}
		client.WillProps.Topic = willTopic
		// will payload, binary data
		_, willPayload, err := client.Rdr.ReadBinaryData()
		if err != nil {
			return err
		}
		client.WillProps.Payload = willPayload
	}
	if client.connectFlags.UserNameFlag {
		_, userName, err := client.Rdr.ReadUtf8Str()
		if err != nil {
			return err
		}
		client.UserName = userName
	}
	if client.connectFlags.PasswordFlag {
		_, password, err := client.Rdr.ReadBinaryData()
		if err != nil {
			return err
		}
		client.Password = password
	}

	// send a CONNACK packet.
	connack, err := client.BuildPacket(mqtt.ConnackCode)
	if err != nil {
		return err
	}
	mqtt.SendPacket(client.Conn, connack)

	return nil
}
func (client *Client) handleConnack() error {
	log.Fatalln("Invalid Operation.")
	return nil
}
func (client *Client) handlePublish() error {
	fmt.Println("Handle Publish")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePuback() error {
	fmt.Println("Handle Puback")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePubrec() error {
	fmt.Println("Handle Pubrec")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePubrel() error {
	fmt.Println("Handle Pubrel")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePubcomp() error {
	fmt.Println("Handle Pubcomp")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleSubscribe() error {
	log.Fatalln("Invalid Operation.")
	return nil
}
func (client *Client) handleSuback() error {
	fmt.Println("Handle Suback")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleUnsubscribe() error {
	log.Fatalln("Invalid Operation.")
	return nil
}
func (client *Client) handleUnsuback() error {
	fmt.Println("Handle Unsuback")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePingreq() error {
	log.Fatalln("Invalid Operation.")
	return nil
}
func (client *Client) handlePingresp() error {
	fmt.Println("Handle Pingresp")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleDisconnect() error {
	fmt.Println("Handle Disconnect")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleAuth() error {
	fmt.Println("Handle Auth")
	log.Fatalln("Not yet implemented.")
	return nil
}
