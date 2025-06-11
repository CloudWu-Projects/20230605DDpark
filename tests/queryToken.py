import requests

qJson={
  "OperatorID": "123456789",
  "Data": "mYvffpNoFf4E/ZTC1tOw41TC5OlkEobfAYCm5N8hEusaLUaUIqOrXtdbMrSck0DSmfM7mRuOGMoCQzH0nWPGuw==",
  "TimeStamp": "20180120165755",
  "Seq": "0001",
  "Sig": "D2D584A14F3F284445DF85D0E8C0697C"
}

r = requests.post("http://127.0.0.1:8081/query_token", json=qJson)

print(r.text)


print("=======================================")
headers={
    'Authorization': 'Bearer FFDD0D8C66C7A6C1C537FD4D6F52CB15',
}
r= requests.post("http://127.0.0.1:8081/test",headers=headers, json=qJson)

print(r.text)