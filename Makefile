build:
	go build -o tui-hub
cp:
	cp tui-hub ~/.local/bin/

install: build cp