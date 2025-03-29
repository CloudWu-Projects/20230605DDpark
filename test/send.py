import requests


ChargeInfo={
    "parkId":"100515570",
"plateNo":"川ab34556"
}


r= requests.post("http://shuyun.ddpark.fun/chargePile/chargingRecord",json=ChargeInfo)
print(r.json())