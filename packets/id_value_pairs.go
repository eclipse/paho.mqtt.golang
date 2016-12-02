package packets

type idvpID byte

const (
	idvpPayloadFormat         idvpID = 1
	idvpPubExpiry                    = 2
	idvpReplyTopic                   = 8
	idvpCorrelationData              = 9
	idvpSessionExpiryInterval        = 17
	idvpAssignedClientID             = 18
	idvpServerKeepAlive              = 19
	idvpAuthMethod                   = 21
	idvpAuthData                     = 22
	idvpRequestProblemInfo           = 23
	idvpWillDelayInterval            = 24
	idvpRequestReplyInfo             = 25
	idvpReplyInfo                    = 26
	idvpServerReference              = 28
	idvpReasonString                 = 31
	idvpReceiveMaximum               = 33
	idvpTopicAliasMaximum            = 34
	idvpTopicAlias                   = 35
	idvpMaximumQOS                   = 36
	idvpRetainUnavailable            = 37
	idvpUserDefinedPair              = 38
)

type payloadFormat byte
type pubExpiry uint32
type replyTopic string
type correlationData []byte
type sessionExpiryInterval uint32
type assignedClientID string
type serverKeepAlive uint16
type authMethod string
type authData []byte
type requestProblemInfo byte
type willDelayInterval uint32
type requestReplyInfo byte
type replyInfo string
type serverReference string
type reasonString string
type receiveMaximum uint16
type topicAliasMaximum uint16
type topicAlias uint16
type maximumQOS byte
type retainUnavailable bool
type userDefinedPair map[string]string

var validIDValuePairs = map[idvpID]map[packetType]struct{}{
	idvpPayloadFormat:         {PUBLISH: {}},
	idvpPubExpiry:             {PUBLISH: {}},
	idvpReplyTopic:            {PUBLISH: {}},
	idvpCorrelationData:       {PUBLISH: {}},
	idvpSessionExpiryInterval: {CONNECT: {}, DISCONNECT: {}},
	idvpAssignedClientID:      {CONNACK: {}},
	idvpServerKeepAlive:       {CONNACK: {}},
	idvpAuthMethod:            {CONNECT: {}, CONNACK: {}, AUTH: {}},
	idvpAuthData:              {CONNECT: {}, CONNACK: {}, AUTH: {}},
	idvpRequestProblemInfo:    {CONNECT: {}},
	idvpWillDelayInterval:     {CONNECT: {}},
	idvpRequestReplyInfo:      {CONNECT: {}},
	idvpReplyInfo:             {CONNACK: {}},
	idvpServerReference:       {CONNACK: {}, DISCONNECT: {}},
	idvpReasonString:          {CONNACK: {}, DISCONNECT: {}},
	idvpReceiveMaximum:        {CONNECT: {}, CONNACK: {}},
	idvpTopicAliasMaximum:     {CONNECT: {}, CONNACK: {}},
	idvpTopicAlias:            {PUBLISH: {}},
	idvpMaximumQOS:            {CONNECT: {}, CONNACK: {}},
	idvpRetainUnavailable:     {CONNACK: {}},
	idvpUserDefinedPair:       {CONNECT: {}, PUBLISH: {}},
}

func validateID(p packetType, i idvpID) bool {
	_, ok := validIDValuePairs[i][p]
	return ok
}
