// Code generated by "enumer -type=BranchAccess -yaml"; DO NOT EDIT.

package main

import (
	"fmt"
)

const (
	_BranchAccessName_0 = "NoAccess"
	_BranchAccessName_1 = "Developer"
	_BranchAccessName_2 = "Maintainer"
	_BranchAccessName_3 = "Admin"
)

var (
	_BranchAccessIndex_0 = [...]uint8{0, 8}
	_BranchAccessIndex_1 = [...]uint8{0, 9}
	_BranchAccessIndex_2 = [...]uint8{0, 10}
	_BranchAccessIndex_3 = [...]uint8{0, 5}
)

func (i BranchAccess) String() string {
	switch {
	case i == 0:
		return _BranchAccessName_0
	case i == 30:
		return _BranchAccessName_1
	case i == 40:
		return _BranchAccessName_2
	case i == 60:
		return _BranchAccessName_3
	default:
		return fmt.Sprintf("BranchAccess(%d)", i)
	}
}

var _BranchAccessValues = []BranchAccess{0, 30, 40, 60}

var _BranchAccessNameToValueMap = map[string]BranchAccess{
	_BranchAccessName_0[0:8]:  0,
	_BranchAccessName_1[0:9]:  30,
	_BranchAccessName_2[0:10]: 40,
	_BranchAccessName_3[0:5]:  60,
}

// BranchAccessString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func BranchAccessString(s string) (BranchAccess, error) {
	if val, ok := _BranchAccessNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to BranchAccess values", s)
}

// BranchAccessValues returns all values of the enum
func BranchAccessValues() []BranchAccess {
	return _BranchAccessValues
}

// IsABranchAccess returns "true" if the value is listed in the enum definition. "false" otherwise
func (i BranchAccess) IsABranchAccess() bool {
	for _, v := range _BranchAccessValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalYAML implements a YAML Marshaler for BranchAccess
func (i BranchAccess) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for BranchAccess
func (i *BranchAccess) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = BranchAccessString(s)
	return err
}
