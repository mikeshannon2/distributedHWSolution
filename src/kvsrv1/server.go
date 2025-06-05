package kvsrv

import (
	"log"
	"sync"

	"6.5840/kvsrv1/rpc"
	"6.5840/labrpc"
	"6.5840/raft"
	tester "6.5840/tester1"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type valueInfo struct {
	value   string
	version rpc.Tversion
}

type KVServer struct {
	mu sync.Mutex

	// Your definitions here.
	data map[string]valueInfo
}

func MakeKVServer() *KVServer {
	kv := &KVServer{data: make(map[string]valueInfo)}
	// Your code here.
	return kv
}

// Get returns the value and version for args.Key, if args.Key
// exists. Otherwise, Get returns ErrNoKey.
func (kv *KVServer) Get(args *rpc.GetArgs, reply *rpc.GetReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()
	v, ok := kv.data[args.Key]
	if !ok {
		reply.Err = rpc.ErrNoKey
	} else {
		reply.Err = rpc.OK
		reply.Value = v.value
		reply.Version = v.version
	}
}

// Update the value for a key if args.Version matches the version of
// the key on the server. If versions don't match, return ErrVersion.
// If the key doesn't exist, Put installs the value if the
// Args.Version is 0.
func (kv *KVServer) Put(args *rpc.PutArgs, reply *rpc.PutReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()

	key := args.Key
	version := args.Version
	v, ok := kv.data[key]
	if !ok && version != 0 {
		reply.Err = rpc.ErrNoKey
	} else if !ok && version == 0 {
		kv.data[key] = valueInfo{value: args.Value, version: version + 1}
		reply.Err = rpc.OK
	} else if ok && version != v.version {
		reply.Err = rpc.ErrVersion
	} else if ok && version == v.version {
		kv.data[key] = valueInfo{value: args.Value, version: version + 1}
	}
}

// You can ignore for this lab
func (kv *KVServer) Kill() {
}

// You can ignore for this lab
func (kv *KVServer) Raft() *raft.Raft {
	return nil
}

// You can ignore all arguments; they are for replicated KVservers in lab 4
func StartKVServer(ends []*labrpc.ClientEnd, gid tester.Tgid, srv int, persister *raft.Persister, maxraftstate int) tester.IKVServer {
	kv := MakeKVServer()
	return kv
}
