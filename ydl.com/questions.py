import requests
from lsm import LSM
db = LSM("questions.ldb")
#for i in range(1, 100000):
for i in range(100000, 200000):
    if i not in db:
        url = "https://www.ydl.com/ask/%d" % i
        print(url)
        ret = requests.get(url)
        if ret.status_code == 200:
            db[i] = ret.text
        else:
            print(ret.status_code)
            db[i] = "404"
