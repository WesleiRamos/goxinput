package goxinput

import "fmt"
import "syscall"

type xboxAPI struct {
	dll *syscall.LazyDLL

	isVBusExists       *syscall.LazyProc
	isControllerExists *syscall.LazyProc
	isControllerOwned  *syscall.LazyProc
	plugin             *syscall.LazyProc
	unplug             *syscall.LazyProc
	unplugForce        *syscall.LazyProc
	setDpad            *syscall.LazyProc
	buttons            map[byte]*syscall.LazyProc
	triggers           map[byte]*syscall.LazyProc
	axis               map[byte]*syscall.LazyProc
}

func (self *xboxAPI) LoadDLL() {
	self.dll = syscall.NewLazyDLL(fmt.Sprintf("dll/vXboxInterface-%s/vXboxInterface.dll", systemArch()))
	self.procs()
}

func (self *xboxAPI) procs() {
	// Status
	self.isVBusExists = self.dll.NewProc("isVBusExists")
	self.isControllerExists = self.dll.NewProc("isControllerExists")
	self.isControllerOwned = self.dll.NewProc("isControllerOwned")

	// Virtual device Plug-In/Unplug
	self.plugin = self.dll.NewProc("PlugIn")
	self.unplug = self.dll.NewProc("UnPlug")
	self.unplugForce = self.dll.NewProc("UnPlugForce")

	// Data Transfer (Data to the device)
	self.buttons = map[byte]*syscall.LazyProc{
		BUTTON_A:     self.dll.NewProc("SetBtnA"),
		BUTTON_B:     self.dll.NewProc("SetBtnB"),
		BUTTON_X:     self.dll.NewProc("SetBtnX"),
		BUTTON_Y:     self.dll.NewProc("SetBtnY"),
		BUTTON_START: self.dll.NewProc("SetBtnStart"),
		BUTTON_BACK:  self.dll.NewProc("SetBtnBack"),
		BUTTON_LS:    self.dll.NewProc("SetBtnThumbL"),
		BUTTON_RS:    self.dll.NewProc("SetBtnThumbR"),
		BUTTON_LB:    self.dll.NewProc("SetBtnShoulderL"),
		BUTTON_RB:    self.dll.NewProc("SetBtnShoulderR"),
	}

	self.triggers = map[byte]*syscall.LazyProc{
		BUTTON_LT: self.dll.NewProc("SetTriggerL"),
		BUTTON_RT: self.dll.NewProc("SetTriggerR"),
	}

	self.axis = map[byte]*syscall.LazyProc{
		AXIS_LX: self.dll.NewProc("SetAxisLx"),
		AXIS_LY: self.dll.NewProc("SetAxisLy"),
		AXIS_RX: self.dll.NewProc("SetAxisRx"),
		AXIS_RY: self.dll.NewProc("SetAxisRy"),
	}

	self.setDpad = self.dll.NewProc("SetDpad")
}

func systemArch() string {
	intSize := 32 << (^uint(0) >> 63)
	if intSize == 32 {
		return "x32"
	}
	return "x64"
}
