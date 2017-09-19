package packets

import "testing"

func TestIDValuePairs(t *testing.T) {
	if !ValidateID(PUBLISH, idvpPayloadFormat) {
		t.Fatalf("'payloadFormat' is valid for 'PUBLISH' packets")
	}

	if !ValidateID(PUBLISH, idvpPubExpiry) {
		t.Fatalf("'pubExpiry' is valid for 'PUBLISH' packets")
	}

	if !ValidateID(PUBLISH, idvpReplyTopic) {
		t.Fatalf("'replyTopic' is valid for 'PUBLISH' packets")
	}

	if !ValidateID(PUBLISH, idvpCorrelationData) {
		t.Fatalf("'correlationData' is valid for 'PUBLISH' packets")
	}

	if !ValidateID(CONNECT, idvpSessionExpiryInterval) {
		t.Fatalf("'sessionExpiryInterval' is valid for 'CONNECT' packets")
	}

	if !ValidateID(DISCONNECT, idvpSessionExpiryInterval) {
		t.Fatalf("'sessionExpiryInterval' is valid for 'DISCONNECT' packets")
	}

	if !ValidateID(CONNACK, idvpAssignedClientID) {
		t.Fatalf("'assignedClientID' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNACK, idvpServerKeepAlive) {
		t.Fatalf("'serverKeepAlive' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNECT, idvpAuthMethod) {
		t.Fatalf("'authMethod' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, idvpAuthMethod) {
		t.Fatalf("'authMethod' is valid for 'CONNACK' packets")
	}

	if !ValidateID(AUTH, idvpAuthMethod) {
		t.Fatalf("'authMethod' is valid for 'auth' packets")
	}

	if !ValidateID(CONNECT, idvpAuthData) {
		t.Fatalf("'authData' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, idvpAuthData) {
		t.Fatalf("'authData' is valid for 'CONNACK' packets")
	}

	if !ValidateID(AUTH, idvpAuthData) {
		t.Fatalf("'authData' is valid for 'auth' packets")
	}

	if !ValidateID(CONNECT, idvpRequestProblemInfo) {
		t.Fatalf("'requestProblemInfo' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNECT, idvpWillDelayInterval) {
		t.Fatalf("'willDelayInterval' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNECT, idvpRequestResponseInfo) {
		t.Fatalf("'requestResponseInfo' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, idvpResponseInfo) {
		t.Fatalf("'ResponseInfo' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNACK, idvpServerReference) {
		t.Fatalf("'serverReference' is valid for 'CONNACK' packets")
	}

	if !ValidateID(DISCONNECT, idvpServerReference) {
		t.Fatalf("'serverReference' is valid for 'DISCONNECT' packets")
	}

	if !ValidateID(CONNACK, idvpReasonString) {
		t.Fatalf("'reasonString' is valid for 'CONNACK' packets")
	}

	if !ValidateID(DISCONNECT, idvpReasonString) {
		t.Fatalf("'reasonString' is valid for 'DISCONNECT' packets")
	}

	if !ValidateID(CONNECT, idvpReceiveMaximum) {
		t.Fatalf("'receiveMaximum' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, idvpReceiveMaximum) {
		t.Fatalf("'receiveMaximum' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNECT, idvpTopicAliasMaximum) {
		t.Fatalf("'topicAliasMaximum' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, idvpTopicAliasMaximum) {
		t.Fatalf("'topicAliasMaximum' is valid for 'CONNACK' packets")
	}

	if !ValidateID(PUBLISH, idvpTopicAlias) {
		t.Fatalf("'topicAlias' is valid for 'PUBLISH' packets")
	}

	if !ValidateID(CONNECT, idvpMaximumQOS) {
		t.Fatalf("'maximumQOS' is valid for 'CONNECT' packets")
	}

	if !ValidateID(CONNACK, idvpMaximumQOS) {
		t.Fatalf("'maximumQOS' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNACK, idvpRetainAvailable) {
		t.Fatalf("'retainAvailable' is valid for 'CONNACK' packets")
	}

	if !ValidateID(CONNECT, idvpUserProperty) {
		t.Fatalf("'userProperty' is valid for 'CONNECT' packets")
	}

	if !ValidateID(PUBLISH, idvpUserProperty) {
		t.Fatalf("'userProperty' is valid for 'PUBLISH' packets")
	}
}

func TestInvalidIDValuePairs(t *testing.T) {
	if ValidateID(PUBLISH, idvpRequestResponseInfo) {
		t.Fatalf("'requestReplyInfo' is invalid for 'PUBLISH' packets")
	}
}
