package matcher

import (
	"fmt"
	"regexp"

	"github.com/golang/mock/gomock"
)

type guid struct{}

// IsGUID checks whether the string is a valid gUID
func IsGUID() gomock.Matcher {
	return &guid{}
}

func (o *guid) Matches(x interface{}) bool {
	str := fmt.Sprintf("%v", x)
	r, _ := regexp.Compile(`([0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12})$`)
	return r.MatchString(str)
}

func (o *guid) String() string {
	return "is a GUID"
}
