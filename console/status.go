package console

import (
	"fmt"
)

type CmdStatus struct {
	rpcMethod	string
	rpcParams	string
	rpcResult	string
}

func (self *CmdStatus) usage(name string) string {
	return fmt.Sprintf("usage: %s status", name)
}

func (self *CmdStatus) defaults() error {
	self.rpcMethod = "Responder.Status"
	return nil
}

func( self *CmdStatus) idxArgsToFields() map[int]string {
	return nil
}

func (self *CmdStatus) FromArgs(args []string) error {
	self.defaults()
	return nil
}

func (self *CmdStatus) RpcMethod () string {
	return self.rpcMethod
}

func (self *CmdStatus) RpcParams() interface{} {
	return self.rpcParams
}

func (self *CmdStatus) RpcResult() interface{} {
	return &self.rpcResult
}
