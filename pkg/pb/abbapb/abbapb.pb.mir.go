package abbapb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_InputValue)(nil)),
		reflect.TypeOf((*Event_Deliver)(nil)),
		reflect.TypeOf((*Event_Round)(nil)),
	}
}

func (*Message) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Message_Finish)(nil)),
		reflect.TypeOf((*Message_Round)(nil)),
	}
}

func (*Origin) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Origin_ContextStore)(nil)),
		reflect.TypeOf((*Origin_Dsl)(nil)),
		reflect.TypeOf((*Origin_AleaAg)(nil)),
	}
}
