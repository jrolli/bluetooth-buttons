package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "os/exec"
    "regexp"

    "github.com/sashko/go-uinput"
)

func main() {
    dev := flag.String("device", "", "(required) mac address of button")
    handle := flag.String("handle", "", "(required) HCI handle for characteristic")
    gatt := flag.String("gatttool", "/usr/bin/gatttool", "path to gatttool")
    blctl := flag.String("bluetoothctl", "/usr/bin/bluetoothctl", "path to bluetoothctl")
    flag.Parse()

    if *dev == "" && *handle == "" {
        fmt.Fprintln(os.Stderr, "Error: 'dev' and 'handle' are required")
        flag.PrintDefaults()
        os.Exit(1)
    }

    devRe := regexp.MustCompile("^([[:xdigit:]]{2}:){5}[[:xdigit:]]{2}$")

    if devRe.FindString(*dev) == "" {
        fmt.Fprintln(os.Stderr, "Error: 'dev' must be a valid mac (with ':')")
        flag.PrintDefaults()
        os.Exit(1)
    }

    err := resetDevice(*blctl, *dev)
    if err != nil {
        log.Fatal(err)
    }

    kbd, err := uinput.CreateKeyboard()
    if err != nil {
        log.Fatal(err)
    }
    defer kbd.Close()



    err = monitorDevice(*gatt, *dev, *handle, kbd)
    defer resetDevice(*blctl, *dev)
    if err != nil {
        log.Fatal(err)
    }

    os.Exit(0)
}

func monitorDevice(gatttool, device, handle string, kbd uinput.Keyboard) error {
    monCmd := exec.Command(gatttool, "--listen", "--device", device, "--handle", handle, "--char-read")

    stdout, err := monCmd.StdoutPipe()
    if err != nil {
        return err
    }

    stderr, err := monCmd.StderrPipe()
    if err != nil {
        return err
    }

    err = monCmd.Start()
    if err != nil {
        return err
    }

    go pipeLogger(stderr)
    go stateMachine(stdout, handle, kbd)

    return monCmd.Wait()
}

func resetDevice(bluetoothctl, device string) error {
    log.Printf("Disconnecting device (%s)", device)
    resetCmd := exec.Command(bluetoothctl, "disconnect", device)
    resetRe := regexp.MustCompile("^Successful disconnected")

    resetRaw, err := resetCmd.StdoutPipe()
    if err != nil {
        return err
    }

    resetScan := bufio.NewScanner(resetRaw)

    err = resetCmd.Start()
    if err != nil {
        return err
    }

    deviceFound := false
    for resetScan.Scan() {
        if resetRe.FindString(resetScan.Text()) != "" {
            deviceFound = true
        }
    }

    if resetScan.Err() != nil {
        return resetScan.Err()
    }
    if !deviceFound {
        return fmt.Errorf("error: device '%s' not found", device)
    }

    err = resetCmd.Wait()
    if err != nil {
        return err
    }

    log.Printf("Device (%s) successfully disconnected", device)
    return nil
}

func pipeLogger(pipe io.Reader) {
    scanner := bufio.NewScanner(pipe)

    for scanner.Scan() {
        log.Print(scanner.Text())
    }

    log.Fatal("error: pipe logger closed unexpectedly")
}

func stateMachine(stdout io.Reader, handle string, kbd uinput.Keyboard) {
    scanner := bufio.NewScanner(stdout)
    capture := regexp.MustCompile("handle = " + handle + " value: 08 00 6(.) 00 00 00 00 00")

    for scanner.Scan() {
        line := capture.FindStringSubmatch(scanner.Text())
        if line != nil {
            switch line[1] {
                case "d": // Long press
                    kbd.KeyPress(uinput.KeyProg3) // BTN_BASE3
                    log.Print("Button 3")
                case "e": // Double press
                    kbd.KeyPress(uinput.KeyPageUp) // BTN_BASE2
                    log.Print("Button 2")
                case "f": // Single press
                    kbd.KeyPress(uinput.KeyPageDown) // BTN_BASE
                    log.Print("Button 1")
            }
        }
    }

    log.Fatal(scanner.Err())
}
