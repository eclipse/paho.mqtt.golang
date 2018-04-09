package packets

import "testing"

func TestPropertiess(t *testing.T) {
	if !ValidateID(PUBLISH, PropPayloadFormat) {
		t.Fatalf("'payloadFormat' is valid for 'PUBLISH' packets")
	}

	if !ValidateID(PUBLISH, PropPubExpiry) {
		t.Fatalf("'pubExpiry' is valid for 'PUBLISH' packets")
	}

	if !ValidateID(PUBLISH, PropReplyTopic) {
		t.Fatalf("'replyTopic' is valid for 'PUBLISH' packets")
	}

	if !ValidateID(PUBLISH, PropCorrelationData) {
		t.Fatalf("'correlationData' is valid for 'PUBLISH' packets")
	}

	if !ValidateID(CONNECT, PropSessionExpiryInterval) {
		t.Fatalf("'sessionExpiryInterval' is valid for 'CONNECT' packets")
	}

	if !ValidateID(DISCONNECT, PropSessionExpiryInterval) {
		t.Fatalf("'sessionExpiryInterval' is valid for 'DISCONNECT' packets")
	}

	if !ValidateID(CONNACK, PropAssignedClientID) {
		t.Fatalf("'assignedClientID' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNACK, PropServerKeepAlive) {
		t.Fatalf("'serverKeepAlive' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNECT, PropAuthMethod) {
		t.Fatalf("'authMethod' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, PropAuthMethod) {
		t.Fatalf("'authMethod' is valid for 'CONNACK' packets")
	}

	if !ValidateID(AUTH, PropAuthMethod) {
		t.Fatalf("'authMethod' is valid for 'auth' packets")
	}

	if !ValidateID(CONNECT, PropAuthData) {
		t.Fatalf("'authData' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, PropAuthData) {
		t.Fatalf("'authData' is valid for 'CONNACK' packets")
	}

	if !ValidateID(AUTH, PropAuthData) {
		t.Fatalf("'authData' is valid for 'auth' packets")
	}

	if !ValidateID(CONNECT, PropRequestProblemInfo) {
		t.Fatalf("'requestProblemInfo' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNECT, PropWillDelayInterval) {
		t.Fatalf("'willDelayInterval' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNECT, PropRequestResponseInfo) {
		t.Fatalf("'requestResponseInfo' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, PropResponseInfo) {
		t.Fatalf("'ResponseInfo' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNACK, PropServerReference) {
		t.Fatalf("'serverReference' is valid for 'CONNACK' packets")
	}

	if !ValidateID(DISCONNECT, PropServerReference) {
		t.Fatalf("'serverReference' is valid for 'DISCONNECT' packets")
	}

	if !ValidateID(CONNACK, PropReasonString) {
		t.Fatalf("'reasonString' is valid for 'CONNACK' packets")
	}

	if !ValidateID(DISCONNECT, PropReasonString) {
		t.Fatalf("'reasonString' is valid for 'DISCONNECT' packets")
	}

	if !ValidateID(CONNECT, PropReceiveMaximum) {
		t.Fatalf("'receiveMaximum' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, PropReceiveMaximum) {
		t.Fatalf("'receiveMaximum' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNECT, PropTopicAliasMaximum) {
		t.Fatalf("'topicAliasMaximum' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, PropTopicAliasMaximum) {
		t.Fatalf("'topicAliasMaximum' is valid for 'CONNACK' packets")
	}

	if !ValidateID(PUBLISH, PropTopicAlias) {
		t.Fatalf("'topicAlias' is valid for 'PUBLISH' packets")
	}

	if !ValidateID(CONNECT, PropMaximumQOS) {
		t.Fatalf("'maximumQOS' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, PropMaximumQOS) {
		t.Fatalf("'maximumQOS' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNACK, PropRetainAvailable) {
		t.Fatalf("'retainAvailable' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNECT, PropUser) {
		t.Fatalf("'user' is valid for 'CONNECT' packets")
	}

	if !ValidateID(PUBLISH, PropUser) {
		t.Fatalf("'user' is valid for 'PUBLISH' packets")
	}
}

func TestInvalidPropertiess(t *testing.T) {
	if ValidateID(PUBLISH, PropRequestResponseInfo) {
		t.Fatalf("'requestReplyInfo' is invalid for 'PUBLISH' packets")
	}
}
