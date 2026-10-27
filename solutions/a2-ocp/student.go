//go:build ignore

package ocp

// The extension: a new tier lives in a new file. No existing file changes.
func init() {
	Register("student", func(cents int64) int64 { return cents * 75 / 100 })
}
