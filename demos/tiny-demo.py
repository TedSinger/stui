import zmq
import json
import subprocess

g = subprocess.Popen(['guise'],
    stdin=subprocess.PIPE,
    stdout=subprocess.PIPE)

def send(command):
    g.stdin.write(bytes(json.dumps(command), 'utf-8'))
    g.stdin.flush()

are_we_gui = False

while True:
    event = json.loads(g.stdout.readline())
    if event[0] == "hi":
        send(
            ["PostElem", "#app", -1, 
                ["div", [
                    ["label", "Are we GUI?"], 
                    ["input", {"type":"checkbox", "className":"foo"}],
                    ["button", "Confirm"]]]])
        send(["Subscribe", ".foo", "onchange", ["target.checked"]])
        send(["Subscribe", "button", "onclick", []])
    elif event[0] == "bye":
        break
    elif event[1] == ".foo":
        are_we_gui = event[3]["target.checked"]
    elif event[1] == "button":
        send(["Close"])
        break

print(are_we_gui)