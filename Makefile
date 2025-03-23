# Makefile for building SHON v0.6 tools

TOOLS = csv2shon json2shon shon2json shonfmt

all: $(TOOLS)

csv2shon: ./tooling/csv2shon.go
	go build -o bin/csv2shon ./tooling/csv2shon.go

json2shon: ./tooling/json2shon.go
	go build -o bin/json2shon ./tooling/json2shon.go

shon2json: ./tooling/shon2json.go
	go build -o bin/shon2json ./tooling/shon2json.go

shonfmt: ./tooling/shonfmt.go
	go build -o bin/shonfmt ./tooling/shonfmt.go

clean:
	rm -rf bin
	rm -f $(TOOLS)
