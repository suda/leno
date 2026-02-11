.PHONY: build dev clean

build:
	npm run build
	go build -o leno .

dev:
	npm run dev

clean:
	rm -rf public/build
	rm -f leno
