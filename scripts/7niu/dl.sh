systemctl stop yianqi_go.service

python3 down7niu.py yianqi_go/yianqi_go.zip

mv yianqi_go/yianqi_go.zip yianqi_go/yianqi_go -f

chmod +x ./yianqi_go/yianqi_go
./yianqi_go/yianqi_go -install

md5sum ./yianqi_go/yianqi_go

systemctl start yianqi_go.service