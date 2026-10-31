//go:build ignore

package consumeriface

// HOW IT WORKS
//
// billing now declares the Invoice type and a one-method Invoices interface,
// in its own package. store imports billing and returns billing.Invoice, so
// *store.Mem satisfies billing.Invoices without either package naming that
// fact. go list shows the arrow the other way round: store depends on billing,
// and billing depends on nothing of store's.
//
// In A5 the provider returned a plain string, so the arrow simply vanished.
// Here the provider hands back the consumer's own type, so it has to import
// the consumer. When that happens the interface must live with the consumer,
// or Go reports an import cycle.
//
// The solution files are after/billing/billing.go and after/store/store.go.
//
// WHY IT HELPS
//
// billing owns its own record and its own question. A second store, or a test
// fake, imports billing and nothing else. store can use any billing type it
// needs, because billing will never import it back.
//
// WHEN NOT TO USE
//
// If the interface is a shared vocabulary that many unrelated packages use,
// like io.Reader, it belongs in a small package of its own, not in any one
// consumer. And with one consumer, one implementation and no fake, you may not
// need the package split or the interface at all.
