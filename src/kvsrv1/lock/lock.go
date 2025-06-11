package lock

import (
	"time"

	"6.5840/kvsrv1/rpc"
	kvtest "6.5840/kvtest1"
)

type Lock struct {
	// IKVClerk is a go interface for k/v clerks: the interfaces hides
	// the specific Clerk type of ck but promises that ck supports
	// Put and Get.  The tester passes the clerk in when calling
	// MakeLock().
	ck kvtest.IKVClerk
	// You may add code here
	lockKey string
	lockId  string
}

// The tester calls MakeLock() and passes in a k/v clerk; you code can
// perform a Put or Get by calling lk.ck.Put() or lk.ck.Get().
func MakeLock(ck kvtest.IKVClerk, l string) *Lock {
	lk := &Lock{ck: ck, lockKey: l, lockId: kvtest.RandValue(8)}
	// You may add code here
	return lk
}

func (lk *Lock) Acquire() {
	for {
		value, version, err := lk.ck.Get(lk.lockKey)
		var putErr rpc.Err = rpc.ErrVersion
		if err == rpc.ErrNoKey {
			putErr = lk.ck.Put(lk.lockKey, lk.lockId, 0)
		} else if value == "0" {
			putErr = lk.ck.Put(lk.lockKey, lk.lockId, version)
		} else if value == lk.lockId {
			break
		}

		if putErr == rpc.ErrVersion || putErr == rpc.ErrMaybe {
			time.Sleep(50 * time.Millisecond)
		} else {
			break
		}

	}
}

func (lk *Lock) Release() {
	// Your code here
	value, version, _ := lk.ck.Get(lk.lockKey)
	if value == lk.lockId {
		lk.ck.Put(lk.lockKey, "0", version)
	}
}
