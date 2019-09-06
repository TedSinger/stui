#!/usr/bin/env bash
rm /tmp/guiseIn
mkfifo /tmp/guiseIn
rm /tmp/guiseOut
mkfifo /tmp/guiseOut
tail -fz /tmp/guiseIn | guise > /tmp/guiseOut &

choice=false

while true
do
    if read msg </tmp/guiseOut; then
        echo $msg >&2
        if [ '["hi"]' = "$msg" ]; then
            echo '["PostElem", "#app", -1, 
                ["label", {"textContent":"I can make a GUI in *bash*?!"},[]]]'  > /tmp/guiseIn
            echo '["PostElem", "#app", -1,  
                ["input", {"type":"checkbox"}, []]]' > /tmp/guiseIn
            echo '["PostElem", "#app", -1,
                ["button", {"textContent": "Confirm"}, []]]' > /tmp/guiseIn
            echo '["Subscribe", "input", "onchange", ["target.checked"]]' > /tmp/guiseIn
            echo '["Subscribe", "button", "onclick", []]' > /tmp/guiseIn
        elif [ '["bye"]' = "$msg" ]; then
            echo "quitting..." >&2
            break
        else
            evType=$(echo $msg | jq .[2])
            
            if [ '"onclick"' = "$evType" ]; then
                echo $choice 
                echo '["Close"]' > /tmp/guiseIn
            elif [ '"onchange"' = "$evType" ]; then
                choice=$(echo $msg | jq .[3].\"target.checked\")
            fi
            
        fi
    fi
done