import requests
import json
import sys
sys.path.append(".")
import et

out = open("content.txt", "w")
for line in open("primarylist.json"):
    items = json.loads(line)
    for item in items:
        url = "https://yjrx.bjedu.cn/" + item['accessPath']
        ret = requests.get(url)
        print(url)
        r = et.ParseExt("school.json", url, ret.text)
        out.write("%s\n" % r)
out.close()

