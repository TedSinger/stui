package main

import (
	"sync"
	"github.com/pebbe/zmq4"
	"encoding/json"
	"os"
	"io"
)

type RawMessage []interface{}

func (r RawMessage) toCommand() Command {
	var m Command
	kind := r[0].(string)

	if kind == "Subscribe" {
		m = NewSubscribeCommand(r)
	} else if kind == "PatchAttrs" {
		m = NewPatchAttrsCommand(r)
	} else if kind == "PostElem" {
		m = NewPostElemCommand(r)
	} else if kind == "PutElem" {
		m = NewPutElemCommand(r)
	} else if kind == "DeleteElem" {
		m = NewDeleteElemCommand(r)
	} else if kind == "PatchStyles" {
		m = NewPatchStylesCommand(r)
	} else if kind == "Close" {
		m = CloseCommand{}
	} else {
		m = NewErrCommand(r)
	}
	return m
}

type Duplex struct {
	In chan Command
	Out chan string
	Done *sync.WaitGroup
}

func NewDuplex() Duplex {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	d := Duplex{make(chan Command, 9), make(chan string, 9), wg}
	return d
}

type Conn interface {
	Start(*Duplex)
}

type ZMQConn struct {
	addr string
	sock *zmq4.Socket
	duplex *Duplex
}

func NewZMQConn(addr string) ZMQConn {
	sock, _ := zmq4.NewSocket(zmq4.PAIR)
	return ZMQConn{addr, sock, nil}
}

func (z ZMQConn) Start(duplex *Duplex) {
	z.duplex = duplex
	z.sock.Bind(z.addr)
	go z.send()
	go z.recv()
}

func (z ZMQConn) send() {
	for msg := range z.duplex.Out {
		z.sock.Send(msg, 0)
	}
	z.sock.Close()
	z.duplex.Done.Done()
}

func (z ZMQConn) recv() {
	for {
		someBytes, err := z.sock.RecvBytes(0)
		if err != nil {
			close(z.duplex.In)
			break
		} else {
			var r RawMessage
			json.Unmarshal(someBytes, &r)
			z.duplex.In <- r.toCommand()
		}
	}
}

type ToBeFile func() (*os.File, error)

type FileConn struct {
	tobein ToBeFile
	tobeout ToBeFile
	duplex *Duplex
	in io.Reader
	decoder *json.Decoder
	out io.Writer
}


func NewFileConn(in ToBeFile, out ToBeFile) FileConn {
		return FileConn{in, out, nil, nil, nil, nil}
}

func (f FileConn) Start(duplex *Duplex) {
	f.duplex = duplex
	go f.send()
	go f.recv()
}

func (f FileConn) send() {
	f.out, _ = f.tobeout()
	for msg := range f.duplex.Out {
		f.out.Write([]byte(msg + "\n"))
	}
	f.duplex.Done.Done()
}

func (f FileConn) recv() {
	var err error
	f.in, err = f.tobein()
	if err != nil {
		panic(err)
	}
	f.decoder = json.NewDecoder(f.in)
	for {
		var r RawMessage
		err := f.decoder.Decode(&r)
		if r != nil && (len(r) != 0){
			f.duplex.In <- r.toCommand()
		} else if err == io.EOF {
			break
		}
	}
}
