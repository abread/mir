package modringpb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_Free)(nil)),
		reflect.TypeOf((*Event_Freed)(nil)),
		reflect.TypeOf((*Event_PastMessagesRecvd)(nil)),
	}
}

func (*FreeSubmoduleOrigin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*FreeSubmoduleOrigin_ContextStore)(nil)),
		reflect.TypeOf((*FreeSubmoduleOrigin_Dsl)(nil)),
	}
}
