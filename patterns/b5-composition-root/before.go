package compositionroot

// PROBLEM
//
// Look in before/. Each package builds the thing it needs:
//
//	handler.New() calls service.New()
//	service.New() calls store.New()
//	store.New()   reads the ORDERS_FILE environment variable
//
// So the wiring of the program is spread over three packages, and nobody can
// see it in one place. To test the handler you must set an environment
// variable and write a real file, because the handler secretly builds the
// store. And nothing can be swapped: there is no way to hand the handler a
// different service, or the service a different store.
//
// TRY IT
//
// Read TestBefore_HandlerNeedsEnvAndFile in root_test.go. A test of an HTTP
// handler has to know about an environment variable and a JSON file. Then run
// the grep from after.go against before/ and see where the wiring lives.
//
// Do not edit before/. Your work goes in after/.
