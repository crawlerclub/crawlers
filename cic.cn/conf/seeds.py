network = [{"parser_name": "network", "url": "http://www.cic.cn/network/index_%d.jhtml" % x} for x in range(1, 5)]
news = [{"parser_name": "news", "url": "http://www.cic.cn/companyNews/index_%d.jhtml" % x} for x in range(1, 23)]
reports = [{"parser_name": "reports", "url": "http://www.cic.cn/mediaReports/index_%d.jhtml" % x} for x in range(1, 13)]
x708 = [{"parser_name": "708", "url": "http://www.cic.cn/708/index_%d.jhtml" % x} for x in range(1, 3)]
seeds = network + news + reports + x708
import json
ret = json.dumps(seeds, indent=2)
print(ret)
