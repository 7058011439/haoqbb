package Net

import (
	"encoding/json"
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

const (
	keyRpcGo   = "rpcGo"
	keyRpcCall = "rpcCall"
	keyData    = "data"
	keyEvent   = "event"
)

type NatsEx struct {
	nc      *nats.Conn
	showLog bool
	speed
}

func NewNats(url string, showLog bool) *NatsEx {
	nc, err := nats.Connect(url)
	if err != nil {
		Log.Error("Failed to NewNats, url = %v, err = %v", url, err)
		return nil
	}
	ret := &NatsEx{nc: nc, showLog: showLog}
	return ret
}

func formatDataKey(dataType int) string {
	return fmt.Sprintf("%v_%v", keyData, dataType)
}

func formatEventKey(eventType int) string {
	return fmt.Sprintf("%v_%v", keyEvent, eventType)
}

func formatRpcGoKey(methodName string) string {
	return fmt.Sprintf("%v_%v", keyRpcGo, methodName)
}

func formatRpcCallKey(methodName string) string {
	return fmt.Sprintf("%v_%v", keyRpcCall, methodName)
}

func (n *NatsEx) handleRpc(mType *methodType, m *nats.Msg, needReply bool) {
	var argv reflect.Value
	var valueType bool
	if mType.argType.Kind() == reflect.Ptr {
		argv = reflect.New(mType.argType.Elem())
	} else {
		argv = reflect.New(mType.argType)
		valueType = true
	}
	reply := reflect.New(mType.replyType.Elem())
	err := json.Unmarshal(m.Data, argv.Interface())
	if err == nil {
		if valueType {
			argv = argv.Elem()
		}
		tmp := mType.method.Func.Call([]reflect.Value{mType.river, argv, reply})
		if !tmp[0].IsNil() {
			err = tmp[0].Interface().(error)
		}
	} else {
		Log.Error("Failed to handleRpc, argv err, err = %v", err)
	}

	if needReply {
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		replyData, _ := json.Marshal(reply.Elem().Interface())
		sendData, _ := json.Marshal(rpcMsgReply{
			Reply: replyData,
			Error: errMsg,
		})
		m.Respond(sendData)
	}
}

func (n *NatsEx) RegeditRpc(river interface{}) error {
	sName := reflect.Indirect(reflect.ValueOf(river)).Type().Name()
	method := suitableMethods(reflect.TypeOf(river), false)
	for name, mType := range method {
		mType.river = reflect.ValueOf(river)
		keyGo := formatRpcGoKey(sName + "." + name)
		n.nc.Subscribe(keyGo, func(m *nats.Msg) {
			n.RefreshMsgLen(len(m.Data))
			n.handleRpc(mType, m, false)
		})
		keyCall := formatRpcCallKey(sName + "." + name)
		n.nc.Subscribe(keyCall, func(m *nats.Msg) {
			n.RefreshMsgLen(len(m.Data))
			n.handleRpc(mType, m, true)
		})
	}

	return nil
}

// 每类型消息只能有一个注册接受
var mapMsg = make(map[string]bool)

func (n *NatsEx) RegeditRecvMsg(messageType int, fun func(int, []byte)) error {
	key := formatDataKey(messageType)
	if mapMsg[key] {
		return errors.Errorf("Failed to RegeditRecvMsg, messageType repeated")
	}
	n.nc.Subscribe(key, func(m *nats.Msg) {
		cost := Timer.NewTiming(Timer.Millisecond)
		n.RefreshMsgLen(len(m.Data))
		go fun(messageType, m.Data)
		cost.PrintCost(10, true, "Subscribe")
	})
	mapMsg[key] = true
	return nil
}

// 事件可以触发多个模块
func (n *NatsEx) RegeditEvent(eventType int, fn interface{}) error {
	value := reflect.ValueOf(fn)
	if value.Kind() != reflect.Func {
		return errors.New(fmt.Sprintf("Failed to RegeditEvent, param is not function"))
	}
	if value.Type().NumIn() != 2 {
		return errors.New(fmt.Sprintf("Failed to RegeditEvent, input param must two"))
	}
	argvType := value.Type().In(1)
	n.nc.Subscribe(formatEventKey(eventType), func(m *nats.Msg) {
		n.RefreshMsgLen(len(m.Data))
		var argv reflect.Value
		valueType := false
		if argvType.Kind() == reflect.Ptr {
			argv = reflect.New(argvType.Elem())
		} else {
			argv = reflect.New(argvType)
			valueType = true
		}
		if argvType.Kind() == reflect.Interface {
			argv = reflect.ValueOf(m.Data)
		} else {
			if err := json.Unmarshal(m.Data, argv.Interface()); err != nil {
				Log.Error("Failed to RegeditEvent, argv err, err = %v", err)
				return
			}
			if valueType {
				argv = argv.Elem()
			}
		}
		value.Call([]reflect.Value{reflect.ValueOf(eventType), argv})
	})
	return nil
}

func (n *NatsEx) RpcGo(moduleName string, data interface{}) {
	sendData, _ := json.Marshal(data)
	n.nc.Publish(formatRpcGoKey(moduleName), sendData)
}

func (n *NatsEx) RpcCall(moduleName string, data interface{}, reply interface{}) error {
	replyType := reflect.TypeOf(reply)
	if replyType.Kind() != reflect.Ptr {
		return errors.Errorf("Reply type of method %q is not a pointer", replyType)
	}
	sendData, _ := json.Marshal(data)
	rpcReply, err := n.nc.Request(formatRpcCallKey(moduleName), sendData, 1*time.Millisecond)
	if err != nil {
		return err
	}
	re := rpcMsgReply{}
	json.Unmarshal(rpcReply.Data, &re)
	err = json.Unmarshal(re.Reply, reply)
	if err != nil {
		return err
	}
	return errors.New(re.Error)
}

func (n *NatsEx) SendData(messageType int, data []byte) {
	cost := Timer.NewTiming(Timer.Millisecond)
	n.nc.Publish(formatDataKey(messageType), data)
	cost.PrintCost(10, false, "natsEx sendData")
}

func (n *NatsEx) PublishEvent(eventType int, data interface{}) {
	cost := Timer.NewTiming(Timer.Millisecond)
	sendData, _ := json.Marshal(data)
	n.nc.Publish(formatEventKey(eventType), sendData)
	cost.PrintCost(10, false, "natsEx PublishEvent")
}
