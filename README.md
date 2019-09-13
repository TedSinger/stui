# Stui

## A language-agnostic GUI driver

Pass messages through stdio or 0mq sockets

```bash
echo '["PostElem", "#app", -1, ["button", "stdio is all you need!"]]
      ["Subscribe", "button", "onmousemove", ["x", "y"]]' | stui
# ["hi"]
# ["ui","button","onmousemove",{"x":78,"y":23}]
# ["ui","button","onmousemove",{"x":77,"y":23}]
# ["ui","button","onmousemove",{"x":76,"y":24}]
# ["bye"]
```
![A simple graphical UI demo. The mouse cursor moves over the button and "onmousemove" events with coordinates are printed in the terminal](/demos/tiny-bash-demo.gif)

### Why?

Typical GUI toolkits require language-native bindings. This limits the options available to smaller languages, and forces every application to tightly couple itself to both the toolkit and the particular binding used.

With `stui`, any language that can form and parse JSON can build a GUI.

`stui` is primarily a message scheme. This message scheme does *not* specify what element tags or attributes to use. This implementation, using the wonderful [WebView](https://github.com/zserge/webview), requires HTML tags and CSS selectors. Another driver could use a `GtkHBox` or `QBoxLayout` or `Tk.Frame` or `wx.BoxSizer` instead of a `div`, but keep the same message structure and semantics.

You can write your own - _trivial_ - headless driver for testing your data model. Just feed your program a simple script of pre-fabricated events.

### Performance?
Depends entirely on the widget toolkit. `gtkwebview` is much smaller than a Chromium Embedded Framework app, but I've noticed it struggles with dozens of simultaneous CSS transitions.

### Stable?
Not quite yet. See the issues tagged _Blocker_. Mostly I want a second pair of eyes on the message scheme.