package abbapb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_InputValue)(nil)),
		reflect.TypeOf((*Event_Deliver)(nil)),
	}
}

func (*Message) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Message_FinishMessage)(nil)),
		reflect.TypeOf((*Message_InitMessage)(nil)),
		reflect.TypeOf((*Message_AuxMessage)(nil)),
		reflect.TypeOf((*Message_ConfMessage)(nil)),
		reflect.TypeOf((*Message_CoinMessage)(nil)),
	}
}
