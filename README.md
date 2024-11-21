# Fuzzing in Go using [GFuzz](https://github.com/system-pclub/GFuzz)

Name: Jeremy Schiemann (61325)
Professor: Prof. Martin Sulzmann


## What is GFuzz?

GFuzz is a fuzzing framework specifically built to detect concurrency bugs in Go programs. Unlike traditional fuzzers that target program inputs, GFuzz focuses on non-deterministic execution schedules caused by concurrent goroutines. By systematically manipulating goroutine schedules, GFuzz can uncover edge-case bugs that are otherwise hard to reproduce.


## How To

You can run the different examples by using the [Makefile](/./Makefile) targets.
This will alter the source code, so make sure to commit all changes beforehand and undo the changes after the run.
Gfuzz will run indefinetly and try to find bugs until you abort the execution.

### Output

After running Gfuzz the [out](./out/) directory will contain the result.

example:

```
[oraclert] selefcm timeout: 1500=== RUN   TestChannelBug
[oraclert] started
Should be buggy
[oraclert]: AfterRun
[oraclert]: AfterRunFuzz
[oraclert]: 1 selects
[oraclert]: CheckBugEnd...
End of unit test. Check bugs
-----New Blocking Bug:
---Primitive location:
/fuzz/target/hello_test.go:15
---Primitive pointer:
0xc000064ac0
-----End Bug

-----FOUND BLOCKING


---Stack:
goroutine 6 [running]:
runtime.DumpAllStack()
	/usr/local/go/src/runtime/myoracle_tmp.go:207 +0x85
gfuzz/pkg/oraclert.AfterRunFuzz(0xc00000a1b0)
	/usr/local/go/src/gfuzz/pkg/oraclert/oracle.go:369 +0xfd
gfuzz/pkg/oraclert.AfterRun(0xc00000a1b0)
	/usr/local/go/src/gfuzz/pkg/oraclert/oracle.go:322 +0x73
hello.TestChannelBug(0xc0000d0500)
	/fuzz/target/hello_test.go:58 +0x325
testing.tRunner(0xc0000d0500, 0x58b9e0)
	/usr/local/go/src/testing/testing.go:1193 +0xef
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:1238 +0x2b5

goroutine 1 [chan receive]:
testing.(*T).Run(0xc0000d0500, 0x5824dd, 0xe, 0x58b9e0, 0x48fe46)
	/usr/local/go/src/testing/testing.go:1239 +0x2dc
testing.runTests.func1(0xc0000d0280)
	/usr/local/go/src/testing/testing.go:1511 +0x78
testing.tRunner(0xc0000d0280, 0xc0000c1de0)
	/usr/local/go/src/testing/testing.go:1193 +0xef
testing.runTests(0xc00000a198, 0x673960, 0x2, 0x2, 0xc1c792ebe3fc48c8, 0x6fca82589, 0x67c060, 0x5822e1)
	/usr/local/go/src/testing/testing.go:1509 +0x305
testing.(*M).Run(0xc0000c5790, 0x0)
	/usr/local/go/src/testing/testing.go:1417 +0x1eb
main.main()
	_testmain.go:45 +0x138

goroutine 8 [chan send]:
hello.TestChannelBug.func1(0xc000096580)
	/fuzz/target/hello_test.go:19 +0x9f
created by hello.TestChannelBug
	/fuzz/target/hello_test.go:16 +0xeb

--- PASS: TestChannelBug (0.30s)
PASS
```

At the top in the `-----New Blocking Bug:` section you can see it found a bug, the interesting part is `---Primitive location:`. It shows you where the bug was found in the code. `---Primitive pointer:` can be ignored.

Below is a stacktrace leading to the bug.


### What to do with it?

This is a rather easy example, here is the code:

```go
ch := make(chan int)
go func() {
    ch <- 1
}()  //<-- this func is were GFuzz is pointing to

select {
case <-ch:
    fmt.Println("Normal")
case <-time.After(300 * time.Millisecond):
    fmt.Println("Should be buggy")
}
```

The problem is that channel `ch` is unbuffered. This means that a send (`ch <- 1`) will block until a corresponding receive (`<-ch`) is ready to accept the value.

The blocking can occur if:
- The send (`ch <- 1`) in the goroutine is not scheduled in time.
- The timeout is not triggered yet because the `select` is stuck waiting on the channel.

It is fixed easily by using a buffered channel instead:
```go
ch := make(chan int, 1)
```

## How does it work? (WIP)

This is the original select statementof the `hello-example`:

```go
select {
	case <-ch:
		fmt.Println("Normal")
	case <-time.After(300 * time.Millisecond):
		fmt.Println("Should be buggy")
	}
```


There are 2 possible routes: a receive from the channel and a timeout.
GFuzz changes the source code so that every branch will be on its own switch case.


```go
switch oraclert.GetSelEfcmSwitchCaseIdx("/fuzz/target/hello_test.go", "17", 2) {
	case 0:
		select {
		case <-ch:
			fmt.Println("Normal")
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
			select {
			case <-ch:
				fmt.Println("Normal")
			case <-time.After(300 * time.Millisecond):
				fmt.Println("Should be buggy")
			}
		}
	case 1:
		select {
		case <-time.After(300 * time.Millisecond):
			fmt.Println("Should be buggy")
		case <-oraclert.SelEfcmTimeout():
			oraclert.StoreLastMySwitchChoice(-1)
			select {
			case <-ch:
				fmt.Println("Normal")
			case <-time.After(300 * time.Millisecond):
				fmt.Println("Should be buggy")
			}
		}
	default:
		oraclert.StoreLastMySwitchChoice(-1)
		select {
		case <-ch:
			fmt.Println("Normal")
		case <-time.After(300 * time.Millisecond):
			fmt.Println("Should be buggy")
		}
	}
```

GFuzz will then generate a random index to choose which case and therefore which select branch is used.
The Timout will always contain the original select, so does the default. The default is used when GFuzz has not yet set an order for this operation.



## Conclusion

| Pro | Contra |
|:---:|:---:|
|No additional code required | No easy setup (quirks, no explanation whats wrong)  |
| seems to find more nieche bugs | doesnt revert changes after run |
|  |  |


## Whats next

- update to gfuzz to latest go version
- write "bigger" go app to test fuzzing
- see if result will be helpful
- compare to default fuzzer
  - does gfuzz find more bugs?
  - whats easier to run?