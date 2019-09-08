# Guise

## A language-agnostic GUI driver using messages passed through stdio or a 0mq socket

```bash
echo '["PostElem", "#app", -1, ["button", "stdio is all you need!"]]["Subscribe", "button", "onmousemove", ["x", "y"]]' | guise
# ["hi"]
# ["ui","button","onmousemove",{"x":78,"y":23}]
# ["ui","button","onmousemove",{"x":77,"y":23}]
# ["ui","button","onmousemove",{"x":76,"y":24}]
# ["bye"]
```

### Why?

Typical GUI toolkits require language-native bindings. This limits the options available to smaller languages, and forces every application to tightly couple itself to both the toolkit and the particular binding used.

_Guise_ is primarily a message format. Any language that can form and parse JSON can build a GUI. This message format does *not* specify what element tags or attributes to use. The implementation here, using the wonderful [WebView](github.com/zserge/webview), requires HTML tags and CSS selectors, but if another implementation wants to use `GtkHBox` or `QBoxLayout` instead of `div`s, that's perfectly fine.

Another implementation could be headless. To test your model, just write an `expect` script with pre-fabricated events. To test your view, use pre-fabricated commands.

### Performance?
Eh. It's fine. The bottleneck is the GTK webkit. IPC message passing is really light.

### Production-ready?
Not quite yet. See the issues tagged _Blocker_