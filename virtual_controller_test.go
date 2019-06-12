package goxinput

import "fmt"
import "time"
import "testing"

var controller = NewController()

var buttonMap map[byte]string = map[byte]string{
	BUTTON_A:     "A",
	BUTTON_B:     "B",
	BUTTON_X:     "X",
	BUTTON_Y:     "Y",
	BUTTON_START: "START",
	BUTTON_BACK:  "BACK",
	BUTTON_LS:    "LS",
	BUTTON_RS:    "RS",
	BUTTON_LB:    "LB",
	BUTTON_RB:    "RB",
	BUTTON_LT:    "LT",
	BUTTON_RT:    "RT",
	AXIS_LX:      "LX",
	AXIS_LY:      "LY",
	AXIS_RX:      "RX",
	AXIS_RY:      "RY",
}

func TestCheckDriver(t *testing.T) {
	if !controller.IsVBusExists() {
		t.Fatalf("Driver não instalado")
	}
}

func TestPluginController(t *testing.T) {
	if error := controller.PlugIn(); error != nil {
		t.Fatalf("Não foi possível plugar o controle")
	}

	time.Sleep(time.Second)
}

func TestDpads(t *testing.T) {
	_dpads := []int{DPAD_UP, DPAD_DOWN, DPAD_LEFT, DPAD_RIGHT, DPAD_OFF}

	for _, dpad := range _dpads {
		if ok := controller.SetDpad(dpad); !ok {
			t.Fatalf("Não foi possível pressionar o dpad %d", dpad)
		}

		time.Sleep(time.Second)
	}
}

func TestPressButtons(t *testing.T) {
	buttons := []byte{
		BUTTON_A, BUTTON_B, BUTTON_X, BUTTON_Y,
		BUTTON_START, BUTTON_BACK,
		BUTTON_LS, BUTTON_RS, BUTTON_LB, BUTTON_RB,
	}

	for _, button := range buttons {

		if ok := controller.SetBtn(button, true); !ok {
			t.Fatalf("Não foi possível pressionar o botão %s", buttonMap[button])
		}

		time.Sleep(time.Second)

		if ok := controller.SetBtn(button, false); !ok {
			t.Fatalf("Não foi possível soltar o botão %s", buttonMap[button])
		}
	}
}

func TestSetTriggers(t *testing.T) {
	triggers := []byte{BUTTON_LT, BUTTON_RT}

	for _, trigger := range triggers {

		btn := buttonMap[trigger]

		fmt.Printf("Aplicando pressão ao botão %s", btn)

		// Simular pressão progressiva
		for i := 0; i <= 10; i++ {

			fmt.Printf(".")

			if ok := controller.SetTrigger(trigger, float32(i)/10); !ok {
				t.Fatalf("Não foi possível adicionar pressão ao trigger %s", btn)
			}

			time.Sleep(time.Millisecond * 100)
		}

		fmt.Printf("OK\nRemovendo pressão..")

		// Remove a pressão do botão
		if ok := controller.SetTrigger(trigger, 0); !ok {
			t.Fatalf("Não foi possível resetar a pressão do trigger %s", btn)
		}

		println("OK")
	}
}

func TestSetAxis(t *testing.T) {
	_axis := []byte{AXIS_LX, AXIS_LY, AXIS_RX, AXIS_RY}

	for _, axis := range _axis {

		ax := buttonMap[axis]

		fmt.Printf("Aplicando movimento ao analogico %s", ax)

		// Simula movimento progressivo
		for i := -10; i <= 10; i++ {

			if i%2 == 0 {
				fmt.Printf(".")
			}

			if ok := controller.SetAxis(axis, float32(i)/10); !ok {
				t.Fatalf("Não foi possível adicionar movimento ao analogico %s", ax)
			}

			time.Sleep(time.Millisecond * 50)
		}

		fmt.Printf("OK\nRemovendo movimento..")

		// Remove o movimento do analogico
		if ok := controller.SetAxis(axis, 0); !ok {
			t.Fatalf("Não foi possível resetar o movimento do analogico %s", ax)
		}

		println("OK")
	}
}

func TestUnplug(t *testing.T) {
	if !controller.Unplug() {
		t.Fatalf("Não foi possível desplugar o controle com id %d", controller.id)
	}
}
