import requests

photo_json = {
    "requestId": "NA",
    "recordSN": "48232102",
    "dataType": "URL",
    "photosData": "",
    "timestamp": 1723612146
}
url = "http://222.92.40.243:8066/mgr/data/access/car/project/photos"



picArr={
    "https://ist-falcon.oss-cn-shenzhen.aliyuncs.com/order-images/48446/in/48232102.jpg?Expires=1723628781&OSSAccessKeyId=LTAIQQrl6GICP0QX&Signature=oBL0hV7CumVW%2FlzLx7rvGiGqpqI%3D",
    "http://ist-falcon.oss-cn-shenzhen.aliyuncs.com/order-images/48446/out/48232102.jpg?Expires=1723633749&OSSAccessKeyId=LTAIQQrl6GICP0QX&Signature=3HBQS42W7v1Z2dCvB%2BunC9xB1hE%3D"
}

allPhoto=""
for i in picArr:
    print("i--",i)
    curPhoto=f"{i}\n"
    allPhoto += curPhoto
    print("allPhoto--",allPhoto)
    
    photo_json["photosData"]=curPhoto
    x = requests.post(url, json=photo_json)
    print(x.text)
    
    photo_json["photosData"]=allPhoto
    x = requests.post(url, json=photo_json)
    print(x.text)
