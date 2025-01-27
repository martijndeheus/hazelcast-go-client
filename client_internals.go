//go:build hazelcastinternal
// +build hazelcastinternal

/*
 * Copyright (c) 2008-2021, Hazelcast, Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hazelcast

import (
	"github.com/hazelcast/hazelcast-go-client/internal/cluster"
	"github.com/hazelcast/hazelcast-go-client/internal/event"
	"github.com/hazelcast/hazelcast-go-client/internal/invocation"
)

type ClientInternal struct {
	client *Client
}

func NewClientInternal(c *Client) *ClientInternal {
	return &ClientInternal{client: c}
}

func (ci *ClientInternal) ConnectionManager() *cluster.ConnectionManager {
	return ci.client.connectionManager
}

func (ci *ClientInternal) DispatchService() *event.DispatchService {
	return ci.client.eventDispatcher
}

func (ci *ClientInternal) InvocationService() *invocation.Service {
	return ci.client.invocationService
}

func (ci *ClientInternal) InvocationHandler() invocation.Handler {
	return ci.client.invocationHandler
}
