GO=go

bin:
	$(GO) build -o build/gots.exe main.go

run :
	$(GO) run main.go

# this works on windows, change the build extension based on your os
all: bin
	./build/main.exe