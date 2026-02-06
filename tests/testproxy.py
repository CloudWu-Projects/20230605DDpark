import requests


ChargeInfo={
       "discountNumber": "169157199496923925786",
       "plateNumber": "浙A6Y557",
       "discountTime": 120,
       "timestamp": 1691573927454,
       "sign": "9B2643F390139D67DB98F6D8CE34B6A2",
       "platformNo": "test001"
}


r= requests.post("http://127.0.0.1:9090/api/wec/XJCdiscount",json=ChargeInfo)
print(r.json())