//go:build tools
// +build tools

// This file ensures all build tools are available and included in go.mod.

package tools

import (
	_ "golang.org/x/tools/cmd/stringer"
)
