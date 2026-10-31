package compositionroot

// GOAL
//
// Move every "build what I need" into one place: main.
//
//   - handler, service and store in after/ already take what they need as an
//     argument: handler.New(svc), service.New(st), store.New(path). But each
//     one ignores it and builds its own, the way before/ does. Make each one
//     use what it is given.
//   - None of the three may call another package's New, or read the
//     environment.
//   - after/cmd/server/main.go has a build function. Make it read
//     ORDERS_FILE through getenv and call store.New, service.New and
//     handler.New, in that order. That function is the composition root.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/b5-composition-root/...
//
// Check the rule with grep. It should print nothing:
//
//	grep -rnE '(store|service|handler)\.New|os\.Getenv' patterns/b5-composition-root/after --include=*.go | grep -v /cmd/
//
// Run the same grep on before/ to see the wiring spread out.
