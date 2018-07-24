package sdk

import (
	protorpc "code.google.com/p/protorpc"
	"gopkg.in/fatih/pool.v2"
	"net/rpc"
)

// RPCClientPool Used to maintain a connnected instances
type RPCClientPool struct {
	impl pool.Pool
}

// NewRPCClientPool create a new RPCClietnPool instance
func NewRPCClientPool(initialCount, maxCount int, factory pool.Factory) (*RPCClientPool, error) {
	ret, err := pool.NewChannelPool(initialCount, maxCount, factory)
	if err != nil {
		return nil, err
	}
	return &RPCClientPool{
		impl: ret,
	}, nil
}

func (obj *RPCClientPool) Get() (*rpc.Client, error) {
	conn, err := obj.impl.Get()

	if err != nil {
		return nil, err
	}
	return protorpc.NewClient(conn), nil
}

func (obj *RPCClientPool) Len() int { return obj.impl.Len() }

func (obj *RPCClientPool) Close() { obj.impl.Close() }
