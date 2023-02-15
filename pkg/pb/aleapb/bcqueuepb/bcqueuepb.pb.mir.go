package bcqueuepb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_InputValue)(nil)),
		reflect.TypeOf((*Event_Deliver)(nil)),
		reflect.TypeOf((*Event_FreeSlot)(nil)),
		reflect.TypeOf((*Event_SlotFreed)(nil)),
	}
}

func (*FreeSlotOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*FreeSlotOrigin_ContextStore)(nil)),
		reflect.TypeOf((*FreeSlotOrigin_Dsl)(nil)),
	}
}
