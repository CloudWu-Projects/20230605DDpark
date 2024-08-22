# -*- coding: utf-8 -*-

from flask import Flask
from flask import request
from flask.json.provider import DefaultJSONProvider
from utils.configer import cfg
from conf.config import  make_error_data,make_success_data
from routers.router_API import router_API
from routers.config_API import config_API

app = Flask(__name__)

class CustomJSONProvider(DefaultJSONProvider):
    def dumps(self, obj, **kwargs):
        # 关闭 ensure_ascii 以支持中文
        kwargs['ensure_ascii'] = False
        return super().dumps(obj, **kwargs)
    
app.json = CustomJSONProvider(app)

app.register_blueprint(router_API)
app.register_blueprint(config_API)


@app.route("/",methods=['GET'])
def get_all_router():
    routers = [rule.rule for rule in app.url_map.iter_rules()]
    return make_success_data(routers)

@app.route("/log",methods=['GET'])
def get_log():
    # 获取日志
    #从log/server.log中获取日志
    with open(cfg.get("logging", "logPath"),"r",encoding="utf-8") as f:
        log_content = f.read()
    html_content = "<html><body><pre>{}</pre></body></html>".format(log_content)
    return html_content
    




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
    app.run(host='0.0.0.0', port=port,debug=True)
# uvicorn boot:app --reload
# sqlite 的图形化管理工具
# sqlite_web --password -x -p 8088 data/datas_db.db
