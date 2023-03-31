package bcqueuepb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_InputValue)(nil)),
		reflect.TypeOf((*Event_Deliver)(nil)),
		reflect.TypeOf((*Event_FreeSlot)(nil)),
		reflect.TypeOf((*Event_PastVcbFinal)(nil)),
		reflect.TypeOf((*Event_BcStarted)(nil)),
	}
}
