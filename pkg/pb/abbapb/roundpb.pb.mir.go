package abbapb

import (
	reflect "reflect"
)

func (*RoundEvent) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*RoundEvent_InputValue)(nil)),
		reflect.TypeOf((*RoundEvent_Deliver)(nil)),
		reflect.TypeOf((*RoundEvent_Finish)(nil)),
	}
}

func (*RoundMessage) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*RoundMessage_Init)(nil)),
		reflect.TypeOf((*RoundMessage_Aux)(nil)),
		reflect.TypeOf((*RoundMessage_Conf)(nil)),
		reflect.TypeOf((*RoundMessage_Coin)(nil)),
	}
}
