import requests
import json
import hashlib
import json
queryJson={
    "service_name":"query_order",
    "sign":"9BBD511E3178A91E66AD4C6A043CEC55",
    "park_id":10051557,
    "data":{  
        "car_number":"京A5566TT"          
        }
    }

ukey="7JPWIA1SGV9N17LE"

def makeSign(payload):
    dataS =json.dumps( payload['data'] , separators=(',', ':'),ensure_ascii=False)
    data = dataS + "key=" + ukey

    print("makeSign",data)
    newSign=hashlib.md5(data.encode('utf-8')).hexdigest().upper()
    print("makeSign",newSign)
    return newSign

def queryOrder(dataJson):
    print("\n\n")
    if dataJson is not None:
        queryJson['data']=dataJson
        sign = makeSign(queryJson)
        queryJson['sign']=sign
    print("queryJson",json.dumps(queryJson,ensure_ascii=False))
    url= 'http://istparking.sciseetech.com/public/order/queryOrder'
    print("url",url)

    r = requests.post(url,json=queryJson)
    print(r.json())

queryOrder(None)
queryOrder({"car_number":"川A12345"})
queryOrder({"car_number":"京A5566TT"})

