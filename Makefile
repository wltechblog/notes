.PHONY: build install clean

build:
	go build -o note

install: build
	cp note $(HOME)/.local/bin/
	cd $(HOME)/.local/bin/ && ln -sf note task

clean:
	rm -f note task
