from flask import Blueprint, request,Response,jsonify
from utils.logutil import logger
import json
from utils.configer import cfg

config_API = Blueprint('config_API', __name__)



@config_API.route('/config', methods=['POST','GET'])  
def config():
    return jsonify(cfg)