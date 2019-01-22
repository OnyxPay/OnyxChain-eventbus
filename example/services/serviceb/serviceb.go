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
package serviceb

import (
	"fmt"

	"github.com/OnyxPay/OnyxChain-eventbus/actor"
	message "github.com/OnyxPay/OnyxChain-eventbus/example/services/messages"
)

type ServiceB struct {
}

func (this *ServiceB) Receive(context actor.Context) {
	switch msg := context.Message().(type) {

	case *message.ServiceBRequest:
		fmt.Println("Receive ServiceBRequest:", msg.Message)
		context.Sender().Request(&message.ServiceBResponse{"response from serviceB"}, context.Self())

	case *message.ServiceAResponse:
		fmt.Println("Receive ServiceAResonse:", msg.Message)

	default:
		//fmt.Println("unknown message")
	}
}
