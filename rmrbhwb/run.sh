#go get crawler.club/bcrawler
CURR_PATH=`cd $(dirname $0);pwd;`
cd $CURR_PATH
ts=`date +%Y-%m/%d`
url="http://paper.people.com.cn/rmrbhwb/html/$ts/node_865.htm"
echo $url
rm first.lock
bcrawler -start day -start_url $url -log_dir ./log -alsologtostderr
