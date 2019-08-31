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

disabled = False
while True:
    event = eventSocket.recv_json()
    
    if event[0] == "hi":
        commandSocket.send_json(["SetHtml", "#app", "<button>hi there!</button><textarea>type here...</textarea><label></label>"])
        commandSocket.send_json(["Sub", "button", "onclick", ["x", "y"]])
        commandSocket.send_json(["Sub", "textarea", "onkeyup", ["target.value"]])
        commandSocket.send_json(["Sub", "textarea", "onmousemove", ["x"]])
        commandSocket.send_json(["SetCss", "button", {"transition": "background-color 2s"}])
        commandSocket.send_json(["SetCss", "#app", {"max-width": "100%"}])
    elif event[0] == "bye":
        break
    elif event[1] == "button" and event[2] == "onclick":
        commandSocket.send_json(["SetCss", "button", {"background-color": get_color()}])
    elif event[1] == "textarea" and event[2] == "onkeyup":
        text = event[3]['target.value']
        if 'disable' in text and not disabled:
            commandSocket.send_json(['SetAttrs', 'button', {'disabled':True}])
            disabled = True
        elif disabled:
            commandSocket.send_json(['SetAttrs', 'button', {'disabled':False}])
            disabled = False
        commandSocket.send_json(["SetHtml", "label", text])
        commandSocket.send_json(["SetCss", "button", {"background-color": "revert"}])
    elif event[1] == "textarea" and event[2] == "onmousemove":
        commandSocket.send_json(['SetCss', 'textarea', {'font-size':str(int(event[3]["x"] / 10)) + 'px'}])
    print("demo: " + str(event))