package bitmask

type Context struct {
	Network *Network
}

func NewContext() *Context {
	this := new(Context)
	this.Network = NewNetwork(this)

	return this
}
