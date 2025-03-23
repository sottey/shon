# Makefile for building SHON v0.6 tools

TOOLS = csv2shon json2shon shon2json shonfmt

all: $(TOOLS)

csv2shon: csv2shon.go
	go build -o bin/csv2shon csv2shon.go

json2shon: json2shon.go
	go build -o bin/json2shon json2shon.go

shon2json: shon2json.go
	go build -o bin/shon2json shon2json.go

shonfmt: shonfmt.go
	go build -o bin/shonfmt shonfmt.go

clean:
	rm -rf bin
	rm -f $(TOOLS)
