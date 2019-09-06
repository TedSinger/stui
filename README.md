
HtmlExpr = {selector: string, attrs: Dict String String, children: List HtmlExpr}

Command = Subscribe {selector: string, onWhat: string, values: List String}
    | PostElem {selector: string, index: int, html: HtmlExpr}
    | PutElem {selector: string, html: HtmlExpr}
    | DeleteElem {selector: string}
    | PatchAttrs {selector: string, attrs: Dict String String}
    | PatchStyles {selector: string, styles: Dict String String}
    | Close

Event = Hi | Bye
    | UI {selector: string, onWhat: string, values: List String}
    | Err {errorMessage: string}

TODO:
    app architecture
        write more demos
            bash demo should be a multi-file selector
            drag and drop
            standardized test app+procedure
        subscriptions
            elements can potentially cease matching a sub selector
    reliability/usability/clarity
        implement PutElem
        errors
            report js errors
            improve Err message
        message format
            sane explanation of grammar
            case-insensitivity on input? consistency on output?
            abbreviations:
                omitted attrs or children in HtmlExpr assumed empty
                ["label", "some text"] -> ["label", {"textContent": "some text"}, []]
        caller should choose zmq socket name

    performance
        test performance
            it's fine on the zmq side
            seems 5% cpu/mem for an idle+empty app. high, but acceptable for now
            test scaling with large doms
            button-click spamming that adds elements and css falls behind
                actually, no. only if there are many css-transitions
        css
            accumulating rules? how to delete/overwrite?
            maintain a copy in the driver?
