# goxinput
[![GoDoc](https://godoc.org/github.com/WesleiRamos/goxinput?status.svg)](https://godoc.org/github.com/WesleiRamos/goxinput)

Go library for simulate a xbox controller based on [PYXInput](https://github.com/bayangan1991/PYXInput).

## Prerequisites
1. `ScpVBus` drivers from  [https://github.com/shauleiz/vXboxInterface/releases](https://github.com/shauleiz/vXboxInterface/releases)
2. `vXboxInterface` dlls from [https://github.com/bayangan1991/vXboxInterface](https://github.com/bayangan1991/vXboxInterface)

## How to use
Install goxinput using:
```
go get -u github.com/WesleiRamos/goxinput
```

Usage example:
```go
package main

import "time"
import "github.com/WesleiRamos/goxinput"

func main() {
  controller := goxinput.NewController()

  // Check if driver is installed
  if !controller.IsVBusExists() {
    panic("VBus driver is not installed")
  }

  // Plugin controller
  if error := controller.PlugIn(); error != nil {
    panic(error)
  }

  // Press the buttons A, B, X, Y
  buttons := []byte{goxinput.BUTTON_A, goxinput.BUTTON_B, goxinput.BUTTON_X, goxinput.BUTTON_Y}
  for _, button := range buttons {
    controller.SetBtn(button, true)
    time.Sleep(time.Second)
    controller.SetBtn(button, false)
  }

  // Movement left axis stick
  for i := -10; i <= 10; i++ {
    // SetAxis: receive a axis stick and a value between -1.0 and 1.0
    controller.SetAxis(goxinput.AXIS_LX, float32(i)/10)
    time.Sleep(time.Millisecond * 50)
  }
  // Reset left stick
  controller.SetAxis(goxinput.AXIS_LX, 0)

  // Press left trigger
  for i := 0; i <= 10; i++ {
    // SetTrigger: trigger a button (LT, RT) and a value between 0.0 and 1.0
    controller.SetTrigger(goxinput.BUTTON_LT, float32(i)/10)
    time.Sleep(time.Millisecond * 100)
  }
  // Reset the trigger
  controller.SetTrigger(goxinput.BUTTON_LT, 0)

  // Press a DPAD button
  controller.SetDpad(goxinput.DPAD_UP)
  time.Sleep(time.Second * 2)
  controller.SetDpad(goxinput.DPAD_OFF)
  time.Sleep(time.Second * 2)

  // OK, now unplug the controller
  controller.Unplug()
}
```