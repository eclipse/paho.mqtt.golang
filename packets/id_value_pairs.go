package packets

import (
	"bytes"
	"fmt"
	"io"
)

const (
	idvpPayloadFormat          byte = 1
	idvpPubExpiry                   = 2
	idvpContentType                 = 3
	idvpReplyTopic                  = 8
	idvpCorrelationData             = 9
	idvpSubscriptionIdentifier      = 11
	idvpSessionExpiryInterval       = 17
	idvpAssignedClientID            = 18
	idvpServerKeepAlive             = 19
	idvpAuthMethod                  = 21
	idvpAuthData                    = 22
	idvpRequestProblemInfo          = 23
	idvpWillDelayInterval           = 24
	idvpRequestResponseInfo         = 25
	idvpResponseInfo                = 26
	idvpServerReference             = 28
	idvpReasonString                = 31
	idvpReceiveMaximum              = 33
	idvpTopicAliasMaximum           = 34
	idvpTopicAlias                  = 35
	idvpMaximumQOS                  = 36
	idvpRetainAvailable             = 37
	idvpUserProperty                = 38
	idvpMaximumPacketSize           = 39
	idvpWildcardSubAvailable        = 40
	idvpSubIDAvailable              = 41
	idvpSharedSubAvailable          = 42
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
			b.WriteByte(idvpPayloadFormat)
			b.WriteByte(*i.PayloadFormat)
		}

		if i.PubExpiry != nil {
			b.WriteByte(idvpPubExpiry)
			writeUint32(*i.PubExpiry, &b)
		}

		if i.ContentType != "" {
			b.WriteByte(idvpContentType)
			b.WriteString(i.ContentType)
		}

		if i.ReplyTopic != "" {
			b.WriteByte(idvpReplyTopic)
			b.WriteString(i.ReplyTopic)
		}

		if i.CorrelationData != nil && len(i.CorrelationData) > 0 {
			b.WriteByte(idvpCorrelationData)
			b.Write(i.CorrelationData)
		}

		if i.TopicAlias != nil {
			b.WriteByte(idvpTopicAlias)
			writeUint16(*i.TopicAlias, &b)
		}
	}

	if p == PUBLISH || p == SUBSCRIBE {
		if i.SubscriptionIdentifier != nil {
			b.WriteByte(idvpSubscriptionIdentifier)
			writeUint32(*i.SubscriptionIdentifier, &b)
		}
	}

	if p == CONNECT || p == CONNACK {
		if i.ReceiveMaximum != nil {
			b.WriteByte(idvpReceiveMaximum)
			writeUint16(*i.ReceiveMaximum, &b)
		}

		if i.TopicAliasMaximum != nil {
			b.WriteByte(idvpTopicAliasMaximum)
			writeUint16(*i.TopicAliasMaximum, &b)
		}

		if i.MaximumQOS != nil {
			b.WriteByte(idvpMaximumQOS)
			b.WriteByte(*i.MaximumQOS)
		}

		if i.MaximumPacketSize != nil {
			b.WriteByte(idvpMaximumPacketSize)
			writeUint32(*i.MaximumPacketSize, &b)
		}
	}

	if p == CONNACK {
		if i.AssignedClientID != "" {
			b.WriteByte(idvpAssignedClientID)
			b.WriteString(i.AssignedClientID)
		}

		if i.ServerKeepAlive != nil {
			b.WriteByte(idvpServerKeepAlive)
			writeUint16(*i.ServerKeepAlive, &b)
		}

		if i.WildcardSubAvailable != nil {
			b.WriteByte(idvpWildcardSubAvailable)
			b.WriteByte(*i.WildcardSubAvailable)
		}

		if i.SubIDAvailable != nil {
			b.WriteByte(idvpSubIDAvailable)
			b.WriteByte(*i.SubIDAvailable)
		}

		if i.SharedSubAvailable != nil {
			b.WriteByte(idvpSharedSubAvailable)
			b.WriteByte(*i.SharedSubAvailable)
		}

		if i.RetainAvailable != nil {
			b.WriteByte(idvpRetainAvailable)
			b.WriteByte(*i.RetainAvailable)
		}

		if i.ResponseInfo != "" {
			b.WriteByte(idvpResponseInfo)
			b.WriteString(i.ResponseInfo)
		}
	}

	if p == CONNECT {
		if i.RequestProblemInfo != nil {
			b.WriteByte(idvpRequestProblemInfo)
			b.WriteByte(*i.RequestProblemInfo)
		}

		if i.WillDelayInterval != nil {
			b.WriteByte(idvpWillDelayInterval)
			writeUint32(*i.WillDelayInterval, &b)
		}

		if i.RequestResponseInfo != nil {
			b.WriteByte(idvpRequestResponseInfo)
			b.WriteByte(*i.RequestResponseInfo)
		}
	}

	if p == CONNECT || p == DISCONNECT {
		if i.SessionExpiryInterval != nil {
			b.WriteByte(idvpSessionExpiryInterval)
			writeUint32(*i.SessionExpiryInterval, &b)
		}
	}

	if p == CONNECT || p == CONNACK || p == AUTH {
		if i.AuthMethod != "" {
			b.WriteByte(idvpAuthMethod)
			b.WriteString(i.AuthMethod)
		}

		if i.AuthData != nil && len(i.AuthData) > 0 {
			b.WriteByte(idvpAuthData)
			b.Write(i.AuthData)
		}
	}

	if p == CONNACK || p == DISCONNECT {
		if i.ServerReference != "" {
			b.WriteByte(idvpServerReference)
			b.WriteString(i.ServerReference)
		}
	}

	if p != CONNECT {
		if i.ReasonString != "" {
			b.WriteByte(idvpReasonString)
			b.WriteString(i.ReasonString)
		}
	}

	for k, v := range i.UserProperty {
		b.WriteByte(idvpUserProperty)
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
		idvpType, err := buf.ReadByte()
		if err != nil && err != io.EOF {
			return 0, err
		}
		if err == io.EOF {
			break
		}
		if !ValidateID(p, idvpType) {
			return 0, fmt.Errorf("Invalid IDVP type %d for packet %d", idvpType, p)
		}
		switch idvpType {
		case idvpPayloadFormat:
			pf, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.PayloadFormat = &pf
		case idvpPubExpiry:
			pe, err := readUint32(buf)
			if err != nil {
				return 0, err
			}
			i.PubExpiry = &pe
		case idvpContentType:
			ct, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.ContentType = ct
		case idvpReplyTopic:
			tr, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.ReplyTopic = tr
		case idvpCorrelationData:
			cd, err := readBinary(buf)
			if err != nil {
				return 0, err
			}
			i.CorrelationData = cd
		case idvpSubscriptionIdentifier:
			si, err := readUint32(buf)
			if err != nil {
				return 0, err
			}
			i.SubscriptionIdentifier = &si
		case idvpSessionExpiryInterval:
			se, err := readUint32(buf)
			if err != nil {
				return 0, err
			}
			i.SessionExpiryInterval = &se
		case idvpAssignedClientID:
			ac, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.AssignedClientID = ac
		case idvpServerKeepAlive:
			sk, err := readUint16(buf)
			if err != nil {
				return 0, err
			}
			i.ServerKeepAlive = &sk
		case idvpAuthMethod:
			am, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.AuthMethod = am
		case idvpAuthData:
			ad, err := readBinary(buf)
			if err != nil {
				return 0, err
			}
			i.AuthData = ad
		case idvpRequestProblemInfo:
			rp, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.RequestProblemInfo = &rp
		case idvpWillDelayInterval:
			wd, err := readUint32(buf)
			if err != nil {
				return 0, err
			}
			i.WillDelayInterval = &wd
		case idvpRequestResponseInfo:
			rp, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.RequestResponseInfo = &rp
		case idvpResponseInfo:
			ri, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.ResponseInfo = ri
		case idvpServerReference:
			sr, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.ServerReference = sr
		case idvpReasonString:
			rs, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.ReasonString = rs
		case idvpReceiveMaximum:
			rm, err := readUint16(buf)
			if err != nil {
				return 0, err
			}
			i.ReceiveMaximum = &rm
		case idvpTopicAliasMaximum:
			ta, err := readUint16(buf)
			if err != nil {
				return 0, err
			}
			i.TopicAliasMaximum = &ta
		case idvpTopicAlias:
			ta, err := readUint16(buf)
			if err != nil {
				return 0, err
			}
			i.TopicAlias = &ta
		case idvpMaximumQOS:
			mq, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.MaximumQOS = &mq
		case idvpRetainAvailable:
			ra, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.RetainAvailable = &ra
		case idvpUserProperty:
			k, err := readString(buf)
			if err != nil {
				return 0, err
			}
			v, err := readString(buf)
			if err != nil {
				return 0, err
			}
			i.UserProperty[k] = v
		case idvpMaximumPacketSize:
			mp, err := readUint32(buf)
			if err != nil {
				return 0, err
			}
			i.MaximumPacketSize = &mp
		case idvpWildcardSubAvailable:
			ws, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.WildcardSubAvailable = &ws
		case idvpSubIDAvailable:
			si, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.SubIDAvailable = &si
		case idvpSharedSubAvailable:
			ss, err := buf.ReadByte()
			if err != nil {
				return 0, err
			}
			i.SharedSubAvailable = &ss
		default:
			return 0, fmt.Errorf("Unknown IDVP type %d", idvpType)
		}
	}

	return size + vbiLen, nil
}

