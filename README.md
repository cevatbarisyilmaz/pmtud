# pmtud

Path MTU Discovery Package for Go.

## Download

```
go get github.com/cevatbarisyilmaz/pmtud
```

## Usage

```go
pmtu, err := Pmtud(addr)
```

## Supported Operating Systems

Currently only Windows and Linux. Possibly Darwin in future.

## Supported IP Versions

IPv4. Possibly works for IPv6 too, though it is not properly tested.
