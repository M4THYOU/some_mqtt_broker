package client

import (
	"testing"

	"github.com/M4THYOU/some_mqtt_broker/pkg/mqtt"
	"github.com/google/go-cmp/cmp"
)

func TestSetProperties(t *testing.T) {

}

func checkConnectProps(t *testing.T, expectedClient *Client, props map[int][]byte, shouldPass bool) {
	c := &Client{} // create dummy client
	err := c.setConnectProps(props)
	if err != nil && shouldPass {
		t.Fatalf("setConnectProps failed: %v", err.Error())
	}
	if c.SessionExpiryInterval != expectedClient.SessionExpiryInterval && shouldPass {
		t.Fatalf("incorrect SessionExpiryInterval. Got %v expected %v", c.SessionExpiryInterval, expectedClient.SessionExpiryInterval)
	}
	if c.ReceiveMaximum != expectedClient.ReceiveMaximum && shouldPass {
		t.Fatalf("incorrect ReceiveMaximum. Got %v expected %v", c.ReceiveMaximum, expectedClient.ReceiveMaximum)
	}
	if c.MaxPacketSize != expectedClient.MaxPacketSize && shouldPass {
		t.Fatalf("incorrect MaxPacketSize. Got %v expected %v", c.MaxPacketSize, expectedClient.MaxPacketSize)
	}
	if c.TopicAliasMaximum != expectedClient.TopicAliasMaximum && shouldPass {
		t.Fatalf("incorrect TopicAliasMaximum. Got %v expected %v", c.TopicAliasMaximum, expectedClient.TopicAliasMaximum)
	}
	if c.ReturnResponseInfo != expectedClient.ReturnResponseInfo && shouldPass {
		t.Fatalf("incorrect ReturnResponseInfo. Got %v expected %v", c.ReturnResponseInfo, expectedClient.ReturnResponseInfo)
	}
	if c.ReturnProblemInfo != expectedClient.ReturnProblemInfo && shouldPass {
		t.Fatalf("incorrect ReturnProblemInfo. Got %v expected %v", c.ReturnProblemInfo, expectedClient.ReturnProblemInfo)
	}
	if c.AuthMethod != expectedClient.AuthMethod && shouldPass {
		t.Fatalf("incorrect AuthMethod. Got %v expected %v", c.AuthMethod, expectedClient.AuthMethod)
	}
	if !cmp.Equal(c.AuthData, expectedClient.AuthData) && shouldPass {
		t.Fatalf("incorrect AuthData. Got %v expected %v", c.AuthData, expectedClient.AuthData)
	}
}
func TestSetConnectProps(t *testing.T) {
	// No payload
	var nilSlice []byte
	expectedClient := &Client{
		SessionExpiryInterval: 0,
		ReceiveMaximum:        65535,
		MaxPacketSize:         0,
		TopicAliasMaximum:     0,
		ReturnResponseInfo:    false,
		ReturnProblemInfo:     true,
		AuthMethod:            "",
		AuthData:              nilSlice,
	}
	props := map[int][]byte{}
	checkConnectProps(t, expectedClient, props, true)

	// Full Payload
	expectedClient = &Client{
		SessionExpiryInterval: 500,
		ReceiveMaximum:        54,
		MaxPacketSize:         1999999999,
		TopicAliasMaximum:     4,
		ReturnResponseInfo:    true,
		ReturnProblemInfo:     false,
		AuthMethod:            "SCRAM-SHA-1",
		AuthData:              []byte{0x04, 0x6d},
	}
	props = map[int][]byte{
		mqtt.SessionExpiryIntervalCode: {0x00, 0x00, 0x01, 0xF4},
		mqtt.ReceiveMaxCode:            {0x00, 0x36},
		mqtt.MaxPacketSizeCode:         {0x77, 0x35, 0x93, 0xFF},
		mqtt.TopicAliasMaxCode:         {0x00, 0x04},
		mqtt.RequestResponseInfoCode:   {0x01},
		mqtt.RequestProblemInfoCode:    {0x00},
		mqtt.AuthenticationMethodCode:  {0x53, 0x43, 0x52, 0x41, 0x4d, 0x2d, 0x53, 0x48, 0x41, 0x2d, 0x31},
		mqtt.AuthenticationDataCode:    {0x04, 0x6d},
	}
	checkConnectProps(t, expectedClient, props, true)

	// check each fail condition.
	props[mqtt.ReceiveMaxCode] = []byte{0x00, 0x00}
	checkConnectProps(t, expectedClient, props, false)
	props[mqtt.ReceiveMaxCode] = []byte{0x00, 0x36}

	props[mqtt.MaxPacketSizeCode] = []byte{0x00, 0x00, 0x00, 0x00}
	checkConnectProps(t, expectedClient, props, false)
	props[mqtt.MaxPacketSizeCode] = []byte{0x77, 0x35, 0x93, 0xFF}

	props[mqtt.RequestResponseInfoCode] = []byte{0x03}
	checkConnectProps(t, expectedClient, props, false)
	props[mqtt.RequestResponseInfoCode] = []byte{0x01}
	props[mqtt.RequestProblemInfoCode] = []byte{0x04}
	checkConnectProps(t, expectedClient, props, false)
	props[mqtt.RequestProblemInfoCode] = []byte{0x00}

	delete(props, mqtt.AuthenticationMethodCode)
	checkConnectProps(t, expectedClient, props, false)
	props[mqtt.AuthenticationMethodCode] = []byte{0x53, 0x43, 0x52, 0x41, 0x4d, 0x2d, 0x53, 0x48, 0x41, 0x2d, 0x31}
}

