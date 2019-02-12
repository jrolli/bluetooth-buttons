package main

import (
	"log"
	"os"
	"regexp"

	"github.com/gvalkov/golang-evdev"
	"github.com/sashko/go-uinput"
)

func main() {
	log.Print("Scanning input devices...")
	devices, err := evdev.ListInputDevices()
	if err != nil {
		log.Fatal(err)
	}

	target_re := regexp.MustCompile("Dell Active Pen PN579X Keyboard")
	var target *evdev.InputDevice
	found_device := false
	for _, device := range devices {
		log.Print(device)
		if target_re.MatchString(device.Name) {
			log.Print(device)
			target = device
			err = target.Grab()
			log.Printf("Found PN579X at %s", device.Fn)
			if err != nil {
				log.Fatal(err)
			}

			defer target.Release()
			found_device = true
			break
		}
	}
	if !found_device {
		log.Fatal("Failed to find PN579X")
	}

	log.Print("Creating 'uinput' keyboard...")
	kbd, err := uinput.CreateKeyboard()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Entering read loop...")
	for {
		// log.Print("Reading event...")
		ev, err := target.ReadOne()
		if err != nil {
			log.Fatal(err)
		}

		if ev.Type == evdev.EV_KEY && ev.Value == 0 {
			switch ev.Code {
			case evdev.KEY_F20:
				// Single press
				kbd.KeyPress(evdev.KEY_PAGEDOWN)
			case evdev.KEY_F19:
				// Double press
				kbd.KeyPress(evdev.KEY_PAGEUP)
			case evdev.KEY_F18:
				// Long press
				kbd.KeyDown(evdev.KEY_LEFTCTRL)
				kbd.KeyPress(evdev.KEY_RIGHT)
				kbd.KeyUp(evdev.KEY_LEFTCTRL)
			}
		}
	}

	os.Exit(0)
}
