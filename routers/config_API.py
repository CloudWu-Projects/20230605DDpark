from flask import Blueprint, request,Response,jsonify,render_template_string
from utils.logutil import logger
import json
from utils.configer import cfg

config_API = Blueprint('config_API', __name__)



@config_API.route('/config', methods=['POST','GET'])  
def config():
    return jsonify(cfg.getConfig())


@config_API.route('/log', methods=['POST','GET'])  
def log():
    file_path = cfg.get("logging", "logPath")
    with open(file_path,"r",encoding="utf-8") as f:
        file_content = f.read()
    html="<html><body><pre>{}</pre></body></html>".format(file_content)
    return render_template_string(html)