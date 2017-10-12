package packets

import (
	"bytes"
	"fmt"
	"io"
)

const (
	IDVPPayloadFormat          byte = 1
	IDVPPubExpiry                   = 2
	IDVPContentType                 = 3
	IDVPReplyTopic                  = 8
	IDVPCorrelationData             = 9
	IDVPSubscriptionIdentifier      = 11
	IDVPSessionExpiryInterval       = 17
	IDVPAssignedClientID            = 18
	IDVPServerKeepAlive             = 19
	IDVPAuthMethod                  = 21
	IDVPAuthData                    = 22
	IDVPRequestProblemInfo          = 23
	IDVPWillDelayInterval           = 24
	IDVPRequestResponseInfo         = 25
	IDVPResponseInfo                = 26
	IDVPServerReference             = 28
	IDVPReasonString                = 31
	IDVPReceiveMaximum              = 33
	IDVPTopicAliasMaximum           = 34
	IDVPTopicAlias                  = 35
	IDVPMaximumQOS                  = 36
	IDVPRetainAvailable             = 37
	IDVPUserProperty                = 38
	IDVPMaximumPacketSize           = 39
	IDVPWildcardSubAvailable        = 40
	IDVPSubIDAvailable              = 41
	IDVPSharedSubAvailable          = 42
)

// IDValuePair is a struct representing the all the described properties
// allowed by the MQTT protocol, determining the validity of a property
// relvative to the packettype it was received in is provided by the
// ValidateID function
type IDValuePair struct {
	PayloadFormat          *byte
	PubExpiry              *uint32
	ContentType            string
	ReplyTopic             string
	CorrelationData        []byte
	SubscriptionIdentifier *uint32
	SessionExpiryInterval  *uint32
	AssignedClientID       string
	ServerKeepAlive        *uint16
	AuthMethod             string
	AuthData               []byte
	RequestProblemInfo     *byte
	WillDelayInterval      *uint32
	RequestResponseInfo    *byte
	ResponseInfo           string
	ServerReference        string
	ReasonString           string
	ReceiveMaximum         *uint16
	TopicAliasMaximum      *uint16
	TopicAlias             *uint16
	MaximumQOS             *byte
	RetainAvailable        *byte
	UserProperty           map[string]string
	MaximumPacketSize      *uint32
	WildcardSubAvailable   *byte
	SubIDAvailable         *byte
	SharedSubAvailable     *byte
}

