all: main.go
	go build
.PHONY: all

install: bluetooth-buttons bluetooth-buttons.service
	install -m755 -groot -oroot bluetooth-buttons /usr/local/bin/bluetooth-buttons
	install -m644 -groot -oroot systemd.service /usr/lib/systemd/system/bluetooth-buttons.service
	install -m644 -groot -oroot udev.rules /etc/udev/rules.d/90-bluetooth-buttons.rules
	systemctl daemon-reload
.PHONY: install

uninstall:
	rm /etc/udev/rules.d/90-bluetooth-buttons.rules
	systemctl stop bluetooth-buttons.service
	rm /usr/lib/systemd/system/bluetooth-buttons.service
	rm /usr/local/bin/bluetooth-buttons
	systemctl daemon-reload
