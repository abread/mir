package vcbpb

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
		reflect.TypeOf((*Message_SendMessage)(nil)),
		reflect.TypeOf((*Message_EchoMessage)(nil)),
		reflect.TypeOf((*Message_FinalMessage)(nil)),
	}
}

func (*Origin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Origin_ContextStore)(nil)),
		reflect.TypeOf((*Origin_Dsl)(nil)),
		reflect.TypeOf((*Origin_AleaBcqueue)(nil)),
	}
}
