# Scope of test files

`_test.go` files share the same variables and function scoping. For example, we have a folder structure like this:
```
./project
    main.go
    go.mod
    go.sum
    logic.go
    logic_test.go
    function.go
    function_test.go
```

Variables and functions defined in `logic_test.go` will be accessible from `function_test.go`