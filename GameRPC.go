package rpcxr

import (
    "bytes"
    "context"
    "errors"
    "fmt"
    "github.com/smallnest/rpcx/server"
    "log"
    "reflect"
    "runtime"
    "sync"
    "sync/atomic"
)

var (
    StatRecvBytes    uint64
    StatSendBytes    uint64
    StatRecvNum      uint64
    StatSendNum      uint64
    StatErrNum       uint64
    isStop           bool
    runningCount     int32
    runningWaitGroup sync.WaitGroup
)

func makeRPCFunc(path, name string, fn interface{}, r interface{}, w interface{}) (func(ctx context.Context, args *RPCProxyRequest, reply *RPCProxyReply) error, error) {
    // 检查传入的函数是否符合格式要求
    var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
    f, ok := fn.(reflect.Value)
    if !ok {
        f = reflect.ValueOf(fn)
    }
    t := f.Type()
    if t.NumIn() != 3 { // context/request/response
        return nil, fmt.Errorf("func in param num error")
    }
    if t.NumOut() != 1 {
        return nil, fmt.Errorf("func out param num error")
    }
    if returnType := t.Out(0); returnType != typeOfError {
        return nil, fmt.Errorf("func out param type error")
    }
    // 生成 rpcx server 识别格式的函数
    return func(ctx context.Context, args *RPCProxyRequest, reply *RPCProxyReply) (err error) {
        if isStop {
            return errors.New("service stopped")
        }
        // defer nobug.ExecTime(fmt.Sprintf("RPC Duration: [%v.%v]", path, name), time.Now())
        atomic.AddInt32(&runningCount, 1)
        runningWaitGroup.Add(1)
        defer func() {
            atomic.AddInt32(&runningCount, -1)
            runningWaitGroup.Done()
            if p := recover(); p != nil {
                buffer := bytes.NewBufferString(fmt.Sprintf("%v\n", p))
                // 打印调用栈信息
                buf := make([]byte, 4096)
                n := runtime.Stack(buf, false)
                stackInfo := fmt.Sprintf("%s", buf[:n])
                buffer.WriteString(fmt.Sprintf("panic stack info %s", stackInfo))
                log.Printf("RPC[%s] [%s] 请求失败:%v\n", path, name, buffer)
                atomic.AddUint64(&StatErrNum, 1)
            }
        }()
        atomic.AddUint64(&StatRecvNum, 1)
        atomic.AddUint64(&StatRecvBytes, uint64(len(args.Data)))
        ctx = context.WithValue(ctx, "playerId", args.PlayerId)
        ctx = context.WithValue(ctx, "ext", args.Ext)

        t0 := reflect.ValueOf(ctx)
        t1 := reflect.New(reflect.TypeOf(r))
        t2 := reflect.New(reflect.TypeOf(w))

        MSGPDecode(args.Data, t1.Interface())
        in := []reflect.Value{
            t0, t1, t2,
        }
        rs := f.Call(in)
        r1 := rs[0]
        reply.Data = MSGPEncode(t2.Interface())
        if r1.Interface() == nil {

            atomic.AddUint64(&StatSendNum, 1)
            atomic.AddUint64(&StatSendBytes, uint64(len(reply.Data)))

            return nil
        } else {
            err = r1.Interface().(error)
            atomic.AddUint64(&StatErrNum, 1)
            return err
        }
    }, nil
}

func ExportRPC(i interface{}, path, name string, fn, r, w interface{}) {
    s := i.(*server.Server)
    f, err := makeRPCFunc(path, name, fn, r, w)
    if err != nil {
        log.Printf("ExportRPC:%v\n", err)
        return
    }
    _ = s.RegisterFunctionName(path, name, f, "")
}

func makeRPCFunc2(path, name string, fn interface{}) (func(ctx context.Context, args *RPCProxyRequest, reply *RPCProxyReply) error, error) {
    // 检查传入的函数是否符合格式要求
    var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
    f, ok := fn.(reflect.Value)
    if !ok {
        f = reflect.ValueOf(fn)
    }
    t := f.Type()
    if t.NumIn() != 3 { // context/request/response
        return nil, fmt.Errorf("func in param num error")
    }
    if t.NumOut() != 1 {
        return nil, fmt.Errorf("func out param num error")
    }
    if returnType := t.Out(0); returnType != typeOfError {
        return nil, fmt.Errorf("func out param type error")
    }

    s1 := reflect.New(t.In(1).Elem()).Elem()
    s2 := reflect.New(t.In(2).Elem()).Elem()
    var (
        r = s1.Interface()
        w = s2.Interface()
    )

    return makeRPCFunc(path, name, fn, r, w)
}
func ExportRPC2(i interface{}, path, name string, fn interface{}) {
    s := i.(*server.Server)
    f, err := makeRPCFunc2(path, name, fn)
    if err != nil {
        // logger.Errorf("ExportRPC:%v", err)
        log.Printf("ExportRPC:%v\n", err)
        return
    }
    _ = s.RegisterFunctionName(path, name, f, "")
}

func Status() string {
    return fmt.Sprintf(`
rpc status:
    接收流量:%v
    发送流量:%v
    接收请求:%v
    发送请求:%v
    错误请求:%v
    运行中:%v
`,
        StatRecvBytes,
        StatSendBytes,
        StatRecvNum,
        StatSendNum,
        StatErrNum,
        runningCount)
}
func WaitAllFinish() {
    isStop = true
    runningWaitGroup.Wait()
}

func CallRPC(cli Client, playerId string, routeName string, r interface{}, w interface{}) error {
    var (
        args  = &RPCProxyRequest{}
        reply = &RPCProxyReply{}
    )
    args.PlayerId = playerId
    args.Data = MSGPEncode(r)
    err := cli.Call(context.Background(), routeName, args, reply)
    if err != nil {
        return err
    }
    MSGPDecode(reply.Data, &w)
    return nil
}
