package Net

import (
	"fmt"
	"go/token"
	"reflect"
)

type methodType struct {
	method    reflect.Method
	argType   reflect.Type
	replyType reflect.Type
	river     reflect.Value // receiver of methods for the service
}

type rpcMsgReply struct {
	Error string
	Reply []byte
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return token.IsExported(t.Name()) || t.PkgPath() == ""
}

func suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs three ins: receiver, *args, *Reply.
		if mtype.NumIn() != 3 {
			if reportErr {
				fmt.Printf("rpc.RegisterService: method %q has %d input parameters; needs exactly three\n", mname, mtype.NumIn())
			}
			continue
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			if reportErr {
				fmt.Printf("rpc.RegisterService: argument type of method %q is not exported: %q\n", mname, argType)
			}
			continue
		}
		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			if reportErr {
				fmt.Printf("rpc.RegisterService: Reply type of method %q is not a pointer: %q\n", mname, replyType)
			}
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			if reportErr {
				fmt.Printf("rpc.RegisterService: Reply type of method %q is not exported: %q\n", mname, replyType)
			}
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			if reportErr {
				fmt.Printf("rpc.RegisterService: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
			}
			continue
		}
		if mtype.Out(0).String() != "error" {
			if reportErr {
				fmt.Printf("rpc.RegisterService: return type of method %v is %v, must be error\n", mname, typ.Out(0).String())
			}
		}
		methods[mname] = &methodType{method: method, argType: argType, replyType: replyType}
	}
	return methods
}
