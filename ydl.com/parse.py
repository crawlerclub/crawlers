from tqdm import tqdm
from lsm import LSM
import sys, os
etpath = os.path.join(os.path.abspath("."), "lib")
sys.path.append(etpath)
import et

db = LSM("questions.ldb")

def main():
    out = open("ret.json", "w")
    err = open("err.txt", "w")
    non = open("non.txt", "w")
    cnt = 0
    for i in tqdm(range(1, 639317)):
        if i not in db:
            non.write("%d\n" % i)
            continue
        page = db[i]
        if page != b'404':
            url = "https://www.ydl.com/ask/%d" % i
            try:
                ret = et.ParseExt("ydl.json", url, page.decode('utf-8'))
                out.write("%s\n" % ret)
            except Exception as e:
                print(e)
                err.write("%s\n" % url)
            cnt += 1
    print(cnt)
    out.close()
    err.close()
    non.close()


if __name__ == "__main__":
    main()
