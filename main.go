package main

import (
	"github.com/zserge/webview"
	"flag"
	"io/ioutil"
	"os"
	"sync"
)


type Stui struct {
	View webview.WebView
	*Duplex
	readyWhenClosed chan bool
}

func (s Stui) listenAndApply() {
	subs := make([]SubscribeCommand, 0)
	<- s.readyWhenClosed
	dup := s.Duplex
	for v := range dup.In {
		switch cmd := v.(type) {
		case SubscribeCommand: // odd: the go driver is collecting subs...
			subs = append(subs, cmd)
			cmd.Apply(s.View)
		case PostElemCommand:
			cmd.Apply(s.View)
			for _, sub := range subs {
				sub.Apply(s.View)
			}
		case CloseCommand:
			cmd.Apply(s.View)
			break
		default: // odd: ... but the webview is collecting CSS
			cmd.Apply(s.View)
		}
	}
}

func genStartFile() string {
	f, _ := ioutil.TempFile("", "stui")
	f.WriteString(`
<body>
    <div id="app"></div>
</body>
<script type="text/javascript">
	window.external.invoke('["hi"]');
</script>`)
	f.Close()
	os.Rename(f.Name(), f.Name() + ".html")
	return "file://" + f.Name() + ".html"
}

func NewStui(d *Duplex) Stui {
	readyWhenClosed := make(chan bool, 1)
	cb := func(w webview.WebView, s string) {
		d.Out <- s
		/* I need this callback to signal Stui.
		   I don't want the Conn to know about the Webview readiness,
		   and Stui itself can't exist in time for this function to close over it.
		   So I have to close over some reftype which will be included in Stui.
		   ... and what's the correct reftype for passing messages? A channel!
		   So: here I am using a channel. The first call here sees the channel
		   with no messages, goes to `default:`, and closes the channel. Stui
		   and further calls here see a closed channel and continue with nil
		*/ 
		select {
		case <- readyWhenClosed:
		default:
			close(readyWhenClosed)
		}
	}
	startingFileName := genStartFile()
	view := webview.New(webview.Settings{
		URL: startingFileName,
		Width:     300,
		Height:    400,
		Title:     "Hi Stui",
		Resizable: true,
		ExternalInvokeCallback: cb,
	})

	return Stui{view, d, readyWhenClosed}
}

func main() {
	zmq := flag.String("zmq", "", "Socket name, if using zmq, such as ipc:///tmp/stui. Will use stdio if omitted or blank")
	flag.Parse()
	var conn Conn
	d := Duplex{make(chan Command, 9), make(chan string, 9)}
	if *zmq == "" {
		conn = StdioConn(&d)
	} else {
		conn = NewZMQConn(&d, *zmq)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go conn.Start(wg)
	s := NewStui(&d)
	go s.listenAndApply()
	defer s.View.Exit()

	s.View.Run()
	d.Out <- `["bye"]`
	close(d.Out)
	wg.Wait()
}