// Pack takes all the defined properties for an IDValuePair and produces
// a slice of bytes representing the wire format for the information
func (i *IDValuePair) Pack(p PacketType) []byte {
	var b bytes.Buffer

	if p == PUBLISH {
		if i.PayloadFormat != nil {
			b.WriteByte(IDVPPayloadFormat)
			b.WriteByte(*i.PayloadFormat)
		}

		if i.PubExpiry != nil {
			b.WriteByte(IDVPPubExpiry)
			writeUint32(*i.PubExpiry, &b)
		}

		if i.ContentType != "" {
			b.WriteByte(IDVPContentType)
			writeString(i.ContentType, &b)
		}

		if i.ReplyTopic != "" {
			b.WriteByte(IDVPReplyTopic)
			writeString(i.ReplyTopic, &b)
		}

		if i.CorrelationData != nil && len(i.CorrelationData) > 0 {
			b.WriteByte(IDVPCorrelationData)
			b.Write(i.CorrelationData)
		}

		if i.TopicAlias != nil {
			b.WriteByte(IDVPTopicAlias)
			writeUint16(*i.TopicAlias, &b)
		}
	}

	if p == PUBLISH || p == SUBSCRIBE {
		if i.SubscriptionIdentifier != nil {
			b.WriteByte(IDVPSubscriptionIdentifier)
			writeUint32(*i.SubscriptionIdentifier, &b)
		}
	}

	if p == CONNECT || p == CONNACK {
		if i.ReceiveMaximum != nil {
			b.WriteByte(IDVPReceiveMaximum)
			writeUint16(*i.ReceiveMaximum, &b)
		}

		if i.TopicAliasMaximum != nil {
			b.WriteByte(IDVPTopicAliasMaximum)
			writeUint16(*i.TopicAliasMaximum, &b)
		}

		if i.MaximumQOS != nil {
			b.WriteByte(IDVPMaximumQOS)
			b.WriteByte(*i.MaximumQOS)
		}

		if i.MaximumPacketSize != nil {
			b.WriteByte(IDVPMaximumPacketSize)
			writeUint32(*i.MaximumPacketSize, &b)
		}
	}

	if p == CONNACK {
		if i.AssignedClientID != "" {
			b.WriteByte(IDVPAssignedClientID)
			writeString(i.AssignedClientID, &b)
		}

		if i.ServerKeepAlive != nil {
			b.WriteByte(IDVPServerKeepAlive)
			writeUint16(*i.ServerKeepAlive, &b)
		}

		if i.WildcardSubAvailable != nil {
			b.WriteByte(IDVPWildcardSubAvailable)
			b.WriteByte(*i.WildcardSubAvailable)
		}

		if i.SubIDAvailable != nil {
			b.WriteByte(IDVPSubIDAvailable)
			b.WriteByte(*i.SubIDAvailable)
		}

		if i.SharedSubAvailable != nil {
			b.WriteByte(IDVPSharedSubAvailable)
			b.WriteByte(*i.SharedSubAvailable)
		}

		if i.RetainAvailable != nil {
			b.WriteByte(IDVPRetainAvailable)
			b.WriteByte(*i.RetainAvailable)
		}

		if i.ResponseInfo != "" {
			b.WriteByte(IDVPResponseInfo)
			writeString(i.ResponseInfo, &b)
		}
	}

	if p == CONNECT {
		if i.RequestProblemInfo != nil {
			b.WriteByte(IDVPRequestProblemInfo)
			b.WriteByte(*i.RequestProblemInfo)
		}

		if i.WillDelayInterval != nil {
			b.WriteByte(IDVPWillDelayInterval)
			writeUint32(*i.WillDelayInterval, &b)
		}

		if i.RequestResponseInfo != nil {
			b.WriteByte(IDVPRequestResponseInfo)
			b.WriteByte(*i.RequestResponseInfo)
		}
	}

	if p == CONNECT || p == DISCONNECT {
		if i.SessionExpiryInterval != nil {
			b.WriteByte(IDVPSessionExpiryInterval)
			writeUint32(*i.SessionExpiryInterval, &b)
		}
	}

	if p == CONNECT || p == CONNACK || p == AUTH {
		if i.AuthMethod != "" {
			b.WriteByte(IDVPAuthMethod)
			writeString(i.AuthMethod, &b)
		}

		if i.AuthData != nil && len(i.AuthData) > 0 {
			b.WriteByte(IDVPAuthData)
			b.Write(i.AuthData)
		}
	}

	if p == CONNACK || p == DISCONNECT {
		if i.ServerReference != "" {
			b.WriteByte(IDVPServerReference)
			writeString(i.ServerReference, &b)
		}
	}

	if p != CONNECT {
		if i.ReasonString != "" {
			b.WriteByte(IDVPReasonString)
			writeString(i.ReasonString, &b)
		}
	}

	for k, v := range i.UserProperty {
		b.WriteByte(IDVPUserProperty)
		writeString(k, &b)
		writeString(v, &b)
	}

	return b.Bytes()
}

// Unpack takes a buffer of bytes and reads out the defined properties
// filling in the appropriate entries in the struct, it returns the number
// of bytes used to store the IDVP data and any error in decoding them
func (i *IDValuePair) Unpack(r *bytes.Buffer, p PacketType) (int, error) {
	vbi, err := getVBI(r)
	if err != nil {
		return 0, err
	}
	vbiLen := vbi.Len()
	size, err := decodeVBI(vbi)
	if err != nil {
		return 0, err
	}
	if size == 0 {
		return 1, nil
	}

	buf := bytes.NewBuffer(r.Next(size))
	for {
		IDVPType, err := buf.ReadByte()
		if err != nil && err != io.EOF {
			return 0, err
		}
		if err == io.EOF {
			break
		}
		if !ValidateID(p, IDVPType) {
			return 0, fmt.Errorf("Invalid IDVP type %d for packet %d", IDVPType, p)
		}
		switch IDVPType {
		case IDVPPayloadFormat:
			pf, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.PayloadFormat = &pf
		case IDVPPubExpiry:
			pe, err := readUint32(buf)
			if err != nil {
				return 0, err
			}
			i.PubExpiry = &pe
		case IDVPContentType:
			ct, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.ContentType = ct
		case IDVPReplyTopic:
			tr, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.ReplyTopic = tr
		case IDVPCorrelationData:
			cd, err := readBinary(buf)
			if err != nil {
				return 0, err
			}
			i.CorrelationData = cd
		case IDVPSubscriptionIdentifier:
			si, err := readUint32(buf)
			if err != nil {
				return 0, err
			}
			i.SubscriptionIdentifier = &si
		case IDVPSessionExpiryInterval:
			se, err := readUint32(buf)
			if err != nil {
				return 0, err
			}
			i.SessionExpiryInterval = &se
		case IDVPAssignedClientID:
			ac, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.AssignedClientID = ac
		case IDVPServerKeepAlive:
			sk, err := readUint16(buf)
			if err != nil {
				return 0, err
			}
			i.ServerKeepAlive = &sk
		case IDVPAuthMethod:
			am, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.AuthMethod = am
		case IDVPAuthData:
			ad, err := readBinary(buf)
			if err != nil {
				return 0, err
			}
			i.AuthData = ad
		case IDVPRequestProblemInfo:
			rp, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.RequestProblemInfo = &rp
		case IDVPWillDelayInterval:
			wd, err := readUint32(buf)
			if err != nil {
				return 0, err
			}
			i.WillDelayInterval = &wd
		case IDVPRequestResponseInfo:
			rp, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.RequestResponseInfo = &rp
		case IDVPResponseInfo:
			ri, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.ResponseInfo = ri
		case IDVPServerReference:
			sr, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.ServerReference = sr
		case IDVPReasonString:
			rs, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.ReasonString = rs
		case IDVPReceiveMaximum:
			rm, err := readUint16(buf)
			if err != nil {
				return 0, err
			}
			i.ReceiveMaximum = &rm
		case IDVPTopicAliasMaximum:
			ta, err := readUint16(buf)
			if err != nil {
				return 0, err
			}
			i.TopicAliasMaximum = &ta
		case IDVPTopicAlias:
			ta, err := readUint16(buf)
			if err != nil {
				return 0, err
			}
			i.TopicAlias = &ta
		case IDVPMaximumQOS:
			mq, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.MaximumQOS = &mq
		case IDVPRetainAvailable:
			ra, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.RetainAvailable = &ra
		case IDVPUserProperty:
			k, err := readString(buf)
			if err != nil {
				return 0, err
			}
			v, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.UserProperty[k] = v
		case IDVPMaximumPacketSize:
			mp, err := readUint32(buf)
			if err != nil {
				return 0, err
			}
			i.MaximumPacketSize = &mp
		case IDVPWildcardSubAvailable:
			ws, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.WildcardSubAvailable = &ws
		case IDVPSubIDAvailable:
			si, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.SubIDAvailable = &si
		case IDVPSharedSubAvailable:
			ss, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.SharedSubAvailable = &ss
		default:
			return 0, fmt.Errorf("Unknown IDVP type %d", IDVPType)
		}
	}

	return size + vbiLen, nil
}

