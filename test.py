
import requests
from service import outputServer

outpicAddr = "https://ist-falcon.oss-cn-shenzhen.aliyuncs.com/order-images/48446/out/48203904.jpg?Expires=1723557992&OSSAccessKeyId=LTAIQQrl6GICP0QX&Signature=xcUsiAYQ%2F1PX9M5%2Fe4BtJcpL4JI%3D"


def test():
    os = outputServer.OutputServer()
    os.car_in('48207302', '川ADG8621', 'license_color',
              '1723536379', outpicAddr)
    os.car_out('48207302', '川ADG8621', 'license_color',
               '1723536379', outpicAddr)

    os.photo('48204902', "aaaa")
# test()


ajson = {"requestId": "NA",
         "recordSN": "123456",
         "actionType": "IN",
         "plateNum": "car_number",
         "plateColor": "OR",
         "brand": "",
         "plateType": "",
         "projectId": "projectId",
         "projectName": "projectName",
         "entryTime": "1723602503",
         "timestamp": "1723602503",
         "deviceSn": "",
         "deviceName": ""}
url = "http://222.92.40.243:8066/mgr/data/access/car/project/record"

# x = requests.post(url, json=ajson)
# print(x.text)

# 2024-08-14 10:28:23,220 -__post_to_Server - outputServer.py[line:54] - ERROR: {'requestId': 'NA', 'recordSN': '123456', 'actionType': 'IN', 'plateNum': 'car_number', 'plateColor': 'OR', 'brand': '', 'plateType': '', 'projectId': 'projectId', 'projectName': 'projectName', 'entryTime': 'in_time', 'timestamp': 1723602503174322688, 'deviceSn': '', 'deviceName': ''}
# 2024-08-14 10:28:23,307 -__post_to_Server - outputServer.py[line:56] - DEBUG: {"code":10001,"message":"系统异常，请联系管理员！","success":false}

photo_json = {
    "requestId": "NA",
    "recordSN": "48232102",
    "dataType": "URL",
    "photosData": "https://ist-falcon.oss-cn-shenzhen.aliyuncs.com/order-images/48446/in/48232102.jpg?Expires=1723628781&OSSAccessKeyId=LTAIQQrl6GICP0QX&Signature=oBL0hV7CumVW%2FlzLx7rvGiGqpqI%3D\nhttp://ist-falcon.oss-cn-shenzhen.aliyuncs.com/order-images/48446/out/48232102.jpg?Expires=1723633749&OSSAccessKeyId=LTAIQQrl6GICP0QX&Signature=3HBQS42W7v1Z2dCvB%2BunC9xB1hE%3D\n",
    "timestamp": 1723612146
}
url = "http://222.92.40.243:8066/mgr/data/access/car/project/photos"

x = requests.post(url, json=photo_json)
print(x.text)
