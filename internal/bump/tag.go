// Package bump parses, validates, and bumps strict semantic versions.
//
// Tags must match vMAJOR.MINOR.PATCH exactly: the v prefix is required
// and pre-release or build metadata is rejected.
package bump

// Tag returns the tag that follows tag after applying kind.
func Tag(tag string, kind Kind) (string, error) {
	version, err := Parse(tag)
	if err != nil {
		return "", err
	}
	if kind == nil {
		panic("bump kind must not be nil")
	}
	switch kind.(type) {
	case Major:
		version = Version{Major: version.Major + 1}
	case Minor:
		version = Version{Major: version.Major, Minor: version.Minor + 1}
	case Patch:
		version = Version{Major: version.Major, Minor: version.Minor, Patch: version.Patch + 1}
	}
	return version.String(), nil
}
