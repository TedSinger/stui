package main

import (
	"github.com/zserge/webview"
)


type Guise struct {
	View webview.WebView
	Conn Conn
}

func (g Guise) listenAndApply() {
	subs := make([]SubCommand, 0)
	for {
		v := g.Conn.Recv()
		switch cmd := v.(type) {
		case SubCommand: // odd: the go driver is collecting subs...
			subs = append(subs, cmd)
			cmd.Apply(g.View)
		case PostHTMLCommand:
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
	addr := "ipc:///tmp/guise"
	c := NewZMQConn(addr)
	g := NewGuise(c)
	go g.listenAndApply()
	defer g.View.Exit()
	defer g.Conn.Send(`["bye"]`)
	defer g.Conn.Stop()
	
	g.View.Run()
}