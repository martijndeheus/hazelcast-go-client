// Copyright (c) 2008-2017, Hazelcast, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package codec

import (
	. "github.com/hazelcast/go-client"
)

type ClientGetPartitionsResponseParameters struct {
	Partitions []Pair
}

func (codec *ClientGetPartitionsResponseParameters) calculateSize() int {
	// Calculates the request payload size
	dataSize := 0
	return dataSize
}

func (codec *ClientGetPartitionsResponseParameters) encodeRequest() *ClientMessage {
	// Encode request into clientMessage
	clientMessage := NewClientMessage(nil, codec.calculateSize())
	clientMessage.SetMessageType(CLIENT_GETPARTITIONS)
	clientMessage.IsRetryable = false
	clientMessage.UpdateFrameLength()
	return clientMessage
}

func (codec *ClientGetPartitionsResponseParameters) decodeResponse(clientMessage *ClientMessage) *ClientGetPartitionsResponseParameters {
	// Decode response from client message
	parameters := new(ClientGetPartitionsResponseParameters)

	partitionsSize := clientMessage.ReadInt32()
	partitions := make([]Pair, partitionsSize)
	for partitionsIndex := 0; partitionsIndex < int(partitionsSize); partitionsIndex++ {
		var partitionsItem Pair
		partitionsItemKey := *AddressCodecDecode(clientMessage)

		partitionsItemValSize := clientMessage.ReadInt32()
		partitionsItemVal := make([]int32, partitionsItemValSize)
		for partitionsItemValIndex := 0; partitionsItemValIndex < int(partitionsItemValSize); partitionsItemValIndex++ {
			partitionsItemValItem := clientMessage.ReadInt32()
			partitionsItemVal = append(partitionsItemVal, partitionsItemValItem)
		}

		partitionsItem.Key = partitionsItemKey
		partitionsItem.Value = partitionsItemVal
		partitions = append(partitions, partitionsItem)
	}
	parameters.Partitions = partitions

	return parameters
}
