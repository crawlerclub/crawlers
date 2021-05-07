import requests
import re

link = "https://yjrx.bjedu.cn/json/%s.json"
urls = [
    "https://yjrx.bjedu.cn/portal_dc/primary.htm",
    "https://yjrx.bjedu.cn/portal_xc/primary.htm",
    "https://yjrx.bjedu.cn/portal_cy/primary.htm",
    "https://yjrx.bjedu.cn/portal_hd/primary.htm",
    "https://yjrx.bjedu.cn/portal_ft/primary.htm",
    "https://yjrx.bjedu.cn/portal_sjs/primary.htm",
    "https://yjrx.bjedu.cn/portal_mtg/primary.htm",
    "https://yjrx.bjedu.cn/portal_fs/primary.htm",
    "https://yjrx.bjedu.cn/portal_tz/primary.htm",
    "https://yjrx.bjedu.cn/portal_sy/primary.htm",
    "https://yjrx.bjedu.cn/portal_cp/primary.htm",
    "https://yjrx.bjedu.cn/portal_dx/primary.htm",
    "https://yjrx.bjedu.cn/portal_pg/primary.htm",
    "https://yjrx.bjedu.cn/portal_hr/primary.htm",
    "https://yjrx.bjedu.cn/portal_my/primary.htm",
    "https://yjrx.bjedu.cn/portal_yq/primary.htm",
    "https://yjrx.bjedu.cn/portal_ys/primary.htm",
    "https://yjrx.bjedu.cn/portal_yz/primary.htm",
]
list_file = open("primarylist.json", "w")
for url in urls:
    ret = requests.get(url)
    fid = re.findall(r"initPaginator\('(.+?)'\);", ret.text)    
    if len(fid) == 1:
        u = link % fid[0]
        r = requests.get(u)
        list_file.write(r.text+"\n")
        print(u)
list_file.close()

