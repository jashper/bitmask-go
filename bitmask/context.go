package bitmask

type Context struct {
	Version   Version
	Mode      Mode
	SubnetMap map[Subnet]bool

	Network        *Network
	AdPackets      *Buffer
	PermReqPackets *Buffer
}

func NewContext(v Version, m Mode, sList []Subnet) (this *Context) {
	this = new(Context)

	this.Version = v
	this.Mode = m

	this.SubnetMap = make(map[Subnet]bool)
	for _, s := range sList {
		this.SubnetMap[s] = true
	}

	this.Network = NewNetwork(this)
	this.AdPackets = NewBuffer(1000)
	this.PermReqPackets = NewBuffer(1000)

	return this
}
