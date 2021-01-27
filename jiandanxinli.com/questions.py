import requests
from retrying import retry
from lsm import LSM
db = LSM("questions.ldb")

@retry(stop_max_attempt_number=7)
def get(url):
    ret = requests.get(url)
    db[i] = ret.text

for i in range(1, 371365):
    if i not in db:
        url = "https://www.jiandanxinli.com/questions/%d" % i
        print(url)
        get(url)
