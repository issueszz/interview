package packageB

import "fmt"

type Pb struct {
	pa imlPa
}

type imlPa interface {
	TestA()
}

func NewPb(ipa imlPa) *Pb {
	return &Pb{pa: ipa}
}

func (pb *Pb) TestB() {
	fmt.Println("testB")
}

func (pb *Pb) PbMethod() {
	pb.pa.TestA()
}
