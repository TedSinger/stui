package main
import (
	"github.com/pebbe/zmq4"
	"os/exec"
)

func main() {
	guise := exec.Command("guise")
	guise.Start()
	
	eventSocket, _ := zmq4.NewSocket(zmq4.PULL)
	eventSocket.Connect("ipc:///tmp/guiseEvents")
	commandSocket, _ := zmq4.NewSocket(zmq4.PUSH)
	commandSocket.Bind("ipc:///tmp/guiseCommands")

	for {
		event, _ := eventSocket.Recv(0)
		if event == `["hi"]` {
			commandSocket.Send(`["SetHtml", "#app", "<button>hi there!</button>"]`, 0)
			commandSocket.Send(`["Sub", "button", "onclick", ["y"]]`, 0)
			commandSocket.Send(`["Sub", "button", "onmouseover", ["x"]]`, 0)
		} else if event == `["bye"]` {
			break
		}
		println(event)
	}
}