// ValidIDValuePairs is a map of the various properties and the
// PacketTypes that property is valid for.
var ValidIDValuePairs = map[byte]map[PacketType]struct{}{
	idvpPayloadFormat:          {PUBLISH: {}},
	idvpPubExpiry:              {PUBLISH: {}},
	idvpContentType:            {PUBLISH: {}},
	idvpReplyTopic:             {PUBLISH: {}},
	idvpCorrelationData:        {PUBLISH: {}},
	idvpTopicAlias:             {PUBLISH: {}},
	idvpSubscriptionIdentifier: {PUBLISH: {}, SUBSCRIBE: {}},
	idvpSessionExpiryInterval:  {CONNECT: {}, DISCONNECT: {}},
	idvpAssignedClientID:       {CONNACK: {}},
	idvpServerKeepAlive:        {CONNACK: {}},
	idvpWildcardSubAvailable:   {CONNACK: {}},
	idvpSubIDAvailable:         {CONNACK: {}},
	idvpSharedSubAvailable:     {CONNACK: {}},
	idvpRetainAvailable:        {CONNACK: {}},
	idvpResponseInfo:           {CONNACK: {}},
	idvpAuthMethod:             {CONNECT: {}, CONNACK: {}, AUTH: {}},
	idvpAuthData:               {CONNECT: {}, CONNACK: {}, AUTH: {}},
	idvpRequestProblemInfo:     {CONNECT: {}},
	idvpWillDelayInterval:      {CONNECT: {}},
	idvpRequestResponseInfo:    {CONNECT: {}},
	idvpServerReference:        {CONNACK: {}, DISCONNECT: {}},
	idvpReasonString:           {CONNACK: {}, PUBACK: {}, PUBREC: {}, PUBREL: {}, PUBCOMP: {}, SUBACK: {}, UNSUBACK: {}, DISCONNECT: {}, AUTH: {}},
	idvpReceiveMaximum:         {CONNECT: {}, CONNACK: {}},
	idvpTopicAliasMaximum:      {CONNECT: {}, CONNACK: {}},
	idvpMaximumQOS:             {CONNECT: {}, CONNACK: {}},
	idvpMaximumPacketSize:      {CONNECT: {}, CONNACK: {}},
	idvpUserProperty:           {CONNECT: {}, CONNACK: {}, PUBLISH: {}, PUBACK: {}, PUBREC: {}, PUBREL: {}, PUBCOMP: {}, SUBACK: {}, UNSUBACK: {}, DISCONNECT: {}, AUTH: {}},
}

// ValidateID takes a PacketType and a property name and returns
// a boolean indicating if that property is valid for that
// PacketType
func ValidateID(p PacketType, i byte) bool {
	_, ok := ValidIDValuePairs[i][p]
	return ok
}
