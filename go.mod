module github.com/nnutter/git-bump

go 1.27.0

require (
	charm.land/fang/v2 v2.0.1
	charm.land/lipgloss/v2 v2.0.5
	github.com/charmbracelet/bubbles v0.21.1-0.20250623103423-23b8fd6302d7
	github.com/charmbracelet/bubbletea v1.3.6
	github.com/charmbracelet/huh v1.0.0
	github.com/charmbracelet/lipgloss v1.1.0
	github.com/samber/lo v1.53.0
	github.com/spf13/cobra v1.10.1
	github.com/stretchr/testify v1.12.1
)

tool (
	github.com/Antonboom/testifylint
	github.com/alecthomas/go-check-sumtype/cmd/go-check-sumtype
	github.com/evilmartians/lefthook/v2
	github.com/kisielk/errcheck
	github.com/nnutter/constable/cmd/constable
	github.com/securego/gosec/v2/cmd/gosec
	github.com/zricethezav/gitleaks/v8
	go.uber.org/nilaway/cmd/nilaway
	golang.org/x/tools/cmd/goimports
	golang.org/x/vuln/cmd/govulncheck
	honnef.co/go/tools/cmd/staticcheck
	mvdan.cc/gofumpt
)
