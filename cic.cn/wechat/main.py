import requests
import json
list_url = "https://zhebwx.cic.cn/product/load/list"
post_data = {"id":"YTHSHH001","idType":"USERID","index":1,"productType":"","rows":100,"status":0}
ret = requests.post(list_url, json=post_data)
with open("product_list.txt", "w") as out:
    out.write(ret.text)
data = json.loads(ret.text)
with open("products.txt", "w") as out:
    for item in data['data']['productList']['data']:
        url = "https://zhebwx.cic.cn/product/load/instance/%s" % item['cCode']
        print(url)
        out.write(requests.get(url).text+"\n")
