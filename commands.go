package main

import (
	"github.com/zserge/webview"
	"fmt"
	"encoding/json"
)

type Command interface {
	Apply(webview.WebView)
}

type SubscribeCommand struct {
	selector string
	onWhat string
	values []string
}

func (c SubscribeCommand) Apply(w webview.WebView) {
	jsonDict := "{"
	for _, v := range c.values {
		jsonDict += fmt.Sprintf(`"%s":e.%s,`, v, v) + "\n"
	}
	jsonDict += "}"
	jsCallback := fmt.Sprintf(`function(e) {window.external.invoke(JSON.stringify(["ui", "%s", "%s", %s]))}`, c.selector, c.onWhat, jsonDict)
	jsFunc := fmt.Sprintf(`document.querySelectorAll("%s").forEach(function (elem) {elem.%s = %s})`, 
		c.selector, c.onWhat, jsCallback)
	w.Dispatch(func() {w.Eval(jsFunc)})
}

func NewSubscribeCommand(v []interface{}) SubscribeCommand {
	selector := v[1].(string)
	onWhat := v[2].(string)
	rawValues := v[3].([]interface{})
	values := make([]string, len(rawValues))
	for i, r := range rawValues {
		values[i] = r.(string)
	}
	return SubscribeCommand{selector, onWhat, values}
}

type Elem struct {
	tag string
	attrs map[string]interface{}
	children []Elem
}
func NewElemFromInterface(v interface{}) Elem {
	w := v.([]interface{})
	tag := w[0].(string)
	attrs := w[1].(map[string]interface{})
	rawChildren := w[2].([]interface{})
	actualChildren := make([]Elem, len(rawChildren))
	for i, c := range rawChildren {
		actualChildren[i] = NewElemFromInterface(c)
	}
	return Elem{tag, attrs, actualChildren}
}

func (e Elem) createElement(jsName string) string {
	ret := fmt.Sprintf(`%s = document.createElement("%s");`, jsName, e.tag)
	for k, v := range e.attrs {
		valuePart, _ := json.Marshal(v)
		ret += fmt.Sprintf(`%s.%s = %s;`, jsName, k, string(valuePart)) + "\n"
	}
	for i, c := range e.children {
		childName := fmt.Sprintf(`%s_%d`, jsName, i)
		ret += c.createElement(childName)
		ret += fmt.Sprintf(`%s.appendChild(%s);`, jsName, childName)
	}
	return ret
}

type PostElemCommand struct {
	selector string
	index int
	elem Elem
}

func (c PostElemCommand) Apply(w webview.WebView) {
	elem := c.elem.createElement("tmp")
	var forEachFnBody string
	if c.index == -1 {
		forEachFnBody = elem + `elem.appendChild(tmp);`
	} else {
		forEachFnBody = elem + fmt.Sprintf(`elem.insertBefore(tmp, elem.childNodes[%d]);`, c.index)
	}
	jsFunc := fmt.Sprintf(`document.querySelectorAll("%s").forEach(function (elem) {%s});`,
					 c.selector, forEachFnBody)
	// println(jsFunc)
	w.Dispatch(func() {w.Eval(jsFunc)})
}

func NewPostElemCommand(v []interface{}) PostElemCommand {
	selector := v[1].(string)
	index := v[2].(float64)
	elem := NewElemFromInterface(v[3])
	return PostElemCommand{selector, int(index), elem}
}

type DeleteElemCommand struct {
	selector string
}

func NewDeleteElemCommand(v []interface{}) DeleteElemCommand {
	selector := v[1].(string)
	return DeleteElemCommand{selector}
}

func (c DeleteElemCommand) Apply(w webview.WebView) {
	jsFunc := fmt.Sprintf(`document.querySelectorAll("%s").forEach(function (elem) {elem.remove()});`,
					 c.selector)
	w.Dispatch(func() {w.Eval(jsFunc)})
}

type PatchAttrsCommand struct {
	selector string
	attrs map[string]interface{}
}
func (c PatchAttrsCommand) Apply(w webview.WebView) {
	fnBody := "{"
	for k, v := range c.attrs {
		valuePart, _ := json.Marshal(v)
		fnBody += fmt.Sprintf(`elem.%s = %s;`, k, string(valuePart)) + "\n"
	}
	fnBody += "}"
	jsFunc := fmt.Sprintf(`document.querySelectorAll("%s").forEach(function (elem) %s)`,
					 c.selector, fnBody)
	
	w.Dispatch(func() {w.Eval(jsFunc)})
}
	
func NewPatchAttrsCommand(v []interface{}) PatchAttrsCommand {
	selector := v[1].(string)
	attrs := v[2].(map[string]interface{})
	return PatchAttrsCommand{selector, attrs}
}


type PatchStylesCommand struct {
	selector string
	styles map[string]string
}

func (c PatchStylesCommand) Apply(w webview.WebView) {
	cssText := c.selector + " {\n"
	for attr, value := range c.styles {
		cssText += "  " + attr + ": " + value + ";\n"
	}
	cssText += "}"
	w.Dispatch(func() {w.InjectCSS(cssText)})
}

func NewPatchStylesCommand(v []interface{}) PatchStylesCommand {
	selector := v[1].(string)
	rawStyles := v[2].(map[string]interface{})
	actualStyles := make(map[string]string)
	for attribute, rawValue := range rawStyles {
		actualValue := rawValue.(string)
		actualStyles[attribute] = actualValue
	}
	return PatchStylesCommand{selector, actualStyles}
}

type CloseCommand struct {}

func (c CloseCommand) Apply(w webview.WebView) {
	w.Dispatch(func() {w.Terminate()})
}

type ErrCommand struct {
	original []interface{}
}

func NewErrCommand(v []interface{}) ErrCommand {
	return ErrCommand{v}
}

func (c ErrCommand) Apply(w webview.WebView) {
	jsFunc := fmt.Sprintf(`window.external.invoke(JSON.stringify(["err", "%v"]))`, c.original)
	w.Dispatch(func() {w.Eval(jsFunc)})
}