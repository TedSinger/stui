
TODO:
    app architecture
        message log testing demo
        write more demos
            bash demo should be a multi-file selector
            drag and drop
            standardized test app+procedure
            expect?
        subscriptions
            elements can potentially cease matching a sub selector
    reliability/usability/clarity
        errors
            report js errors
            improve Err message
        message format
            case-insensitivity on input? consistency on output?
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
