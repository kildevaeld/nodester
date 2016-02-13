// Code generated by "stringer -type=Step,Arch,OS,SourceType"; DO NOT EDIT

package nodester

import "fmt"

const _Step_name = "UnpackDownloadCompileConfigureBuildInstall"

var _Step_index = [...]uint8{0, 6, 14, 21, 30, 35, 42}

func (i Step) String() string {
	if i < 0 || i >= Step(len(_Step_index)-1) {
		return fmt.Sprintf("Step(%d)", i)
	}
	return _Step_name[_Step_index[i]:_Step_index[i+1]]
}

const _Arch_name = "X86X64ArmArmbeArm64"

var _Arch_index = [...]uint8{0, 3, 6, 9, 14, 19}

func (i Arch) String() string {
	i -= 200
	if i < 0 || i >= Arch(len(_Arch_index)-1) {
		return fmt.Sprintf("Arch(%d)", i+200)
	}
	return _Arch_name[_Arch_index[i]:_Arch_index[i+1]]
}

const _OS_name = "LinuxDarwinWindowsAndroid"

var _OS_index = [...]uint8{0, 5, 11, 18, 25}

func (i OS) String() string {
	i -= 100
	if i < 0 || i >= OS(len(_OS_index)-1) {
		return fmt.Sprintf("OS(%d)", i+100)
	}
	return _OS_name[_OS_index[i]:_OS_index[i+1]]
}

const _SourceType_name = "GitURL"

var _SourceType_index = [...]uint8{0, 3, 6}

func (i SourceType) String() string {
	i -= 300
	if i < 0 || i >= SourceType(len(_SourceType_index)-1) {
		return fmt.Sprintf("SourceType(%d)", i+300)
	}
	return _SourceType_name[_SourceType_index[i]:_SourceType_index[i+1]]
}
