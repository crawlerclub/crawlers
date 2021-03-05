import datetime
url_ = "http://paper.people.com.cn/rmrb/html/%s/nbs.D110000renmrb_01.htm"
begin = datetime.date(2020, 1, 1)
end = datetime.date(2021, 3, 4)
d = begin
delta = datetime.timedelta(days=1)
urls = []
while d <= end:
    day = d.strftime("%Y-%m/%d")
    url = url_ % day
    d += delta
    urls.append(url)
seeds = [{"parser_name": "day", "url": x} for x in urls]
import json
ret = json.dumps(seeds, indent=2)
print(ret)
