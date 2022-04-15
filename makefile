.PHONY: checkenv distribute

checkenv:
ifndef version
	$(error version is undefined)
endif

clean:
	rm -f csv-chef
	rm -f csv-chef-*.zip

build:
	go build

archives: checkenv csv-chef-mac-intel csv-chef-linux-amd64 csv-chef-linux-arm csv-chef-mac-m1
	mv csv-chef-darwin csv-chef
	zip csv-chef-darwin-amd64-$(version).zip csv-chef README.md LICENSE
	rm csv-chef
	mv csv-chef-arm csv-chef
	zip csv-chef-linux-arm-$(version).zip csv-chef README.md LICENSE
	rm csv-chef
	mv csv-chef-amd64 csv-chef
	zip csv-chef-linux-amd64-$(version).zip csv-chef README.md LICENSE
	rm csv-chef
	mv csv-chef-darwin-m1 csv-chef
	zip csv-chef-darwin-arm64-$(version).zip csv-chef README.md LICENSE
	rm csv-chef

csv-chef-mac-intel:
	GOOS=darwin GOARCH=amd64 go build -o csv-chef-darwin

csv-chef-mac-m1:
	GOOS=darwin GOARCH=arm64 go build -o csv-chef-darwin-m1

csv-chef-linux-arm:
	GOOS=linux GOARCH=arm go build -o csv-chef-arm

csv-chef-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o csv-chef-amd64