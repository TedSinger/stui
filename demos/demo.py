import zmq
import json
import subprocess
import random

g = subprocess.Popen(['guise'], stdout=subprocess.PIPE)
sockets = json.load(g.stdout)
c = zmq.Context()
eventSocket = c.socket(zmq.PULL)
eventSocket.connect(sockets['events'])
commandSocket = c.socket(zmq.PUSH)
commandSocket.bind(sockets['commands'])

def get_color():
    return 'rgb({},{},{})'.format(int(random.random()*256),int(random.random()*256),int(random.random()*256))

send = commandSocket.send_json

disabled = False
while True:
    event = eventSocket.recv_json()
    
    if event[0] == "hi":
        send(["PostHtml", "#app", -1, ["button", {"textContent": "hello"}, []]])
        send(["PostHtml", "#app", -1, ["textarea", {}, []]])
        send(["PostHtml", "#app", 1, ["label", {}, []]])
        send(["Sub", "button", "onclick", ["x", "y"]])
        send(["Sub", "textarea", "onkeyup", ["target.value"]])
        send(["Sub", "textarea", "onmousemove", ["x"]])
        send(["PatchCss", "button", {"transition": "background-color 2s"}])
        send(["PatchCss", "#app", {"max-width": "100%"}])
    elif event[0] == "bye":
        break
    elif event[1] == "button" and event[2] == "onclick":
        send(["PatchCss", "button", {"background-color": get_color()}])
        if random.random() > 0.9:
            send(["DeleteHtml", "label"])
    elif event[1] == "textarea" and event[2] == "onkeyup":
        text = event[3]['target.value']
        if 'disable' in text and not disabled:
            send(['PatchAttrs', 'button', {'disabled':True}])
            disabled = True
        elif 'disable' not in text and disabled:
            send(['PatchAttrs', 'button', {'disabled':False}])
            disabled = False
        send(["PatchAttrs", "label", {"textContent":text}])
        send(["PatchCss", "button", {"background-color": "revert"}])
    elif event[1] == "textarea" and event[2] == "onmousemove":
        send(['PatchCss', 'textarea', {'font-size':str(int(event[3]["x"] / 10)) + 'px'}])
    # print("demo: " + str(event))