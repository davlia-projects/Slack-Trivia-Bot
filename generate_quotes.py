import json


def banned(banned, phrase):
    for word in banned:
        for p in phrase:
            if word in p:
                return True
    return False

bl = ["hm", "hehe", "ha", "heh", "tower", "barracks", "ancient", "first", "bag"]
f = open('data/r.json')
g = open('data/s.json', 'w')
content = json.load(f)
content[-1]['responses'] = content[-1]['responses'][::2]
added = {}
skip = False
for c in content:
    k = c['name'].split('/')[0]
    if len(k.split(" ")) > 2:
        continue
    responses = c['responses']
    added_resp = []
    for resp in responses:
        tokens = set([x.lower() for x in resp.split(" ")])
        tokens.add(k.lower())
        if len(tokens) > 10 and not tokens.intersection(set(bl)) and not k.lower() in resp.lower():
            print k , resp
            added_resp.append(resp)
    added[k] = list(set(added_resp))

topop = ["Warlock's Golem", "Announcer", "Portal Pack", "Shopkeeper"]
for t in topop:
    added.pop(t)

o = open('data/questions.json', 'a')
o.write('[')
id = 0
for k,v in added.items():
    for quote in v:
        out = '{\n\t"id": "%s",\n\t"prompt": "Dota Hero by quote: %s", \n\t"answer": "%s"\n},\n' % (id, quote.replace('\"', "'"), k)
        id += 1
        try:
            o.write(out)
        except:
            continue
o.write(']')
