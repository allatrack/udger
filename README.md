# Udger Go client (data v3)
## Motivation
There is an official golang udger client: [udger](https://github.com/udger/udger).    
But the problem is that it does not provide all the data user may need.      
We looked at the [PHP udger client](https://github.com/udger/udger-php) which provides much more information about
the user agent and crawler. We took the data structure from the PHP client, approaches from golang client
and put it together.   
This client works in two modes:    
 - Fast load - not so fast execution. With this approach we just make DB queries when needed.   
 - Slow first load - faster execution. With this approach we load needed data in memory and then work with it(make as less DB   queries as possible).

## Install
```
go get github.com/allatrack/udger
```
## Usage
If you want to use first approach (fast load - not so fast parsing):
```
import (
     udger "github.com/allatrack/udger/parser"
)

udgerFS, err := udger.NewFS("path to db")
if err != nil {
    log.Fatal(err)
}
userAgent, err = udgerFS.ParseUa("user agent")
if err != nil {
    log.Fatal(err)
}
ipAddress = udgerFS.ParseIp(`101.0.64.0`)
```
If you want to use second approach (slow load - faster parsing):
```
udgerSF, err := udger.NewSF("path to db")
if err != nil {
    log.Fatal(err)
}
userAgent, err := udgerSF.ParseUa("user agent")
if err != nil {
    log.Fatal(err)
}
ipAddress := udgerSF.ParseIp(`101.0.64.0`)
```

## Running tests
```
go test ./...
```
## Documentation
For detailed documentation and basic usage examples, please see the examples folder and tests provided.
