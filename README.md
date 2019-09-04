
HtmlExpr = selector attrs children

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
            standardized test app+procedure
                drag and drop
        who should choose the ports?
            app should be able to
        think about vdom
            the driver should not enforce TEA
    reliability/usability/clarity
        implement PutHTML
        debug flags
        fix Close
        consistent casing on output
        case-insensitivity on input
            generous message format?
        implement ok/err

    performance
        test performance
            it's fine on the zmq side
            seems 5% cpu/mem for an idle+empty app. high, but acceptable for now
            test scaling with large doms
            button-click spamming that adds elements and css falls behind
                actually, no. only if there are many css-transitions
        css
            accumulating rules? how to delete/overwrite?
            should maintain a copy in the driver
        
    rename gluier? guiso?
    explain message format in a sane way

