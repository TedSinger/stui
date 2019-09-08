import zmq
import subprocess
import random

g = subprocess.Popen(['stui', '-conn', 'zmq'])
c = zmq.Context()
sock = c.socket(zmq.PAIR)
sock.connect("ipc:///tmp/stui")

def get_color():
    return 'rgb({},{},{})'.format(int(random.random()*256),int(random.random()*256),int(random.random()*256))

send = sock.send_json

disabled = False
while True:
    event = sock.recv_json()
    
    if event[0] == "hi":
        send(["PostElem", "#app", -1, ["button", {"textContent": "hello", "className":"foo"}, []]])
        send(["PostElem", "#app", -1, ["textarea", {}, []]])
        send(["PostElem", "#app", 1, ["label", {}, []]])
        send(["Subscribe", "button", "onclick", ["x", "y"]])
        send(["Subscribe", "textarea", "onkeyup", ["target.value"]])
        send(["Subscribe", "textarea", "onmousemove", ["x"]])
        send(["Subscribe", ".foo", "onmousemove", ["x", "y"]])
        send(["PatchStyles", "button", {"transition": "background-color 2s"}])
        send(["PatchStyles", "#app", {"max-width": "100%", "width":"100%", "height":"600px"}])
    elif event[0] == "bye":
        break
    elif event[1] == "button" and event[2] == "onclick":
        send(["PatchStyles", "button", {"background-color": get_color()}])
    elif event[1] == ".foo" and event[2] == "onmousemove":
        send(["PatchStyles", ".foo", {"position": "absolute", "left": str(int(event[3]["x"]-20)) + "px", "top": str(int(event[3]["y"])-20)+ "px"}])
    elif event[1] == "textarea" and event[2] == "onkeyup":
        text = event[3]['target.value']
        if 'disable' in text and not disabled:
            send(['PatchAttrs', 'button', {'disabled':True}])
            disabled = True
        elif 'disable' not in text and disabled:
            send(['PatchAttrs', 'button', {'disabled':False}])
            disabled = False
        send(["PatchAttrs", "label", {"textContent":text}])
        send(["PatchStyles", "button", {"background-color": "revert"}])
    elif event[1] == "textarea" and event[2] == "onmousemove":
        send(['PatchStyles', 'textarea', {'font-size':str(int(event[3]["x"] / 10)) + 'px'}])
    else:
        print(event)
    # print("demo: " + str(event))