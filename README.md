
HtmlExpr = [selector, Dict key value, List HtmlExpr] | HtmlString

Command = Sub selector onWhat values
    | PostHtml selector index HtmlExpr
    | PutHtml selector HtmlExpr
    | DeleteHtml selector
    | PatchAttrs selector attrs
    | PatchCss selector styles
    | Close

Event = Hi | Bye
    | UI selector onWhat values
    | Ok hash ??
    | Err errorMessage

TODO:
    app architecture
        write more demos
            drag and drop
        who should choose the ports?
            app should be able to
        think about vdom
            the driver should not enforce TEA
    reliability/usability/clarity
        debug flags
        standardized test app+procedure
        test Close
        consistent casing on output
        case-insensitivity on input
            generous message format?
        implement ok/err

    performance
        test performance
            it's fine on the zmq side...
        css
            accumulating rules? how to delete/overwrite?
            should maintain a copy in the driver
        avoid rebuilding SetHtml string messages
    
    rename gluier? guiso?

