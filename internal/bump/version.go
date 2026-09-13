package bump

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

var strictPattern = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

// Version is a parsed vMAJOR.MINOR.PATCH tag.
type Version struct {
	Major int
	Minor int
	Patch int
}

// Parse validates tag as a strict semantic version.
func Parse(tag string) (Version, error) {
	if !strictPattern.MatchString(tag) || !semver.IsValid(tag) {
		return Version{}, fmt.Errorf("invalid semver tag %q: must match vMAJOR.MINOR.PATCH", tag)
	}
	parts := strings.Split(tag[1:], ".")
	numbers := [3]int{}
	for i, part := range parts {
		number, err := strconv.Atoi(part)
		if err != nil {
			return Version{}, fmt.Errorf("invalid semver tag %q: %w", tag, err)
		}
		numbers[i] = number
	}
	return Version{Major: numbers[0], Minor: numbers[1], Patch: numbers[2]}, nil
}

// String renders the version with its v prefix.
func (v Version) String() string {
	return "v" + strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
}
