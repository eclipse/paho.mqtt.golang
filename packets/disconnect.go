package packets

import (
	"bytes"
	"io"
	"net"
)

// Disconnect is the Variable Header definition for a Disconnect control packet
type Disconnect struct {
	ReasonCode byte
	Properties Properties
}

// DisconnectNormalDisconnection, etc are the list of valid disconnection reason codes.
const (
	DisconnectNormalDisconnection                 = 0x00
	DisconnectDisconnectWithWillMessage           = 0x04
	DisconnectUnspecifiedError                    = 0x80
	DisconnectMalformedPacket                     = 0x81
	DisconnectProtocolError                       = 0x82
	DisconnectImplementationSpecificError         = 0x83
	DisconnectNotAuthorized                       = 0x87
	DisconnectServerBusy                          = 0x89
	DisconnectServerShuttingDown                  = 0x8B
	DisconnectKeepAliveTimeout                    = 0x8D
	DisconnectSessionTakenOver                    = 0x8E
	DisconnectTopicFilterInvalid                  = 0x8F
	DisconnectTopicNameInvalid                    = 0x90
	DisconnectReceiveMaximumExceeded              = 0x93
	DisconnectTopicAliasInvalid                   = 0x94
	DisconnectPacketTooLarge                      = 0x95
	DisconnectMessageRateTooHigh                  = 0x96
	DisconnectQuotaExceeded                       = 0x97
	DisconnectAdministrativeAction                = 0x98
	DisconnectPayloadFormatInvalid                = 0x99
	DisconnectRetainNotSupported                  = 0x9A
	DisconnectQoSNotSupported                     = 0x9B
	DisconnectUseAnotherServer                    = 0x9C
	DisconnectServerMoved                         = 0x9D
	DisconnectSharedSubscriptionNotSupported      = 0x9E
	DisconnectConnectionRateExceeded              = 0x9F
	DisconnectMaximumConnectTime                  = 0xA0
	DisconnectSubscriptionIdentifiersNotSupported = 0xA1
	DisconnectWildcardSubscriptionsNotSupported   = 0xA2
)

// NewDisconnect creates a new Disconnect packet and applies all the
// provided/listed option functions to configure the packet
func NewDisconnect(opts ...func(*Disconnect)) *Disconnect {
	d := &Disconnect{}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

// DisconnectReason is a Disconnect option function that sets the
// reason code for the Disconnect packet
func DisconnectReason(r byte) func(*Disconnect) {
	return func(d *Disconnect) {
		d.ReasonCode = r
	}
}

// DisconnectProperties is a Disconnect option function that sets
// the Properties for the Disconnect packet
func DisconnectProperties(p Properties) func(*Disconnect) {
	return func(d *Disconnect) {
		d.Properties = p
	}
}

// Unpack is the implementation of the interface required function for a packet
func (d *Disconnect) Unpack(r *bytes.Buffer) error {
	var err error
	d.ReasonCode, err = r.ReadByte()
	if err != nil {
		return err
	}

	err = d.Properties.Unpack(r, DISCONNECT)
	if err != nil {
		return err
	}

	return nil
}

// Buffers is the implementation of the interface required function for a packet
func (d *Disconnect) Buffers() net.Buffers {
	idvp := d.Properties.Pack(DISCONNECT)
	propLen := encodeVBI(len(idvp))
	return net.Buffers{[]byte{d.ReasonCode}, propLen, idvp}
}

// Send is the implementation of the interface required function for a packet
func (d *Disconnect) Send(w io.Writer) error {
	cp := &ControlPacket{FixedHeader: FixedHeader{Type: DISCONNECT}}
	cp.Content = d

	return cp.Send(w)
}

// Reason returns a string representation of the meaning of the ReasonCode
func (d *Disconnect) Reason() string {
	switch d.ReasonCode {
	case 0:
		return "Normal disconnection - Close the connection normally. Do not send the Will Message."
	case 4:
		return "Disconnect with Will Message - The Client wishes to disconnect but requires that the Server also publishes its Will Message."
	case 128:
		return "Unspecified error - The Connection is closed but the sender either does not wish to reveal the reason, or none of the other Reason Codes apply."
	case 129:
		return "Malformed Packet - The received packet does not conform to this specification."
	case 130:
		return "Protocol Error - An unexpected or out of order packet was received."
	case 131:
		return "Implementation specific error - The packet received is valid but cannot be processed by this implementation."
	case 135:
		return "Not authorized - The request is not authorized."
	case 137:
		return "Server busy - The Server is busy and cannot continue processing requests from this Client."
	case 139:
		return "Server shutting down - The Server is shutting down."
	case 141:
		return "Keep Alive timeout - The Connection is closed because no packet has been received for 1.5 times the Keepalive time."
	case 142:
		return "Session taken over - Another Connection using the same ClientID has connected causing this Connection to be closed."
	case 143:
		return "Topic Filter invalid - The Topic Filter is correctly formed, but is not accepted by this Sever."
	case 144:
		return "Topic Name invalid - The Topic Name is correctly formed, but is not accepted by this Client or Server."
	case 147:
		return "Receive Maximum exceeded - The Client or Server has received more than Receive Maximum publication for which it has not sent PUBACK or PUBCOMP."
	case 148:
		return "Topic Alias invalid - The Client or Server has received a PUBLISH packet containing a Topic Alias which is greater than the Maximum Topic Alias it sent in the CONNECT or CONNACK packet."
	case 149:
		return "Packet too large - The packet size is greater than Maximum Packet Size for this Client or Server."
	case 150:
		return "Message rate too high - The received data rate is too high."
	case 151:
		return "Quota exceeded - An implementation or administrative imposed limit has been exceeded."
	case 152:
		return "Administrative action - The Connection is closed due to an administrative action."
	case 153:
		return "Payload format invalid - The payload format does not match the one specified by the Payload Format Indicator."
	case 154:
		return "Retain not supported - The Server has does not support retained messages."
	case 155:
		return "QoS not supported - The Client specified a QoS greater than the QoS specified in a Maximum QoS in the CONNACK."
	case 156:
		return "Use another server - The Client should temporarily change its Server."
	case 157:
		return "Server moved - The Server is moved and the Client should permanently change its server location."
	case 158:
		return "Shared Subscription not supported - The Server does not support Shared Subscriptions."
	case 159:
		return "Connection rate exceeded - This connection is closed because the connection rate is too high."
	case 160:
		return "Maximum connect time - The maximum connection time authorized for this connection has been exceeded."
	case 161:
		return "Subscription Identifiers not supported - The Server does not support Subscription Identifiers; the subscription is not accepted."
	case 162:
		return "Wildcard subscriptions not supported - The Server does not support Wildcard subscription; the subscription is not accepted."
	}

	return ""
}
