//go:generate stringer -type=Step,Arch,OS,SourceType
package nodester

import "fmt"

type Step int
type Arch int
type OS int
type SourceType int

const (
	Unpack Step = iota
	Download
	Compile
	Configure
	Build
	Install
)

const (
	Linux OS = iota + 100
	Darwin
	Windows
	Android
)

const (
	X86 Arch = iota + 200
	X64
	Arm
	Armbe
	Arm64
)

const (
	Git SourceType = iota + 300
	URL
)

func (self Step) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", self)), nil
}

func (self *Step) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "Unpack", "unpack":
		*self = Unpack
	}

	return nil
}
