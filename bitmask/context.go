package bitmask

type Context struct {
	Network *Network
}

func NewContext() (this *Context) {
	this = new(Context)
	this.Network = NewNetwork(this)

	return this
}
