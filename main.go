package main

import (
	"github.com/zserge/webview"
	"flag"
)


type Guise struct {
	View webview.WebView
	Conn
}

func (g Guise) listenAndApply() {
	subs := make([]SubscribeCommand, 0)
	for {
		v := g.Conn.Recv()
		switch cmd := v.(type) {
		case SubscribeCommand: // odd: the go driver is collecting subs...
			subs = append(subs, cmd)
			cmd.Apply(g.View)
		case PostElemCommand:
			cmd.Apply(g.View)
			for _, sub := range subs {
				sub.Apply(g.View)
			}
		default: // odd: ... but the webview is collecting CSS
			cmd.Apply(g.View)
		}
	}
}

func NewGuise(c Conn) Guise {
	c.Start()
		
	cb := func(w webview.WebView, s string) {c.Send(s)}

	view := webview.New(webview.Settings{
		Width:     300,
		Height:    400,
		Title:     "Hi Guise",
		Resizable: true,
		ExternalInvokeCallback: cb,
	})
	c.Send(`["hi"]`)
	return Guise{view, c}
}

func main() {
	connectionType := flag.String("conn", "stdio", "stdio or zmq (ipc:///tmp/guise)")
	flag.Parse()
	var conn Conn
	if *connectionType == "stdio" {
		conn = StdioConn()
	} else if *connectionType == "zmq" {
		conn = NewZMQConn("ipc:///tmp/guise")
	} else {
		println("`conn` must be `stdio` or `zmq`")
		return
	}
	g := NewGuise(conn)
	go g.listenAndApply()
	defer g.View.Exit()
	defer g.Send(`["bye"]`)
	defer g.Stop()
	
	g.View.Run()
}