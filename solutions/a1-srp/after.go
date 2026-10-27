//go:build ignore

package srp

import (
	"fmt"
	"os"
	"strings"
)

// HOW IT WORKS
//
// Process is cut into three pieces. Total only does arithmetic, Render only
// builds text, and Save only touches the disk. Render takes the total as an
// argument instead of calculating it, so the text never depends on the tax rule.
//
// WHY IT HELPS
//
//   - A tax change edits Total and nothing else.
//   - The sum is tested with no folder, no cleanup and no error to handle.
//   - If Save fails, you still have the correct total.
//
// WHEN NOT TO USE
//
// If the three parts always change together and are never used apart, one
// function is easier to read. Split when you can name two different reasons
// to change it.

// Total returns subtotal plus tax, in cents. Pure: no files, no printing.
func (inv Invoice) Total() int64 {
	var subtotal int64
	for _, it := range inv.Items {
		subtotal += it.UnitCents * it.Qty
	}
	return subtotal + subtotal*inv.TaxPercent/100
}

// Render builds the receipt text for inv with the given total. Pure.
func Render(inv Invoice, total int64) string {
	var b strings.Builder
	for _, it := range inv.Items {
		fmt.Fprintf(&b, "%s x%d %s\n", it.Name, it.Qty, money(it.UnitCents*it.Qty))
	}
	fmt.Fprintf(&b, "TOTAL %s\n", money(total))
	return b.String()
}

// Save writes receipt to path. The only piece allowed to touch the disk.
func Save(path, receipt string) error {
	if err := os.WriteFile(path, []byte(receipt), 0o644); err != nil {
		return fmt.Errorf("save receipt: %w", err)
	}
	return nil
}
