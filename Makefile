dependencies:
	go get github.com/logrusorgru/aurora
	go get github.com/goccy/go-yaml

ifndef PREFIX
	$(error PREFIX is not set)
endif

target: dependencies
	go build
	cd nekoctl && go build

push:
	sudo qemu-nbd --connect=/dev/nbd0 sandbox/rootfs.img
	sudo mount /dev/nbd0p1 root/
	sudo make install PREFIX=root
	sudo umount /dev/nbd0p1
	sudo qemu-nbd --disconnect /dev/nbd0

install: target
	mkdir $(PREFIX)/etc/nekoRC -p
	cp nekoRC $(PREFIX)/sbin/init
	cp nekoctl/nekoctl $(PREFIX)/usr/bin/
	cp -r services $(PREFIX)/etc/nekoRC
	cp inittab.neko.yml $(PREFIX)/etc/nekoRC/
	cp config.neko.yml $(PREFIX)/etc/nekoRC/
