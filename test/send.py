import requests


ChargeInfo={
    "parkId":"1234565",
"plateNo":"川ab34556"
}


r= requests.post("http://127.0.0.1:8080/chargePile/chargingRecord",json=ChargeInfo)
print(r.json())