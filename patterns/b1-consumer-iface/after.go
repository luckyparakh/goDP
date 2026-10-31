package consumeriface

// GOAL
//
// after/ starts as a working copy of before/: billing imports store. Turn the
// arrow round.
//
//   - billing declares the Invoice type and the one-method interface it calls.
//     It imports nothing from store.
//   - store imports billing and returns billing.Invoice. *Mem must still
//     satisfy billing's interface, without either package saying so.
//   - Keep store.NewMem, (*Mem).Add, billing.New and (*Service).Due working
//     the same way, so the tests still compile.
//
// Edit two files: after/billing/billing.go and after/store/store.go.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/b1-consumer-iface/...
//
// See the arrow for yourself, before and after:
//
//	go list -deps ./patterns/b1-consumer-iface/after/billing | grep b1-consumer
//	go list -deps ./patterns/b1-consumer-iface/after/store | grep b1-consumer
