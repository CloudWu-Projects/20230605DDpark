import requests

jsondataStr="{\"appId\":\"cdist\",\"key\":\"FD2A9F67E2E5B30FDB40437AF84477D6\",\"parkId\":10053274,\"serviceCode\":\"SyncChargePilePay\",\"ts\":\"1753398891322\",\"reqId\":\"1753398891322\",\"orderNo\":\"1753398891322\",\"plateNo\":\"苏ADS9990\",\"startTime\":\"2020-01-19 14:23:11\",\"endTime\":\"2020-01-19 15:23:11\",\"stationId\":\"asdad\",\"stationName\":\"测试场站\",\"deviceId\":\"asdas\",\"deviceName\":\"dasd\",\"spaceNo\":\"1\",\"power\":1,\"elecMoney\":1,\"seviceMoney\":1,\"totalMoney\":1,\"freeType\":0,\"freeMoney\":0,\"freeTime\":7200}"

jsondata=eval(jsondataStr)
print(jsondata)

r = requests.post("http://127.0.0.1:9090/api/wec/SyncChargePilePay", json=jsondata)
print(r.text)