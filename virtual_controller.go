package goxinput

import "time"
import "errors"
import "syscall"

var api *xboxAPI = &xboxAPI{}
var bool_to_int = map[bool]int{true: 1, false: 0}

const (
	BUTTON_A byte = iota + 1
	BUTTON_B
	BUTTON_X
	BUTTON_Y
	BUTTON_START
	BUTTON_BACK
	BUTTON_LS
	BUTTON_RS
	BUTTON_LB
	BUTTON_RB
	BUTTON_LT
	BUTTON_RT
	AXIS_LX
	AXIS_LY
	AXIS_RX
	AXIS_RY
)

const (
	DPAD_OFF int = 0
	DPAD_UP  int = 1 << (iota - 1)
	DPAD_DOWN
	DPAD_LEFT
	DPAD_RIGHT
)

func init() {
	api.LoadDLL()
}

type VirtualController struct {
	Id uint
}

// Checa se o driver do ScpVbus está instalado ou não
func (self VirtualController) IsVBusExists() bool {
	r, _, _ := api.isVBusExists.Call()
	return r != 0
}

// Pluga um controle, retorna o id do controle e um erro
func (self *VirtualController) PlugIn() error {
	// Checa se há algum id disponivel
	if availableId := self.availableId(); availableId > 0 {
		api.plugin.Call(uintptr(availableId))

		for !self.IsIdAvailable(availableId) {
			time.Sleep(time.Second)
		}

		self.Id = availableId
		return nil
	}

	return errors.New("Max inputs reached")
}

// Despluga o controle, recebe um valor booleano que pode ser true para caso seja
// necessário forçar o "desplugue" e falso para o processo normal, por padrão é
// assumido que seja valor falso
func (self VirtualController) Unplug(force ...bool) bool {
	var function *syscall.LazyProc = api.unplug

	if len(force) > 0 {
		if force[0] {
			function = api.unplugForce
		}
	}

	r, _, _ := function.Call(uintptr(self.Id))
	return r != 0
}

// Vai verificar se há algum id disponível para ser usado caso não tenha
// nenhum retorna 0
func (self VirtualController) availableId() uint {
	for i := uint(1); i < 5; i++ {
		if !self.IsIdAvailable(i) {
			return i
		}
	}

	return 0
}

// Checa se o controle esta disponivel
func (self VirtualController) IsIdAvailable(id uint) bool {
	r, _, _ := api.isControllerExists.Call(uintptr(id))
	return r != 0
}

// Retorna quanto é a x% de um numero
func (self VirtualController) getPercentValue(percent float32, val int16) int16 {
	return int16(percent * float32(val))
}

// Define o estado de um botão, pode ser pressionado (true) e solto (false)
func (self VirtualController) SetBtn(button byte, pressed bool) bool {
	if _button, ok := api.buttons[button]; ok {
		r, _, _ := _button.Call(uintptr(self.Id), uintptr(bool_to_int[pressed]))
		return r != 0
	}

	return false
}

// Define o valor de um gatilho, recebe um valor de ponto flutuante entre 0 e 1
func (self VirtualController) SetTrigger(trigger byte, value float32) bool {
	if _trigger, ok := api.triggers[trigger]; ok {
		r, _, _ := _trigger.Call(uintptr(self.Id), uintptr(self.getPercentValue(value, 0xff)&0xff))
		return r != 0
	}

	return false
}

// Define o valor de um dos analogicos, recebe um valor de ponto flutuante entre -1.0 e 1.0
func (self VirtualController) SetAxis(axis byte, value float32) bool {
	if _axis, ok := api.axis[axis]; ok {
		r, _, _ := _axis.Call(uintptr(self.Id), uintptr(self.getPercentValue(value, 0x7fff)))
		return r != 0
	}

	return false
}

// Aperta uma das "setinhas"
func (self VirtualController) SetDpad(dpad int) bool {
	r, _, _ := api.setDpad.Call(uintptr(self.Id), uintptr(dpad))
	return r != 0
}

// Cria, inicializa e retorna um novo controle
func NewController() *VirtualController {
	return &VirtualController{}
}
