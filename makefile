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

archives: checkenv csv-chef csv-chef-arm csv-chef-386
	mv csv-chef-darwin csv-chef
	zip csv-chef-darwin-amd64-$(version).zip csv-chef README.md LICENSE
	rm csv-chef
	mv csv-chef-arm csv-chef
	zip csv-chef-linux-arm-$(version).zip csv-chef README.md LICENSE
	rm csv-chef
	mv csv-chef-386 csv-chef
	zip csv-chef-linux-386-$(version).zip csv-chef README.md LICENSE
	rm csv-chef

csv-chef:
	GOOS=darwin GOARCH=amd64 go build -o csv-chef-darwin

csv-chef-arm:
	GOOS=linux GOARCH=arm go build -o csv-chef-arm

csv-chef-386:
	GOOS=linux GOARCH=386 go build -o csv-chef-386

upload: checkenv archives
	curl -s -u dstockton@i3logix.com -X POST https://api.bitbucket.org/2.0/repositories/i3logix/csv-chef/downloads \
	-F files=@csv-chef-darwin-amd64-$(version).zip \
	-F files=@csv-chef-linux-arm-$(version).zip \
	-F files=@csv-chef-linux-386-$(version).zip