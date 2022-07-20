package rpcxr

import "github.com/vmihailenco/msgpack/v5"

func MSGPEncode(m interface{}) []byte {
    if data, err := msgpack.Marshal(m); err == nil {
        return data
    } else {
        panic(err)
    }
    return nil
}
func MSGPDecode(data []byte, m interface{}) bool {
    err := msgpack.Unmarshal(data, m)
    if err != nil {
        panic(err)
    }
    return true
}
