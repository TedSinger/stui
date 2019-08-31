package main

import (
	"github.com/zserge/webview"
	"fmt"
	"html/template"
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

type SetHTMLCommand struct {
	Selector string
	Html string
}

func (s SetHTMLCommand) Apply(w webview.WebView) {
	jsFunc := fmt.Sprintf(`document.querySelectorAll("%s").forEach(function (elem) {elem.innerHTML = "%s"})`,
					 s.Selector, template.JSEscapeString(s.Html))
	
	w.Dispatch(func() {w.Eval(jsFunc)})
}
	

func NewSetHTMLCommand(v []interface{}) SetHTMLCommand {
	selector := v[1].(string)
	html := v[2].(string)
	return SetHTMLCommand{selector, html}
}

type SetAttrsCommand struct {
	Selector string
	Attrs map[string]interface{}
}
func (s SetAttrsCommand) Apply(w webview.WebView) {
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
	
func NewSetAttrsCommand(v []interface{}) SetAttrsCommand {
	selector := v[1].(string)
	attrs := v[2].(map[string]interface{})
	return SetAttrsCommand{selector, attrs}
}


type SetCSSCommand struct {
	Selector string
	Styles map[string]string
}

func (s SetCSSCommand) Apply(w webview.WebView) {
	cssText := s.Selector + " {\n"
	for attr, value := range s.Styles {
		cssText += "  " + attr + ": " + value + ";\n"
	}
	cssText += "}"
	w.Dispatch(func() {w.InjectCSS(cssText)})
}

func NewSetCSSCommand(v []interface{}) SetCSSCommand {
	selector := v[1].(string)
	rawStyles := v[2].(map[string]interface{})
	actualStyles := make(map[string]string)
	for attribute, rawValue := range rawStyles {
		actualValue := rawValue.(string)
		actualStyles[attribute] = actualValue
	}
	return SetCSSCommand{selector, actualStyles}
}

type CloseCommand struct {}

func (c CloseCommand) Apply(w webview.WebView) {
	w.Exit()
}