// ValidIDValuePairs is a map of the various properties and the
// PacketTypes that property is valid for.
var ValidIDValuePairs = map[byte]map[PacketType]struct{}{
	IDVPPayloadFormat:          {PUBLISH: {}},
	IDVPPubExpiry:              {PUBLISH: {}},
	IDVPContentType:            {PUBLISH: {}},
	IDVPReplyTopic:             {PUBLISH: {}},
	IDVPCorrelationData:        {PUBLISH: {}},
	IDVPTopicAlias:             {PUBLISH: {}},
	IDVPSubscriptionIdentifier: {PUBLISH: {}, SUBSCRIBE: {}},
	IDVPSessionExpiryInterval:  {CONNECT: {}, DISCONNECT: {}},
	IDVPAssignedClientID:       {CONNACK: {}},
	IDVPServerKeepAlive:        {CONNACK: {}},
	IDVPWildcardSubAvailable:   {CONNACK: {}},
	IDVPSubIDAvailable:         {CONNACK: {}},
	IDVPSharedSubAvailable:     {CONNACK: {}},
	IDVPRetainAvailable:        {CONNACK: {}},
	IDVPResponseInfo:           {CONNACK: {}},
	IDVPAuthMethod:             {CONNECT: {}, CONNACK: {}, AUTH: {}},
	IDVPAuthData:               {CONNECT: {}, CONNACK: {}, AUTH: {}},
	IDVPRequestProblemInfo:     {CONNECT: {}},
	IDVPWillDelayInterval:      {CONNECT: {}},
	IDVPRequestResponseInfo:    {CONNECT: {}},
	IDVPServerReference:        {CONNACK: {}, DISCONNECT: {}},
	IDVPReasonString:           {CONNACK: {}, PUBACK: {}, PUBREC: {}, PUBREL: {}, PUBCOMP: {}, SUBACK: {}, UNSUBACK: {}, DISCONNECT: {}, AUTH: {}},
	IDVPReceiveMaximum:         {CONNECT: {}, CONNACK: {}},
	IDVPTopicAliasMaximum:      {CONNECT: {}, CONNACK: {}},
	IDVPMaximumQOS:             {CONNECT: {}, CONNACK: {}},
	IDVPMaximumPacketSize:      {CONNECT: {}, CONNACK: {}},
	IDVPUserProperty:           {CONNECT: {}, CONNACK: {}, PUBLISH: {}, PUBACK: {}, PUBREC: {}, PUBREL: {}, PUBCOMP: {}, SUBACK: {}, UNSUBACK: {}, DISCONNECT: {}, AUTH: {}},
}

// ValidateID takes a PacketType and a property name and returns
// a boolean indicating if that property is valid for that
// PacketType
func ValidateID(p PacketType, i byte) bool {
	_, ok := ValidIDValuePairs[i][p]
	return ok
}
