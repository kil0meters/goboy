package gbio

import (
	"net/http"
	"time"

	"github.com/kil0meters/goboy/pkg/gb"
	"github.com/kil0meters/goboy/pkg/gbio/iopixel"
)

type ServerIOBinding struct {
	mon *iopixel.PixelsIOBinding
}

func (io *ServerIOBinding) ShortPressButton(button gb.Button) {
	io.mon.Gameboy.PressButton(button)
	time.Sleep(100 * time.Millisecond)
	io.mon.Gameboy.ReleaseButton(button)
}

func (io *ServerIOBinding) PressUp(w http.ResponseWriter, r *http.Request) {
	println("Pressed UP")
	io.ShortPressButton(gb.ButtonUp)
}

func (io *ServerIOBinding) PressDown(w http.ResponseWriter, r *http.Request) {
	println("Pressed DOWN")
	io.ShortPressButton(gb.ButtonDown)
}

func (io *ServerIOBinding) PressLeft(w http.ResponseWriter, r *http.Request) {
	println("Pressed LEFT")
	io.ShortPressButton(gb.ButtonLeft)
}

func (io *ServerIOBinding) PressRight(w http.ResponseWriter, r *http.Request) {
	println("Pressed RIGHT")
	io.ShortPressButton(gb.ButtonRight)
}

func (io *ServerIOBinding) PressA(w http.ResponseWriter, r *http.Request) {
	println("Pressed A")
	io.ShortPressButton(gb.ButtonA)
}

func (io *ServerIOBinding) PressB(w http.ResponseWriter, r *http.Request) {
	println("Pressed B")
	io.ShortPressButton(gb.ButtonB)
}

func (io *ServerIOBinding) PressStart(w http.ResponseWriter, r *http.Request) {
	println("Pressed START")
	io.ShortPressButton(gb.ButtonStart)
}

func (io *ServerIOBinding) PressSelect(w http.ResponseWriter, r *http.Request) {
	println("Pressed SELECT")
	io.ShortPressButton(gb.ButtonSelect)
}

func InitializeInputServer(mon *iopixel.PixelsIOBinding) {
	println("Starting server")
	io := &ServerIOBinding{mon: mon}

	http.HandleFunc("/press-up", io.PressUp)
	http.HandleFunc("/press-down", io.PressDown)
	http.HandleFunc("/press-left", io.PressLeft)
	http.HandleFunc("/press-right", io.PressRight)
	http.HandleFunc("/press-a", io.PressA)
	http.HandleFunc("/press-b", io.PressB)
	http.HandleFunc("/press-start", io.PressStart)
	http.HandleFunc("/press-select", io.PressSelect)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
