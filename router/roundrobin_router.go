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
package router

import (
	"sync/atomic"

	"github.com/OnyxPay/OnyxChain-eventbus/actor"
)

type roundRobinGroupRouter struct {
	GroupRouter
}

type roundRobinPoolRouter struct {
	PoolRouter
}

type roundRobinState struct {
	index   int32
	routees *actor.PIDSet
	values  []actor.PID
}

func (state *roundRobinState) SetRoutees(routees *actor.PIDSet) {
	state.routees = routees
	state.values = routees.Values()
}

func (state *roundRobinState) GetRoutees() *actor.PIDSet {
	return state.routees
}

func (state *roundRobinState) RouteMessage(message interface{}) {
	pid := roundRobinRoutee(&state.index, state.values)
	pid.Tell(message)
}

func NewRoundRobinPool(size int) *actor.Props {
	return actor.FromSpawnFunc(spawner(&roundRobinPoolRouter{PoolRouter{PoolSize: size}}))
}

func NewRoundRobinGroup(routees ...*actor.PID) *actor.Props {
	return actor.FromSpawnFunc(spawner(&roundRobinGroupRouter{GroupRouter{Routees: actor.NewPIDSet(routees...)}}))
}

func (config *roundRobinPoolRouter) CreateRouterState() Interface {
	return &roundRobinState{}
}

func (config *roundRobinGroupRouter) CreateRouterState() Interface {
	return &roundRobinState{}
}

func roundRobinRoutee(index *int32, routees []actor.PID) actor.PID {
	i := int(atomic.AddInt32(index, 1))
	if i < 0 {
		*index = 0
		i = 0
	}
	mod := len(routees)
	routee := routees[i%mod]
	return routee
}
