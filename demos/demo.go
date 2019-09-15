package main
import (
	"github.com/pebbe/zmq4"
	"os/exec"
)

func main() {
	stui := exec.Command("stui", "-zmq", "ipc:///tmp/stui")
	stui.Start()
	
	sock, _ := zmq4.NewSocket(zmq4.PAIR)
	sock.Connect("ipc:///tmp/stui")
	
	for {
		event, _ := sock.Recv(0)
		println(event)
		if event == `["hi"]` {
			sock.Send(`["PostElem", "#app", -1, ["button", "hi there!"]]`, 0)
			sock.Send(`["Subscribe", "button", "onclick", ["y"]]`, 0)
			sock.Send(`["Subscribe", "button", "onmouseover", ["x"]]`, 0)
		} else if event == `["bye"]` {
			return
		}
	}
}