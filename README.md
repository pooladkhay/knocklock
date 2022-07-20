# KnockLock
a simple implementation of [Port Knocking](https://en.wikipedia.org/wiki/Port_knocking)

## How to run
Make sure Go compiler is installed on your machine
```bash
git clone https://github.com/pooladkhay/knocklock.git
cd knocklock 
go run .
```
Program wil start listening on four ports which are hard-coded.
Expected knocking sequence is: `2002, 6006, 3003, 1001`

Open a new terminal window and try connecting using ssh:
```bash
ssh root@127.0.0.1 -p 2002
``` 
or netcat:
```bash
nc 127.0.0.1 2002
```

Try different sequences and see the results.

I have also added some comments to explain the code inside `main.go`