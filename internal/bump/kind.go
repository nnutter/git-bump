package bump

// Kind selects which version component to bump. It is sealed: the only
// variants are Major, Minor, and Patch, so type switches over Kind are
// exhaustiveness-checked by go-check-sumtype.
//
//sumtype:decl
type Kind interface {
	sealed()
}

// Major resets minor and patch to zero.
type Major struct{}

// Minor resets patch to zero.
type Minor struct{}

// Patch increments only the patch component.
type Patch struct{}

func (Major) sealed() {}
func (Minor) sealed() {}
func (Patch) sealed() {}

func (Major) String() string { return "major" }
func (Minor) String() string { return "minor" }
func (Patch) String() string { return "patch" }
