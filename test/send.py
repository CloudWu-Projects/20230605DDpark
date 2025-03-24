import requests


ChargeInfo={
    "parkId":"1234565",
"plateNo":"川ab34556"
}


r= requests.post("http://shuyun.ddpark.fun/chargePile/chargingRecord",json=ChargeInfo)
print(r.json())