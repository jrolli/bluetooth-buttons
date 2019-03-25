debug: main.go
	go build
.PHONY: debug

release: main.go
	go build -ldflags="-s -w"
.PHONY: release

install: bluetooth-buttons support/systemd.service support/udev.rules
	install -m755 -groot -oroot bluetooth-buttons /usr/local/bin/bluetooth-buttons
	install -m644 -groot -oroot support/systemd.service /usr/lib/systemd/system/bluetooth-buttons.service
	install -m644 -groot -oroot support/udev.rules /etc/udev/rules.d/90-bluetooth-buttons.rules
	systemctl daemon-reload
.PHONY: install

uninstall:
	rm /etc/udev/rules.d/90-bluetooth-buttons.rules
	systemctl stop bluetooth-buttons.service
	rm /usr/lib/systemd/system/bluetooth-buttons.service
	rm /usr/local/bin/bluetooth-buttons
	systemctl daemon-reload
