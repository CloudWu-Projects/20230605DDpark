# -*- coding: utf-8 -*-

from flask import Flask
from flask import request
from utils.configer import cfg
from conf.config import  make_error_data,make_success_data
from router_API import router_API

app = Flask(__name__)

app.register_blueprint(router_API)

@app.route("/api/get_car_data",methods=['POST'])
def get_car_data():
    data = request.get_json()
    if "park_id" not in data:
        return make_error_data("缺少停车场编号参数")
    if "car_number" not in data:
        return make_error_data("缺少车牌号参数")
    if "vin" not in data:
        return make_error_data("缺少vin参数")
    if "motor" not in data:
        return make_error_data("缺少发动机号参数")
    return make_success_data("OK")

@app.errorhandler(404)
def page_not_found(error):
    return "不接受该访问",404

@app.errorhandler(400)
def param_error(error):
    return "不接受该访问",400

if __name__ == "__main__":
    port = cfg.getInt("server","port")
    # 启动服务
    # http://127.0.0.1:8080/api/get_car_data
    app.run(host='0.0.0.0', port=port)
# uvicorn boot:app --reload
# sqlite 的图形化管理工具
# sqlite_web --password -x -p 8088 data/datas_db.db
