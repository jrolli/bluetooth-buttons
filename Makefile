all: main.go
	go build
.PHONY: all

install: bluetooth-buttons bluetooth-buttons.service udev.rules
	install -m755 -groot -oroot bluetooth-buttons /usr/local/bin/bluetooth-buttons
	install -m644 -groot -oroot bluetooth-buttons.service /usr/lib/systemd/system/bluetooth-buttons.service
	systemctl daemon-reload
.PHONY: install

uninstall:
	systemctl stop bluetooth-buttons.service
	rm /usr/local/bin/bluetooth-buttons
	rm /usr/lib/systemd/system/bluetooth-buttons.service
	systemctl daemon-reload
