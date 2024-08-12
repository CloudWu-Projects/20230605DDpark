from flask import Blueprint, request,Response,jsonify
from utils.logutil import logger
import json

router_API = Blueprint('router_API', __name__)



@router_API.route('/outpark', methods=['POST','GET'])  
@router_API.route('/inpark', methods=['POST','GET']) 
def out_in_park():
    json_body = request.json
    
    logger.debug(f'{request.path},>>>{json.dumps(json_body,ensure_ascii=False)}')
    park_id=json_body['park_id']    
    '''if park_id not in global_LastInfo['parkinfo']:
        global_LastInfo['parkinfo'][park_id]={}
        
    global_LastInfo['parkinfo'][park_id]['car_number']=json_body['data']['car_number']
    global_LastInfo['parkinfo'][park_id]['lastrecv']=json_body
    global_LastInfo['parkinfo'][park_id]['lastUPdateTime']=time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())
    #1：出 0：进
    ]]
    '''
    if request.path == '/inpark' :
        inOutType=0
        device_name="大门入口"
        device_id = json_body['data']['in_channel_id']
    elif request.path == '/outpark' :
        inOutType=1
        device_name="大门出口"
        device_id = json_body['data']['out_channel_id']
    '''
    if str(park_id) not in config.parkInfo    :
        return jsonify({"state":0,"errmsg":"park_id error"})
    
    companyId = config.parkInfo[str(park_id)]['companyId']
    ret = postTo(companyId,inOutType,json_body['data']['car_number'],device_id,device_name)
'''
    ret={"state":1,"errmsg":"OK"}
    reuslt={
  "state": 1,
  "order_id": json_body['data']['order_id'],
  "park_id": park_id,
  "service_name": json_body['service_name'],
  "errmsg": json.dumps(ret,ensure_ascii=False)
}
    return jsonify(reuslt)
    