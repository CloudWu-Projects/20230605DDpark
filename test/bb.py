
import hashlib
import json
s='{"car_number":"川B11125","service_name":"query_price","order_id":"order_1C62AD23B57DAD7B76DC8563_1590119900","park_id":"10051557","pay_scene":0,"query_order_no":"2184020180411145518-263"}key=7JPWIA1SGV9N17LE'
md5_hash = hashlib.md5(s.encode('utf-8')).hexdigest().upper()
print("原始字符串:\n",s)
print("md5_hash: ",md5_hash)
print("需要的sign: 2F7619E3962B4A4597AAB3270909FEC5")

s='{"query_time":1672129283,"car_number":"京A5566TT"}key=7JPWIA1SGV9N17LE'
md5_hash = hashlib.md5(s.encode('utf-8')).hexdigest().upper()
print("\n原始字符串: ",s)
print("md5_hash: ",md5_hash)
print("需要的sign: 9BBD511E3178A91E66AD4C6A043CEC55")

s='{"car_number":"京A5566TT","query_time":1672129283}key=7JPWIA1SGV9N17LE'
md5_hash = hashlib.md5(s.encode('utf-8')).hexdigest().upper()
print("\n原始字符串: ",s)
print("md5_hash: ",md5_hash)
print("需要的sign: 9BBD511E3178A91E66AD4C6A043CEC55")
