package consumeriface

// PROBLEM
//
// Look in before/. The billing package imports the store package, because
// store declares both the Invoice type and the Store interface that billing
// calls. So billing's own record lives in the storage package, and the import
// arrow points from the code that decides (billing) to the code that stores
// (store).
//
// That holds only until store needs something of billing's. The first thing
// it needs is usually the Invoice type itself, on the day someone decides
// Invoice is a billing idea and moves it home. Then billing imports store,
// store imports billing, and Go refuses to build it.
//
// TRY IT
//
// Read TestBefore_BillingImportsStore in consumer_test.go. It runs
// `go list -deps` and finds store in billing's list.
//
// Then move the Invoice type from before/store/store.go into
// before/billing/billing.go, make store use billing.Invoice, and run:
//
//	go build ./patterns/b1-consumer-iface/...
//
// Read the error, then undo:
//
//	git checkout -- patterns/b1-consumer-iface/before
//
// Do not keep that edit. Your work goes in after/.
