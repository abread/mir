package reliablenetpb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_SendMessage)(nil)),
		reflect.TypeOf((*Event_Ack)(nil)),
		reflect.TypeOf((*Event_MarkRecvd)(nil)),
		reflect.TypeOf((*Event_MarkModuleMsgsRecvd)(nil)),
		reflect.TypeOf((*Event_RetransmitAll)(nil)),
	}
}
