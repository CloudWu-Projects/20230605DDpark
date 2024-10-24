from flask import Blueprint, request,Response,jsonify
from utils.logutil import logger
import json
from utils.configer import cfg

router_API = Blueprint('router_API', __name__)

from service import outputServer

outputServer= outputServer.OutputServer()



@router_API.route('/outpark', methods=['POST','GET'])  
def out_park():
    json_body = request.json
    logger.debug(f"{json.dumps(json_body,ensure_ascii=False)}")  

    park_id=json_body['park_id']
    if f'{park_id}' not in cfg.GetConfig()['parkinfo']:
        logger.error(f'park_id {park_id} not allow>>>')
        return jsonify({"state":0,f"errmsg":"park_id {park_id} not allow"})           
    
    
    order_id=json_body['data']['order_id']
    car_number=json_body['data']['car_number']
    license_color=json_body['data']['license_color']
    out_time=json_body['data']['out_time']
    pic_addr=json_body['data']['pic_addr']

    try:
        outputServer.car_out(park_id,order_id,car_number,license_color,out_time,pic_addr)
    except Exception as e:
        logger.error(e)
        
    ret = {"state": 1, "errmsg": "OK"}
    reuslt = {
        "state": 1,
        "order_id": json_body['data']['order_id'],
        "park_id": park_id,
        "service_name": json_body['service_name'],
        "errmsg": json.dumps(ret, ensure_ascii=False)
    }
    return jsonify(reuslt)
    

    
@router_API.route('/inpark', methods=['POST','GET']) 
def in_park():
    json_body = request.json
    logger.debug(f"{json.dumps(json_body,ensure_ascii=False)}") 
    park_id=json_body['park_id']
    b = cfg.getConfig()['parkinfo']
    if f'{park_id}'  not in cfg.getConfig()['parkinfo']:
        logger.error(f'park_id {park_id} not allow>>>')
        return jsonify({"state":0,"errmsg":"park_id error"})           
    logger.error(f"park_id {park_id} ")
    order_id=json_body['data']['order_id']
    car_number=json_body['data']['car_number']
    license_color=json_body['data']['license_color']
    in_time=json_body['data']['in_time']
    pic_addr=json_body['data']['pic_addr']

    try:
        outputServer.car_in(park_id,order_id,car_number,license_color,in_time,pic_addr)
    except Exception as e:
        logger.error(e)
    ret = {"state": 1, "errmsg": "OK"}
    reuslt = {
        "state": 1,
        "order_id": json_body['data']['order_id'],
        "park_id": park_id,
        "service_name": json_body['service_name'],
        "errmsg": json.dumps(ret, ensure_ascii=False)
    }
    return jsonify(reuslt)
    