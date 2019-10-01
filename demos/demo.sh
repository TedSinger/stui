#!/usr/bin/env bash
rm /tmp/stuiIn
mkfifo /tmp/stuiIn

choice=false

app () {
    echo "looping..." >&2
    while read msg ; do
        echo "received an event: " $msg >&2
        if [ '["hi"]' = "$msg" ]; then
            echo '["PostElem", "#app", -1, 
            ["label", "I can make a GUI in *bash*?!"]]'
            echo '["PostElem", "#app", -1,  
                ["input", {"type":"checkbox"}]]'
            echo '["PostElem", "#app", -1,
                ["button", "Confirm"]]'
            echo '["Subscribe", "input", "onchange", ["target.checked"]]'
            echo '["Subscribe", "button", "onclick", []]'
            echo '["Subscribe", "body", "onmousemove", ["x","y"]]'
            echo '["PatchStyles", ".path", {"position":"absolute"}]'
        elif [ '["bye"]' = "$msg" ]; then
            # echo "quitting..." >&2
            break
        else
            evType=$(echo $msg | jq .[2])
            
            if [ '"onclick"' = "$evType" ]; then
                echo $choice >&2
                echo '["Close"]'
            elif [ '"onchange"' = "$evType" ]; then
                choice=$(echo $msg | jq .[3].\"target.checked\")
            elif [ '"onmousemove"' = "$evType" ]; then
                x=`echo $msg | jq '.[3].x' `
                y=`echo $msg | jq '.[3].y'`
                echo '["PostElem", "#app", -1,
                ["button", {"className": "path", "style":"left:'$x'px; top:'$y'px"}]]'
            fi
        fi
    done
}

app < <(stui -in /tmp/stuiIn) > /tmp/stuiIn