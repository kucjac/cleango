// Code generated by "enumer -text -trimprefix=HealthCheckStatus -transform=snake_upper -type=HealthCheckStatus -output health_check_status_string.go"; DO NOT EDIT.

package xservice

import (
	"fmt"
	"strings"
)

const _HealthCheckStatusName = "UNKNOWNSERVINGNOT_SERVINGSERVICE_UNKNOWN"

var _HealthCheckStatusIndex = [...]uint8{0, 7, 14, 25, 40}

const _HealthCheckStatusLowerName = "unknownservingnot_servingservice_unknown"

func (i HealthCheckStatus) String() string {
	if i < 0 || i >= HealthCheckStatus(len(_HealthCheckStatusIndex)-1) {
		return fmt.Sprintf("HealthCheckStatus(%d)", i)
	}
	return _HealthCheckStatusName[_HealthCheckStatusIndex[i]:_HealthCheckStatusIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _HealthCheckStatusNoOp() {
	var x [1]struct{}
	_ = x[HealthCheckStatusUnknown-(0)]
	_ = x[HealthCheckStatusServing-(1)]
	_ = x[HealthCheckStatusNotServing-(2)]
	_ = x[HealthCheckStatusServiceUnknown-(3)]
}

var _HealthCheckStatusValues = []HealthCheckStatus{HealthCheckStatusUnknown, HealthCheckStatusServing, HealthCheckStatusNotServing, HealthCheckStatusServiceUnknown}

var _HealthCheckStatusNameToValueMap = map[string]HealthCheckStatus{
	_HealthCheckStatusName[0:7]:        HealthCheckStatusUnknown,
	_HealthCheckStatusLowerName[0:7]:   HealthCheckStatusUnknown,
	_HealthCheckStatusName[7:14]:       HealthCheckStatusServing,
	_HealthCheckStatusLowerName[7:14]:  HealthCheckStatusServing,
	_HealthCheckStatusName[14:25]:      HealthCheckStatusNotServing,
	_HealthCheckStatusLowerName[14:25]: HealthCheckStatusNotServing,
	_HealthCheckStatusName[25:40]:      HealthCheckStatusServiceUnknown,
	_HealthCheckStatusLowerName[25:40]: HealthCheckStatusServiceUnknown,
}

var _HealthCheckStatusNames = []string{
	_HealthCheckStatusName[0:7],
	_HealthCheckStatusName[7:14],
	_HealthCheckStatusName[14:25],
	_HealthCheckStatusName[25:40],
}

// HealthCheckStatusString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func HealthCheckStatusString(s string) (HealthCheckStatus, error) {
	if val, ok := _HealthCheckStatusNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _HealthCheckStatusNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to HealthCheckStatus values", s)
}

// HealthCheckStatusValues returns all values of the enum
func HealthCheckStatusValues() []HealthCheckStatus {
	return _HealthCheckStatusValues
}

// HealthCheckStatusStrings returns a slice of all String values of the enum
func HealthCheckStatusStrings() []string {
	strs := make([]string, len(_HealthCheckStatusNames))
	copy(strs, _HealthCheckStatusNames)
	return strs
}

// IsAHealthCheckStatus returns "true" if the value is listed in the enum definition. "false" otherwise
func (i HealthCheckStatus) IsAHealthCheckStatus() bool {
	for _, v := range _HealthCheckStatusValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalText implements the encoding.TextMarshaler interface for HealthCheckStatus
func (i HealthCheckStatus) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for HealthCheckStatus
func (i *HealthCheckStatus) UnmarshalText(text []byte) error {
	var err error
	*i, err = HealthCheckStatusString(string(text))
	return err
}