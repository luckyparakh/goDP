package srp

import (
	"fmt"
	"os"
	"strings"
)

// GOAL
//
// Split Process into three pieces, one per reason to change. Keep the tax
// rule and the receipt format exactly as in before.go.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/a1-srp/

// Total returns subtotal plus tax, in cents. Pure: no files, no printing.
func (inv Invoice) Total() int64 {
	var total int64
	for _, in := range inv.Items {
		total = total + in.Qty*in.UnitCents
	}
	total = total + total*inv.TaxPercent/100
	return total
}

// Render builds the receipt text for inv with the given total. Pure.
func Render(inv Invoice, total int64) string {
	var sb strings.Builder
	for _, in := range inv.Items {
		fmt.Fprintf(&sb, "%s x%d %s\n", in.Name, in.Qty, money(in.UnitCents*in.Qty))
	}
	fmt.Fprintf(&sb, "TOTAL %s\n", money(total))

	return sb.String()
}

// Save writes receipt to path. The only piece allowed to touch the disk.
func Save(path, receipt string) error {
	writeErr := os.WriteFile(path, []byte(receipt), 0o644)
	if writeErr != nil {
		return fmt.Errorf("receipt save failed: %w", writeErr)
	}
	return writeErr
}
