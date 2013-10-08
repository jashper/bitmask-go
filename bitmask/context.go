package bitmask

type Context struct {
	Network *Network

	AdPackets      Buffer
	PermReqPackets Buffer
}

func NewContext() (this *Context) {
	this = new(Context)
	this.Network = NewNetwork(this)

	this.AdPackets = NewBuffer(1000)
	this.PermReqPackets = NewBuffer(1000)

	return this
}
