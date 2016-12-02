package packets

import "testing"

func TestIDValuePairs(t *testing.T) {
	if !validateID(publish, idvpPayloadFormat) {
		t.Fatalf("'payloadFormat' is valid for 'publish' packets")
	}

	if !validateID(publish, idvpPubExpiry) {
		t.Fatalf("'pubExpiry' is valid for 'publish' packets")
	}

	if !validateID(publish, idvpReplyTopic) {
		t.Fatalf("'replyTopic' is valid for 'publish' packets")
	}

	if !validateID(publish, idvpCorrelationData) {
		t.Fatalf("'correlationData' is valid for 'publish' packets")
	}

	if !validateID(connect, idvpSessionExpiryInterval) {
		t.Fatalf("'sessionExpiryInterval' is valid for 'connect' packets")
	}

	if !validateID(disconnect, idvpSessionExpiryInterval) {
		t.Fatalf("'sessionExpiryInterval' is valid for 'disconnect' packets")
	}

	if !validateID(connack, idvpAssignedClientID) {
		t.Fatalf("'assignedClientID' is valid for 'connack' packets")
	}

	if !validateID(connack, idvpServerKeepAlive) {
		t.Fatalf("'serverKeepAlive' is valid for 'connack' packets")
	}

	if !validateID(connect, idvpAuthMethod) {
		t.Fatalf("'authMethod' is valid for 'connect' packets")
	}

	if !validateID(connack, idvpAuthMethod) {
		t.Fatalf("'authMethod' is valid for 'connack' packets")
	}

	if !validateID(auth, idvpAuthMethod) {
		t.Fatalf("'authMethod' is valid for 'auth' packets")
	}

	if !validateID(connect, idvpAuthData) {
		t.Fatalf("'authData' is valid for 'connect' packets")
	}

	if !validateID(connack, idvpAuthData) {
		t.Fatalf("'authData' is valid for 'connack' packets")
	}

	if !validateID(auth, idvpAuthData) {
		t.Fatalf("'authData' is valid for 'auth' packets")
	}

	if !validateID(connect, idvpRequestProblemInfo) {
		t.Fatalf("'requestProblemInfo' is valid for 'connect' packets")
	}

	if !validateID(connect, idvpWillDelayInterval) {
		t.Fatalf("'willDelayInterval' is valid for 'connect' packets")
	}

	if !validateID(connect, idvpRequestReplyInfo) {
		t.Fatalf("'requestReplyInfo' is valid for 'connect' packets")
	}

	if !validateID(connack, idvpReplyInfo) {
		t.Fatalf("'replyInfo' is valid for 'connack' packets")
	}

	if !validateID(connack, idvpServerReference) {
		t.Fatalf("'serverReference' is valid for 'connack' packets")
	}

	if !validateID(disconnect, idvpServerReference) {
		t.Fatalf("'serverReference' is valid for 'disconnect' packets")
	}

	if !validateID(connack, idvpReasonString) {
		t.Fatalf("'reasonString' is valid for 'connack' packets")
	}

	if !validateID(disconnect, idvpReasonString) {
		t.Fatalf("'reasonString' is valid for 'disconnect' packets")
	}

	if !validateID(connect, idvpReceiveMaximum) {
		t.Fatalf("'receiveMaximum' is valid for 'connect' packets")
	}

	if !validateID(connack, idvpReceiveMaximum) {
		t.Fatalf("'receiveMaximum' is valid for 'connack' packets")
	}

	if !validateID(connect, idvpTopicAliasMaximum) {
		t.Fatalf("'topicAliasMaximum' is valid for 'connect' packets")
	}

	if !validateID(connack, idvpTopicAliasMaximum) {
		t.Fatalf("'topicAliasMaximum' is valid for 'connack' packets")
	}

	if !validateID(publish, idvpTopicAlias) {
		t.Fatalf("'topicAlias' is valid for 'publish' packets")
	}

	if !validateID(connect, idvpMaximumQOS) {
		t.Fatalf("'maximumQOS' is valid for 'connect' packets")
	}

	if !validateID(connack, idvpMaximumQOS) {
		t.Fatalf("'maximumQOS' is valid for 'connack' packets")
	}

	if !validateID(connack, idvpRetainUnavailable) {
		t.Fatalf("'retainUnavailable' is valid for 'connack' packets")
	}

	if !validateID(connect, idvpUserDefinedPair) {
		t.Fatalf("'userDefinedPair' is valid for 'connect' packets")
	}

	if !validateID(publish, idvpUserDefinedPair) {
		t.Fatalf("'userDefinedPair' is valid for 'publish' packets")
	}
}

func TestInvalidIDValuePairs(t *testing.T) {
	if validateID(publish, idvpRequestReplyInfo) {
		t.Fatalf("'requestReplyInfo' is invalid for 'publish' packets")
	}
}
