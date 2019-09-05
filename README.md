
HtmlExpr = {selector: string, attrs: Dict String String, children: List HtmlExpr}

Command = Sub {selector: string, onWhat: string, values: List String}
    | PostHtml {selector: string, index: int, html: HtmlExpr}
    | PutHtml {selector: string, html: HtmlExpr}
    | DeleteHtml {selector: string}
    | PatchAttrs {selector: string, attrs: Dict String String}
    | PatchCss {selector: string, styles: Dict String String}
    | Close

Event = Hi | Bye
    | UI {selector: string, onWhat: string, values: List String}
    | Err {errorMessage: string}

TODO:
    app architecture
        write more demos
            bash demo should be a multi-file selector
            standardized test app+procedure
                drag and drop
    reliability/usability/clarity
        caller should choose transport+names/urls
        change messages to refer to Elems, Attrs, and Styles
        implement Err message
        implement PutHTML
        debug flags
        consistent casing on output
        generous message format
            case-insensitivity?
            abbreviations:
                omitted attrs or children in HtmlExpr assumed empty
                ["label", "some text"] -> ["label", {"textContent": "some text"}, []]

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

