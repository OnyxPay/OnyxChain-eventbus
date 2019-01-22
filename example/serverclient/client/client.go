/****************************************************
Copyright 2019 The OnyxChain-eventbus Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*****************************************************/

/***************************************************
Copyright 2016 https://github.com/AsynkronIT/protoactor-go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*****************************************************/
package client

import (
	"fmt"
	"time"

	"github.com/OnyxPay/OnyxChain-eventbus/actor"
	"github.com/OnyxPay/OnyxChain-eventbus/example/serverclient/message"
)

type Client struct{}

//Call the server synchronous
func (client *Client) SyncCall(serverPID *actor.PID) (interface{}, error) {
	future := serverPID.RequestFuture(&message.Request{Who: "ONYXIO"}, 10*time.Second)
	result, err := future.Result()
	return result, err
}

//Call the server asynchronous
func (client *Client) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize client actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *message.Response:
		fmt.Println("Receive message", msg.Welcome)
	}
}

func (client *Client) AsyncCall(serverPID *actor.PID) *actor.PID {
	props := actor.FromProducer(func() actor.Actor { return &Client{} })
	clientPID := actor.Spawn(props)
	serverPID.Request(&message.Request{Who: "ONYXIO"}, clientPID)
	return clientPID
}

func (client *Client) Stop(pid *actor.PID) {
	pid.Stop()
}
