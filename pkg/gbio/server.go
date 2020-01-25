package gbio

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kil0meters/goboy/pkg/gb"
	"github.com/kil0meters/goboy/pkg/gbio/iopixel"
)

// ServerIOBinding A struct used for containing things - wow
type ServerIOBinding struct {
	mon        *iopixel.PixelsIOBinding
	wsSessions []*websocket.Conn
	upgrader   websocket.Upgrader
}

// SendToSockets sends a message to all sockets in an array
func (io *ServerIOBinding) SendToSockets(data []byte) {
	for i := 0; i < len(io.wsSessions); i++ {
		io.wsSessions[i].WriteMessage(1, data)
	}
}

// ShortPressButton trigger a button for a short period of time
func (io *ServerIOBinding) ShortPressButton(button gb.Button) {
	key := ""

	switch button {
	case gb.ButtonUp:
		key = "up-button"
	case gb.ButtonDown:
		key = "down-button"
	case gb.ButtonLeft:
		key = "left-button"
	case gb.ButtonRight:
		key = "right-button"
	case gb.ButtonA:
		key = "a-button"
	case gb.ButtonB:
		key = "b-button"
	case gb.ButtonStart:
		key = "start-button"
	case gb.ButtonSelect:
		key = "select-button"
	default:
		key = ""
	}

	io.SendToSockets([]byte(fmt.Sprintf("{\"pressed_key\": \"%s\"}", key)))
	io.mon.Gameboy.PressButton(button)
	time.Sleep(100 * time.Millisecond)
	io.mon.Gameboy.ReleaseButton(button)
	io.SendToSockets([]byte(fmt.Sprintf("{\"released_key\": \"%s\"}", key)))
}

// PressUp On up press
func (io *ServerIOBinding) PressUp(w http.ResponseWriter, r *http.Request) {
	println("Pressed UP")
	io.ShortPressButton(gb.ButtonUp)
}

// PressDown On down press
func (io *ServerIOBinding) PressDown(w http.ResponseWriter, r *http.Request) {
	println("Pressed DOWN")
	io.ShortPressButton(gb.ButtonDown)
}

// PressLeft On left press
func (io *ServerIOBinding) PressLeft(w http.ResponseWriter, r *http.Request) {
	println("Pressed LEFT")
	io.ShortPressButton(gb.ButtonLeft)
}

// PressRight On right press
func (io *ServerIOBinding) PressRight(w http.ResponseWriter, r *http.Request) {
	println("Pressed RIGHT")
	io.ShortPressButton(gb.ButtonRight)
}

// PressA On a press
func (io *ServerIOBinding) PressA(w http.ResponseWriter, r *http.Request) {
	println("Pressed A")
	io.ShortPressButton(gb.ButtonA)
}

// PressB On b press
func (io *ServerIOBinding) PressB(w http.ResponseWriter, r *http.Request) {
	println("Pressed B")
	io.ShortPressButton(gb.ButtonB)
}

// PressStart On start press
func (io *ServerIOBinding) PressStart(w http.ResponseWriter, r *http.Request) {
	println("Pressed START")

	io.ShortPressButton(gb.ButtonStart)
}

// PressSelect On select press
func (io *ServerIOBinding) PressSelect(w http.ResponseWriter, r *http.Request) {
	println("Pressed SELECT")
	io.ShortPressButton(gb.ButtonSelect)
}

func remove(slice []*websocket.Conn, ws *websocket.Conn) []*websocket.Conn {
	s := 0
	for s = 0; s < len(slice); s++ {
		if slice[s] == ws {
			break
		}
	}

	return append(slice[:s], slice[s+1:]...)
}

func internalError(ws *websocket.Conn, msg string, err error) {
	println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

// ServeWS Serves websocket
func (io *ServerIOBinding) ServeWS(w http.ResponseWriter, r *http.Request) {
	ws, err := io.upgrader.Upgrade(w, r, nil)
	if err != nil {
		println(err.Error())
		return
	}

	index := len(io.wsSessions)
	io.wsSessions = append(io.wsSessions, ws)

	println("Opening WS session", index)
	defer ws.Close()
	for {
		// io.wsSessions;
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			println("read:", err)
			break
		}
		println("Received: ", msg)
		err = ws.WriteMessage(msgType, msg)
		if err != nil {
			println("write:", err)
			break
		}
	}

	io.wsSessions = remove(io.wsSessions, ws)
	println("Closing WS session", index)
}

// InitializeInputServer Initializes button input web server
func InitializeInputServer(mon *iopixel.PixelsIOBinding) {
	println("Starting server on localhost:8080")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	io := &ServerIOBinding{mon: mon, wsSessions: make([]*websocket.Conn, 0), upgrader: upgrader}

	http.HandleFunc("/press-up", io.PressUp)
	http.HandleFunc("/press-down", io.PressDown)
	http.HandleFunc("/press-left", io.PressLeft)
	http.HandleFunc("/press-right", io.PressRight)
	http.HandleFunc("/press-a", io.PressA)
	http.HandleFunc("/press-b", io.PressB)
	http.HandleFunc("/press-start", io.PressStart)
	http.HandleFunc("/press-select", io.PressSelect)

	http.HandleFunc("/ws", io.ServeWS)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
