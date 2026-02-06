export serverName=hbproxy
systemctl stop $serverName.service

python3 down7niu.py $serverName/$serverName.zip

mv $serverName/$serverName.zip $serverName/$serverName -f

chmod +x ./$serverName/$serverName
./$serverName/$serverName -install

md5sum ./$serverName/$serverName

systemctl start $serverName.service