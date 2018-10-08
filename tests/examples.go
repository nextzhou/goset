//go:generate goderive
package tests

import (
	"net/http"
	t "time"

	"github.com/nextzhou/goderive/plugin"
)

// derive-set
type Int = int

// derive-set:Rename=intOrderSet;Order=Append
type Int2 = int

// derive-set:Order=Key
type Int3 = int

// unexported type, from imported package
// derive-set
type h = http.Handler

// from renamed imported package
// derive-set
type T = t.Time

// from this package
// derive-set: !Export
type A struct{ s string }

// derive-set:Order=Key
type S = string

// from remote package
// derive-set: Export
type p = plugin.Plugin
