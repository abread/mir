package agreementpb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_RequestInput)(nil)),
		reflect.TypeOf((*Event_InputValue)(nil)),
		reflect.TypeOf((*Event_Deliver)(nil)),
	}
}

func (*Message) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Message_FinishAbba)(nil)),
	}
}
