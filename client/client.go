package client

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/M4THYOU/some_mqtt_broker/mqtt"
	"github.com/M4THYOU/some_mqtt_broker/packet"
)

type Client struct {
	Conn         net.Conn
	Rdr          *packet.Reader
	connectFlags *mqtt.ConnectFlags
	KeepAlive    uint16

	// Connect Properties
	SessionExpiryInterval uint32
	ReceiveMaximum        uint16
	MaxPacketSize         uint32 // because go is go, the default value is 0. However, according to spec the client must not set this to zero. Therefore, 0 => no limit.
	TopicAliasMaximum     uint16
	ReturnResponseInfo    bool // if true, return response info to the client in connack.
	ReturnProblemInfo     bool // if true, maybe send proper reason string to client (CONNACK & DISCONNECT). Depends on log level though.
	AuthMethod            string
	AuthData              []byte

	WillProps *mqtt.WillProps
}

func (client *Client) handleConnect() error {
	fmt.Println("Handle Connect")

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
		// _, willTopic, err := client.Rdr.ReadUtf8Str()
		if err != nil {
			return err
		}
		// will payload, binary data

	}

	return nil
}
func (client *Client) handleConnack() error {
	fmt.Println("Handle Connack")
	log.Fatalln("Not yet implemented.")
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
	fmt.Println("Handle Subscribe")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleSuback() error {
	fmt.Println("Handle Suback")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleUnsubscribe() error {
	fmt.Println("Handle Unsubscribe")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handleUnsuback() error {
	fmt.Println("Handle Unsuback")
	log.Fatalln("Not yet implemented.")
	return nil
}
func (client *Client) handlePingreq() error {
	fmt.Println("Handle Pingreq")
	log.Fatalln("Not yet implemented.")
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

// processFixedHeader processes the fixed header.
// Returns the request type code, remaining length of the packet, and maybe an error.
func (client *Client) processFixedHeader() (byte, int, error) {
	fmt.Println("(fixed header)")
	b1, err := client.Rdr.ReadByte()
	if err != nil {
		return 0x00, 0, err
	}

	reqType := mqtt.GetRequestType(b1)
	if reqType == mqtt.PublishCode {
		// TODO Do thing with flags!
		// flags := (b1 & 0xF)
		log.Fatalln("Not yet implemented.")
	}

	_, remainingLength, err := client.Rdr.ReadVarByteInt()
	if err != nil {
		return 0x00, 0, err
	}

	return reqType, int(remainingLength), nil
}

// processVarHeader processes the variable header and payload of the packet (if payload exists)
func (client *Client) processVarHeader(reqType byte) (err error) {
	fmt.Println("(var header and payload)")

	switch reqType {
	case mqtt.ConnectCode:
		err = client.handleConnect()
	case mqtt.ConnackCode:
		err = client.handleConnack()
	case mqtt.PublishCode:
		err = client.handlePublish()
	case mqtt.PubackCode:
		err = client.handlePuback()
	case mqtt.PubrecCode:
		err = client.handlePubrec()
	case mqtt.PubrelCode:
		err = client.handlePubrel()
	case mqtt.PubcompCode:
		err = client.handlePubcomp()
	case mqtt.SubscribeCode:
		err = client.handleSubscribe()
	case mqtt.SubackCode:
		err = client.handleSuback()
	case mqtt.UnsubscribeCode:
		err = client.handleUnsubscribe()
	case mqtt.UnsubackCode:
		err = client.handleUnsuback()
	case mqtt.PingreqCode:
		err = client.handlePingreq()
	case mqtt.PingrespCode:
		err = client.handlePingresp()
	case mqtt.DisconnectCode:
		err = client.handleDisconnect()
	case mqtt.AuthCode:
		err = client.handleAuth()
	default:
		msg := fmt.Sprintf("No matching case for request type: %d", reqType)
		return errors.New(msg)
	}
	return err

}

func (client *Client) ProcessPacket() error {
	// fixed header can be up to 5 bytes, so set that as the limit.
	client.Rdr.SetRemainingLength(5)
	reqType, remLen, err := client.processFixedHeader() // make this guy return remaining length!
	if err != nil {
		return err
	}
	client.Rdr.SetRemainingLength(remLen)
	err = client.processVarHeader(reqType)
	return err
}
