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
package remote

import (
	"github.com/OnyxPay/OnyxChain-eventbus/actor"
	"github.com/OnyxPay/OnyxChain-eventbus/eventhub"
	"github.com/OnyxPay/OnyxChain-eventbus/log"
)

func newEndpointWatcher(address string) actor.Producer {
	return func() actor.Actor {
		return &endpointWatcher{
			address: address,
		}
	}
}

type endpointWatcher struct {
	address string
	watched map[string]*actor.PIDSet //key is the watching PID string, value is the watched PID
}

func (state *endpointWatcher) initialize() {
	plog.Info("Started EndpointWatcher", log.String("address", state.address))
	state.watched = make(map[string]*actor.PIDSet)
}

func (state *endpointWatcher) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.initialize()

	case *remoteTerminate:
		//delete the watch entries
		if pidSet, ok := state.watched[msg.Watcher.Id]; ok {
			pidSet.Remove(msg.Watchee)
			if pidSet.Len() == 0 {
				delete(state.watched, msg.Watcher.Id)
			}
		}

		terminated := &actor.Terminated{
			Who:               msg.Watchee,
			AddressTerminated: false,
		}
		ref, ok := actor.ProcessRegistry.GetLocal(msg.Watcher.Id)
		if ok {
			ref.SendSystemMessage(msg.Watcher, terminated)
		}
	case *EndpointConnectedEvent:
		//Already connected, pass
	case *EndpointTerminatedEvent:
		plog.Info("EndpointWatcher handling terminated", log.String("address", state.address))

		for id, pidSet := range state.watched {
			//try to find the watcher ID in the local actor registry
			ref, ok := actor.ProcessRegistry.GetLocal(id)
			if ok {
				pidSet.ForEach(func(i int, pid actor.PID) {
					eventhub.GlobalEventHub.RemovePID(pid)
					//create a terminated event for the Watched actor
					terminated := &actor.Terminated{
						Who:               &pid,
						AddressTerminated: true,
					}

					watcher := actor.NewLocalPID(id)
					//send the address Terminated event to the Watcher
					ref.SendSystemMessage(watcher, terminated)
				})
			}
		}

		//Clear watcher's map
		state.watched = make(map[string]*actor.PIDSet)
		ctx.SetBehavior(state.Terminated)

	case *remoteWatch:
		//add watchee to watcher's map
		if pidSet, ok := state.watched[msg.Watcher.Id]; ok {
			pidSet.Add(msg.Watchee)
		} else {
			state.watched[msg.Watcher.Id] = actor.NewPIDSet(msg.Watchee)
		}

		//recreate the Watch command
		w := &actor.Watch{
			Watcher: msg.Watcher,
		}

		//pass it off to the remote PID
		SendMessage(msg.Watchee, nil, w, nil, -1)

	case *remoteUnwatch:
		//delete the watch entries
		if pidSet, ok := state.watched[msg.Watcher.Id]; ok {
			pidSet.Remove(msg.Watchee)
			if pidSet.Len() == 0 {
				delete(state.watched, msg.Watcher.Id)
			}
		}

		//recreate the Unwatch command
		uw := &actor.Unwatch{
			Watcher: msg.Watcher,
		}

		//pass it off to the remote PID
		SendMessage(msg.Watchee, nil, uw, nil, -1)
	case actor.SystemMessage, actor.AutoReceiveMessage:
		//ignore
	default:
		plog.Error("EndpointWatcher received unknown message", log.String("address", state.address), log.Message(msg))
	}
}

func (state *endpointWatcher) Terminated(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *remoteWatch:
		//try to find the watcher ID in the local actor registry
		ref, ok := actor.ProcessRegistry.GetLocal(msg.Watcher.Id)
		if ok {

			//create a terminated event for the Watched actor
			terminated := &actor.Terminated{
				Who:               msg.Watchee,
				AddressTerminated: true,
			}
			//send the address Terminated event to the Watcher
			ref.SendSystemMessage(msg.Watcher, terminated)
		}
	case *EndpointConnectedEvent:
		plog.Info("EndpointWatcher handling restart", log.String("address", state.address))
		ctx.SetBehavior(state.Receive)
	case *remoteTerminate, *EndpointTerminatedEvent, *remoteUnwatch:
		// pass
		plog.Error("EndpointWatcher receive message for already terminated endpoint", log.String("address", state.address), log.Message(msg))
	case actor.SystemMessage, actor.AutoReceiveMessage:
		//ignore
	default:
		plog.Error("EndpointWatcher received unknown message", log.String("address", state.address), log.TypeOf("type", msg), log.Message(msg))
	}
}