func checkWillProps(t *testing.T, expectedClient *Client, props map[int][]byte, shouldPass bool) {
	c := &Client{} // create dummy client
	err := c.setWillProps(props)
	if err != nil && shouldPass {
		t.Fatalf("setWillProps failed: %v", err.Error())
	}
	if c.WillProps.WillDelayInterval != expectedClient.WillProps.WillDelayInterval && shouldPass {
		t.Fatalf("incorrect WillDelayInterval. Got %v expected %v", c.WillProps.WillDelayInterval, expectedClient.WillProps.WillDelayInterval)
	}
	if c.WillProps.PayloadFormatIndicator != expectedClient.WillProps.PayloadFormatIndicator && shouldPass {
		t.Fatalf("incorrect PayloadFormatIndicator. Got %v expected %v", c.WillProps.PayloadFormatIndicator, expectedClient.WillProps.PayloadFormatIndicator)
	}
	if c.WillProps.MessageExpiryInterval != expectedClient.WillProps.MessageExpiryInterval && shouldPass {
		t.Fatalf("incorrect MessageExpiryInterval. Got %v expected %v", c.WillProps.MessageExpiryInterval, expectedClient.WillProps.MessageExpiryInterval)
	}
	if c.WillProps.ContentType != expectedClient.WillProps.ContentType && shouldPass {
		t.Fatalf("incorrect ContentType. Got %v expected %v", c.WillProps.ContentType, expectedClient.WillProps.ContentType)
	}
	if c.WillProps.ResponseTopic != expectedClient.WillProps.ResponseTopic && shouldPass {
		t.Fatalf("incorrect ResponseTopic. Got %v expected %v", c.WillProps.ResponseTopic, expectedClient.WillProps.ResponseTopic)
	}
	if !cmp.Equal(c.WillProps.CorrelationData, expectedClient.WillProps.CorrelationData) && shouldPass {
		t.Fatalf("incorrect CorrelationData. Got %v expected %v", c.WillProps.CorrelationData, expectedClient.WillProps.CorrelationData)
	}
}
func TestSetWillProps(t *testing.T) {
	// No payload
	var nilSlice []byte
	expectedClient := &Client{
		WillProps: &mqtt.WillProps{
			WillDelayInterval:      0,
			PayloadFormatIndicator: 0,
			MessageExpiryInterval:  0,
			ContentType:            "",
			ResponseTopic:          "",
			CorrelationData:        nilSlice,
		},
	}
	props := map[int][]byte{}
	checkWillProps(t, expectedClient, props, true)

	// Full Payload
	expectedClient = &Client{
		WillProps: &mqtt.WillProps{
			WillDelayInterval:      2148343340,
			PayloadFormatIndicator: 1,
			MessageExpiryInterval:  60,
			ContentType:            "json",
			ResponseTopic:          "my/response/topic",
			CorrelationData:        []byte{0x02, 0xFF, 0x6B},
		},
	}
	props = map[int][]byte{
		mqtt.WillDelayIntervalCode:      {0x80, 0x0D, 0x1E, 0x2C},
		mqtt.PayloadFormatIndicatorCode: {0x01},
		mqtt.MessageExpiryIntervalCode:  {0x00, 0x00, 0x00, 0x3C},
		mqtt.ContentTypeCode:            {0x6a, 0x73, 0x6f, 0x6e},
		mqtt.ResponseTopicCode:          {0x6d, 0x79, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2f, 0x74, 0x6f, 0x70, 0x69, 0x63},
		mqtt.CorrelationDataCode:        {0x02, 0xFF, 0x6B},
	}
	checkWillProps(t, expectedClient, props, true)

	// check each fail condition.
	props[mqtt.PayloadFormatIndicatorCode] = []byte{0x03}
	checkWillProps(t, expectedClient, props, false)
	props[mqtt.PayloadFormatIndicatorCode] = []byte{0x01}
}
