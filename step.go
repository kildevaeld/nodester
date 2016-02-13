//go:generate stringer -type=Step
package nodester

type Step int

const (
	Unpack Step = iota
	Compile
	Configure
	Build
	Install
)
