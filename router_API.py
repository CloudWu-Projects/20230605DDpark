from flask import Blueprint

router_API = Blueprint('router_API', __name__)



@router_API.route('/outpark', methods=['POST','GET'])  
@router_API.route('/inpark', methods=['POST','GET']) 
def out_in_park():
    return 'router_API'