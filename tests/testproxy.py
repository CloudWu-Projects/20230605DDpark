import requests


ChargeInfo={
          "appId":"007B45C733038000",
	"startTime": "2025-07-24 00:00:00",
	"endTime": "2025-07-25 00:00:00"

}


r= requests.post("http://127.0.0.1:9090/proxyHandler",json=ChargeInfo)
print(r.json())