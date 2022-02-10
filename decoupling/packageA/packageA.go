package packageA

import "fmt"

type Pa struct {
	pb imlPb
}

type imlPb interface {
	TestB()
}

func NewPa(ipb imlPb) *Pa {
	return &Pa{
		pb: ipb,
	}
}
func (pa *Pa) TestA() {
	fmt.Println("testA")
}

func (pa *Pa) PaMethod() {
	pa.pb.TestB()
}
