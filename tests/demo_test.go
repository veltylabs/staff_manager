//go:build wasm
package tests
import (
    "testing"
    "github.com/veltylabs/staff_manager/seed"
)
func TestDemo_Valid(t *testing.T) {
    demo := seed.Demo()
    if demo.Label != "Staff" { t.Errorf("expected Staff") }
}
