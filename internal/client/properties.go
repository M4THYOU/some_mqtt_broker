package client

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/M4THYOU/some_mqtt_broker/internal/defaults"
	"github.com/M4THYOU/some_mqtt_broker/pkg/mqtt"
)

func (client *Client) setConnectProps(props map[int][]byte) error {
	if v, ok := props[mqtt.SessionExpiryIntervalCode]; ok {
		client.SessionExpiryInterval = binary.BigEndian.Uint32(v)
	} else {
		client.SessionExpiryInterval = defaults.DefaultSessionExpiryInterval
	}
	if v, ok := props[mqtt.ReceiveMaxCode]; ok {
		v_i := binary.BigEndian.Uint16(v)
		if v_i <= 0 {
			msg := fmt.Sprintf("setProperties: invalid ReceiveMax %d for Connect", v_i)
			return errors.New(msg)
		}
		client.ReceiveMaximum = v_i
	} else {
		client.ReceiveMaximum = defaults.DefaultReceiveMaximum
	}
	if v, ok := props[mqtt.MaxPacketSizeCode]; ok {
		v_i := binary.BigEndian.Uint32(v)
		if v_i <= 0 {
			msg := fmt.Sprintf("setProperties: invalid MaxPacketSize %d for Connect", v_i)
			return errors.New(msg)
		}
		client.MaxPacketSize = v_i
	} // no default, unlimited.
	if v, ok := props[mqtt.TopicAliasMaxCode]; ok {
		client.TopicAliasMaximum = binary.BigEndian.Uint16(v)
	} else {
		client.TopicAliasMaximum = defaults.DefaultTopicAliasMaximum
	}
	if v, ok := props[mqtt.RequestResponseInfoCode]; ok {
		v_i := uint(v[0]) // it's a single byte that can only be 0 or 1.
		if v_i != 0 && v_i != 1 {
			msg := fmt.Sprintf("setProperties: invalid RequestResponseInfo %d for Connect", v_i)
			return errors.New(msg)
		}
		client.ReturnResponseInfo = (v_i == 1)
	} else {
		client.ReturnResponseInfo = defaults.DefaultRequestResponseInfo
	}
	if v, ok := props[mqtt.RequestProblemInfoCode]; ok {
		v_i := uint(v[0]) // it's a single byte that can only be 0 or 1.
		if v_i != 0 && v_i != 1 {
			msg := fmt.Sprintf("setProperties: invalid RequestProblemInfo %d for Connect", v_i)
			return errors.New(msg)
		}
		client.ReturnProblemInfo = (v_i == 1)
	} else {
		client.ReturnProblemInfo = defaults.DefaultRequestProblemInfo
	}
	if v, ok := props[mqtt.AuthenticationMethodCode]; ok {
		client.AuthMethod = string(v)
	}
	if v, ok := props[mqtt.AuthenticationDataCode]; ok {
		if client.AuthMethod == "" {
			return errors.New("setProperties: cannot set AuthData without AuthMethod")
		}
		client.AuthData = v
	}
	return nil
}

func (client *Client) setConnackProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Connack")
}
func (client *Client) setPublishProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Publish")
}
func (client *Client) setPubackProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Puback")
}
func (client *Client) setPubrecProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Pubrec")
}
func (client *Client) setPubrelProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Pubrel")
}
func (client *Client) setPubcompProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Pubcomp")
}
func (client *Client) setSubscribeProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Subscribe")
}
func (client *Client) setSubackProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Suback")
}
func (client *Client) setUnsubscribeProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Unsubscribe")
}
func (client *Client) setUnsubackProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Unsuback")
}
func (client *Client) setPingreqProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Pingreq")
}
func (client *Client) setPingrespProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Pingresp")
}
func (client *Client) setDisconnectProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Disconnect")
}
func (client *Client) setAuthProps(props map[int][]byte) error {
	return errors.New("setProperties: Case not implemented yet for Auth")
}

// Not part of spec.
func (client *Client) setWillProps(props map[int][]byte) error {
	if client.WillProps == nil {
		client.WillProps = &mqtt.WillProps{}
	}
	if v, ok := props[mqtt.PayloadFormatIndicatorCode]; ok {
		v_i := uint8(v[0]) // it's a single byte that can only be 0 or 1.
		if v_i != 0 && v_i != 1 {
			msg := fmt.Sprintf("setProperties: invalid PayloadFormatIndicator %d for WillProps", v_i)
			return errors.New(msg)
		}
		client.WillProps.PayloadFormatIndicator = v_i
	}
	if v, ok := props[mqtt.MessageExpiryIntervalCode]; ok {
		client.WillProps.MessageExpiryInterval = binary.BigEndian.Uint32(v)
	} else {
		client.WillProps.MessageExpiryInterval = defaults.DefaultMessageExpiryInterval
	}
	if v, ok := props[mqtt.ContentTypeCode]; ok {
		client.WillProps.ContentType = string(v)
	}
	if v, ok := props[mqtt.ResponseTopicCode]; ok {
		client.WillProps.ResponseTopic = string(v)
	}
	if v, ok := props[mqtt.CorrelationDataCode]; ok {
		client.WillProps.CorrelationData = v
	}
	if v, ok := props[mqtt.WillDelayIntervalCode]; ok {
		client.WillProps.WillDelayInterval = binary.BigEndian.Uint32(v)
	} else {
		client.WillProps.WillDelayInterval = defaults.DefaultWillDelayInterval
	}
	return nil
}

func (client *Client) setProperties(packetType int, props map[int][]byte) (err error) {
	// need a switch for each packet type.
	switch packetType {
	case mqtt.ConnectCode:
		err = client.setConnectProps(props)
	case mqtt.ConnackCode:
		err = client.setConnackProps(props)
	case mqtt.PublishCode:
		err = client.setPublishProps(props)
	case mqtt.PubackCode:
		err = client.setPubackProps(props)
	case mqtt.PubrecCode:
		err = client.setPubrecProps(props)
	case mqtt.PubrelCode:
		err = client.setPubrelProps(props)
	case mqtt.PubcompCode:
		err = client.setPubcompProps(props)
	case mqtt.SubscribeCode:
		err = client.setSubscribeProps(props)
	case mqtt.SubackCode:
		err = client.setSubackProps(props)
	case mqtt.UnsubscribeCode:
		err = client.setUnsubscribeProps(props)
	case mqtt.UnsubackCode:
		err = client.setUnsubackProps(props)
	case mqtt.PingreqCode:
		err = client.setPingreqProps(props)
	case mqtt.PingrespCode:
		err = client.setPingrespProps(props)
	case mqtt.DisconnectCode:
		err = client.setDisconnectProps(props)
	case mqtt.AuthCode:
		err = client.setAuthProps(props)
	case mqtt.WillPropsCode: // not defined by the spec, but we use this in getProps.:
		err = client.setWillProps(props)
	default:
		msg := fmt.Sprintf("setProperties: No matching case for packet type: %d", packetType)
		err = errors.New(msg)
	}
	return err
}
