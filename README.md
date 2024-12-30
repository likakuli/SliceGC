# SliceGC

An example program to demonstrate the memory leak issue and its fix in the reflector.

You can demonstrate the leak using the following command:
```bash
go run -ldflags="-X 'main.flag=leak'" main.go
```

To demonstrate the issue being fixed, use the following command:
```bash
go run main.go
```