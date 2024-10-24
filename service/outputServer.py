import requests
import time
import json

from utils.configer import cfg
from utils.logutil import logger
from utils.sqlite_dao import SqliteUtil
Car_InOUT_Json = {
    "requestId": "NA",
    "recordSN": "停车云2.2.order_id",
    "actionType": "IN",
    "plateNum": "停车云2.2.Car_number",
    "plateColor": "BE",
    "brand": "",
    "plateType": "",
    "projectId": "",
    "projectName": "",
    "entryTime": "停车云2.2.in_time 停车云2.3.out_time",
    "timestamp": "now",
    "deviceSn": "",
    "deviceName": ""
}
PHOTO_JSON = {
    "requestId": "NA",
    "recordSN": "停车云2.2.order_id",
    "dataType": "URL",
    "photosData": " →停车云2.2.pic_addr、停车云2.3.pic_addr",
    "timestamp": " →当前时间，十位时间戳"
}
# SELECT id,orderid,picUrl FROM PICInfo p
SELECT_CAR_INFO_SQL = "select id,orderid,picUrl from `PICInfo` where orderid = ?"
INSERT_CAR_INFO_SQL = "insert into `PICInfo` (`orderid`, `picUrl`) VALUES (?,?)"
REMOVE_CAR_INFO_SQL = "delete from `PICInfo` where orderid = ?"


class OutputServer:
    def __init__(self):
        self.dao = SqliteUtil()

    def __convertColor(self, license_color):
        if license_color == 0:  # "蓝":
            return "BE"
        elif license_color == 1:  # "黄":
            return "YW"
        elif license_color == 3:  # "黑":
            return "BK"
        elif license_color == 2:  # "白":
            return "WE"
        elif license_color == 4:  # "绿":
            return "GN"
        else:
            return "OR"

    def __post_to_Server(self, url, outJson):
        try:
            logger.debug(f"__post_to_Server {url}")
            logger.debug(f"{json.dumps(outJson, ensure_ascii=False)}")
            ret = requests.post(url, json=outJson)
            logger.debug(ret.text)
        except Exception as e:
            logger.error("post_to_server error")
            logger.error(url)
            logger.error(outJson)
            logger.error(e)

    def __getProjectid(self,parkid):
        return cfg.getParkinfo(parkid)['projectId']
    def __getProjectName(self,parkid):
        return cfg.getParkinfo(parkid)['projectName']
    
    def car_in(self, park_id, order_id, car_number, license_color, in_time, pic_addr):
        logger.debug(f"car_in {park_id} {order_id} {car_number} {license_color} {in_time} {pic_addr}")
        outJson = Car_InOUT_Json.copy()
        outJson['recordSN'] = order_id
        outJson['plateNum'] = car_number
        outJson['entryTime'] = in_time
        outJson['timestamp'] = int(time.time())
        outJson['plateColor'] = self.__convertColor(license_color)
        outJson['projectId'] =self.__getProjectid(park_id)
        outJson['projectName'] = self.__getProjectName(park_id)

        self.__insert_car_info(f"{park_id}-{order_id}", pic_addr)

        self.__post_to_Server(cfg.get("platform", "url_inout"), outJson)
        self.photo(order_id, pic_addr,None)

    def car_out(self, park_id, order_id, car_number, license_color, out_time, pic_addr):
        logger.debug(f"car_out {park_id} {order_id} {car_number} {license_color} {out_time} {pic_addr}")
        outJson = Car_InOUT_Json.copy()
        outJson['recordSN'] = order_id
        outJson['plateNum'] = car_number
        outJson['entryTime'] = out_time
        outJson['actionType'] = "OU"
        outJson['timestamp'] = int(time.time())
        outJson['plateColor'] = self.__convertColor(license_color)
        dbOrderID = f"{park_id}-{order_id}"
        in_picUrl = self.__query(dbOrderID)
        if in_picUrl is None:
            logger.error(f"不能找到 order_id {dbOrderID} 的记录")
            return
        self.__post_to_Server(cfg.get("platform", "url_inout"), outJson)
        self.photo(order_id, in_picUrl, pic_addr)

    def __insert_car_info(self, order_id, picurl):
        logger.debug(f"insert {order_id} {picurl}")
        self.dao.updateBySql(INSERT_CAR_INFO_SQL, (order_id, picurl))

    def __query(self, order_id):
        try:
            data = self.dao.getOneBySql(SELECT_CAR_INFO_SQL, (order_id,))
        except Exception as e:
            logger.error("查询车牌【%s】数据异常：%s", order_id, e)
            raise e

        if data:
           # return {"oss":data["oss"],"pfjd":data["pfjd"],"order_id":order_id}
           return data['picUrl']
        return None

    def __remove(self, order_id):
        try:
            data = self.dao.getOneBySql(REMOVE_CAR_INFO_SQL, (order_id,))
        except Exception as e:
            logger.error("删除【%s】数据异常：%s", order_id, e)

    def photo(self, order_id, in_picUrl, out_pic_addr):

        outJson = PHOTO_JSON.copy()

        outJson['recordSN'] = order_id
        outJson['timestamp'] = int(time.time())
        outJson['photosData'] =""
        if in_picUrl is not None:
            outJson['photosData'] = f"{in_picUrl}\n"        
        
        if out_pic_addr is not None:
            outJson['photosData'] += f"{out_pic_addr}\n"

        self.__post_to_Server(cfg.get("platform", "url_pic"), outJson)
        self.__remove(order_id)
