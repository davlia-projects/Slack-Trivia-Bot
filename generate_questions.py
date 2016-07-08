import json


def banned(banned, phrase):
    for word in banned:
        for p in phrase:
            if word in p:
                return True
    return False

bl = ["hm", "hehe", "ha ha", "he he", "tower", "barracks", "ancient", "first blood"]
f = open('r.json')
g = open('s.json', 'w')
content = json.load(f)
added = {}
skip = False
for c in content:
    k = c['name'].split('/')[0]
    if len(k.split(" ")) > 2:
        continue
    responses = c['responses']
    added_resp = []
    for resp in responses:
        if len(resp.split(" ")) > 12 and not banned(bl, resp.lower()):
            print k , resp, len(resp.split(" "))
            added_resp.append(resp)
    added[k] = added_resp


added.pop("Warlock's Golem")
