GO=go

bin:
	$(GO) build -o build/ main.go

# this works on windows, change the build extension based on your os
all: bin
	./build/main.exe