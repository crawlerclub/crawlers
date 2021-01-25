import requests
from lsm import LSM
db = LSM("questions.ldb")
for i in range(1, 371365):
    if i not in db:
        url = "https://www.jiandanxinli.com/questions/%d" % i
        print(url)
        ret = requests.get(url)
        db[i] = ret.text
