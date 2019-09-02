package main

import (
	"github.com/zserge/webview"
	"fmt"
	// "html/template"
	"encoding/json"
)

type Command interface {
	Apply(webview.WebView)
}

type SubCommand struct {
	Selector string
	OnWhat string
	Values []string
}

func (s SubCommand) Apply(w webview.WebView) {
	jsonDict := "{"
	for _, v := range s.Values {
		jsonDict += fmt.Sprintf(`"%s":e.%s,`, v, v) + "\n"
	}
	jsonDict += "}"
	jsCallback := fmt.Sprintf(`function(e) {window.external.invoke(JSON.stringify(["ui", "%s", "%s", %s]))}`, s.Selector, s.OnWhat, jsonDict)
	jsFunc := fmt.Sprintf(`document.querySelectorAll("%s").forEach(function (elem) {elem.%s = %s})`, 
		s.Selector, s.OnWhat, jsCallback)
	w.Dispatch(func() {w.Eval(jsFunc)})
}

func NewSubCommand(v []interface{}) SubCommand {
	selector := v[1].(string)
	onWhat := v[2].(string)
	rawValues := v[3].([]interface{})
	values := make([]string, len(rawValues))
	for i, r := range rawValues {
		values[i] = r.(string)
	}
	return SubCommand{selector, onWhat, values}
}

type HTML struct {
	Tag string
	Attrs map[string]interface{}
	Children []HTML
}
func NewHTMLFromInterface(v interface{}) HTML {
	w := v.([]interface{})
	tag := w[0].(string)
	attrs := w[1].(map[string]interface{})
	rawChildren := w[2].([]interface{})
	actualChildren := make([]HTML, len(rawChildren))
	for i, c := range rawChildren {
		actualChildren[i] = NewHTMLFromInterface(c)
	}
	return HTML{tag, attrs, actualChildren}
}

func (h HTML) createElement(jsName string) string {
	ret := fmt.Sprintf(`%s = document.createElement("%s");`, jsName, h.Tag)
	for k, v := range h.Attrs {
		valuePart, _ := json.Marshal(v)
		ret += fmt.Sprintf(`%s.%s = %s;`, jsName, k, string(valuePart)) + "\n"
	}
	for i, c := range h.Children {
		childName := fmt.Sprintf(`%s_%d`, jsName, i)
		ret += c.createElement(childName)
		ret += fmt.Sprintf(`%s.appendChild(%s);`, jsName, childName)
	}
	return ret
}

type PostHTMLCommand struct {
	Selector string
	Index int
	Html HTML
}

func (s PostHTMLCommand) Apply(w webview.WebView) {
	elem := s.Html.createElement("tmp")
	var forEachFnBody string
	if s.Index == -1 {
		forEachFnBody = elem + `elem.appendChild(tmp);`
	} else {
		forEachFnBody = elem + fmt.Sprintf(`elem.insertBefore(tmp, elem.childNodes[%d]);`, s.Index)
	}
	jsFunc := fmt.Sprintf(`document.querySelectorAll("%s").forEach(function (elem) {%s});`,
					 s.Selector, forEachFnBody)
	// println(jsFunc)
	w.Dispatch(func() {w.Eval(jsFunc)})
}

func NewPostHTMLCommand(v []interface{}) PostHTMLCommand {
	selector := v[1].(string)
	index := v[2].(float64)
	html := NewHTMLFromInterface(v[3])
	return PostHTMLCommand{selector, int(index), html}
}

type DeleteHTMLCommand struct {
	Selector string
}

func NewDeleteHTMLCommand(v []interface{}) DeleteHTMLCommand {
	selector := v[1].(string)
	return DeleteHTMLCommand{selector}
}

func (s DeleteHTMLCommand) Apply(w webview.WebView) {
	jsFunc := fmt.Sprintf(`document.querySelectorAll("%s").forEach(function (elem) {elem.remove()});`,
					 s.Selector)
	w.Dispatch(func() {w.Eval(jsFunc)})
}

type PatchAttrsCommand struct {
	Selector string
	Attrs map[string]interface{}
}
func (s PatchAttrsCommand) Apply(w webview.WebView) {
	fnBody := "{"
	for k, v := range s.Attrs {
		valuePart, _ := json.Marshal(v)
		fnBody += fmt.Sprintf(`elem.%s = %s;`, k, string(valuePart)) + "\n"
	}
	fnBody += "}"
	jsFunc := fmt.Sprintf(`document.querySelectorAll("%s").forEach(function (elem) %s)`,
					 s.Selector, fnBody)
	
	w.Dispatch(func() {w.Eval(jsFunc)})
}
	
func NewPatchAttrsCommand(v []interface{}) PatchAttrsCommand {
	selector := v[1].(string)
	attrs := v[2].(map[string]interface{})
	return PatchAttrsCommand{selector, attrs}
}


type PatchCSSCommand struct {
	Selector string
	Styles map[string]string
}

func (s PatchCSSCommand) Apply(w webview.WebView) {
	cssText := s.Selector + " {\n"
	for attr, value := range s.Styles {
		cssText += "  " + attr + ": " + value + ";\n"
	}
	cssText += "}"
	w.Dispatch(func() {w.InjectCSS(cssText)})
}

func NewPatchCSSCommand(v []interface{}) PatchCSSCommand {
	selector := v[1].(string)
	rawStyles := v[2].(map[string]interface{})
	actualStyles := make(map[string]string)
	for attribute, rawValue := range rawStyles {
		actualValue := rawValue.(string)
		actualStyles[attribute] = actualValue
	}
	return PatchCSSCommand{selector, actualStyles}
}

type CloseCommand struct {}

func (c CloseCommand) Apply(w webview.WebView) {
	w.Dispatch(func() {w.Exit()})
}
