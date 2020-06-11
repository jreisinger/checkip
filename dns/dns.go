package dns

import "net"

type DNS struct {
	Names []string
}

func New() *DNS {
	return &DNS{}
}

func (d *DNS) ForIP(addr net.IP) error {
	names, err := net.LookupAddr(addr.String())
	if err != nil {
		return err
	}
	d.Names = names
	return nil
}
