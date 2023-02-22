package bcpb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_StartBroadcast)(nil)),
		reflect.TypeOf((*Event_Deliver)(nil)),
		reflect.TypeOf((*Event_FreeSlot)(nil)),
		reflect.TypeOf((*Event_FillGap)(nil)),
	}
}
