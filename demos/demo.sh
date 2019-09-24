#!/usr/bin/env bash
rm /tmp/stuiIn
mkfifo /tmp/stuiIn
rm /tmp/stuiOut
mkfifo /tmp/stuiOut
tail -f /tmp/stuiIn | stui > /tmp/stuiOut &

choice=false

while true
do
    if read msg </tmp/stuiOut; then
        # echo $msg >&2
        if [ '["hi"]' = "$msg" ]; then
            echo '["PostElem", "#app", -1, 
                ["label", "I can make a GUI in *bash*?!"]]'  > /tmp/stuiIn
            echo '["PostElem", "#app", -1,  
                ["input", {"type":"checkbox"}]]' > /tmp/stuiIn
            echo '["PostElem", "#app", -1,
                ["button", "Confirm"]]' > /tmp/stuiIn
            echo '["Subscribe", "input", "onchange", ["target.checked"]]' > /tmp/stuiIn
            echo '["Subscribe", "button", "onclick", []]' > /tmp/stuiIn
        elif [ '["bye"]' = "$msg" ]; then
            # echo "quitting..." >&2
            break
        else
            evType=$(echo $msg | jq .[2])
            
            if [ '"onclick"' = "$evType" ]; then
                echo $choice 
                echo '["Close"]' > /tmp/stuiIn
            elif [ '"onchange"' = "$evType" ]; then
                choice=$(echo $msg | jq .[3].\"target.checked\")
            fi
            
        fi
    fi
done
