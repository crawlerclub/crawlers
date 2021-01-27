import requests
from lsm import LSM
db = LSM("questions.ldb")
fail = 0
#for i in range(1, 100000):
#for i in range(100000, 200000):
for i in range(602000, 639317):
    if i not in db:
        url = "https://www.ydl.com/ask/%d" % i
        print(url)
        ret = requests.get(url)
        if ret.status_code == 200:
            db[i] = ret.text
            fail = 0
        else:
            print(ret.status_code)
            fail += 1
            if fail > 100 and fail % 10 == 0:
                print("fail=%d" % fail)
            db[i] = "404"
