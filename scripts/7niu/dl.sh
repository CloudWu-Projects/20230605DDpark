export servername="yilingshequ_go"

echo $servername
systemctl stop $servername.service

python3 down7niu.py $servername/$servername.zip

mv $servername/$servername.zip $servername/$servername -f

chmod +x ./$servername/$servername
./$servername/$servername -install

md5sum ./$servername/$servername

systemctl start $servername.service