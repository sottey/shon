# Makefile for building SHON v0.6 tools

TOOLS = shon vsc

all: $(TOOLS)

shon: ./tooling/
	go build -o bin/shon ./tooling/shon/

vsc: ./tooling/vsextension
	cd ~/projects/shon/tooling/vsextension && vsce package --out ../../bin

clean:
	rm -rf bin
	rm -f $(TOOLS)
