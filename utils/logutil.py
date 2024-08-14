# -*- coding: utf-8 -*-
import logging
from logging import handlers
from utils.configer import cfg
import os

class LogUtil(object):
    #  日志级别关系映射
    level_relations = {
        'debug': logging.DEBUG,
        'info': logging.INFO,
        'warning': logging.WARNING,
        'error': logging.ERROR,
        'crit': logging.CRITICAL
    }
    log_conf_section = "logging"
    
    my_instance=None
    
    def __init__(self):
        self._do_init()
        
    def _do_init(self):
        self.__log_file_path = cfg.get(self.log_conf_section, "logPath")
        log_dir  = os.path.dirname(self.__log_file_path)
        if not os.path.exists(log_dir):
            os.makedirs(log_dir)
        self.__log2stream = cfg.getInt(self.log_conf_section, "log2stream")

    #def getLogger(self, fmt='%(asctime)s - %(process)d - %(thread)d - %(funcName)s - %(filename)s[line:%(lineno)d] - %(levelname)s: %(message)s'):
    def getLogger(self, fmt='%(asctime)s - %(thread)d - %(filename)s:%(lineno)d - %(levelname)s - %(funcName)s: %(message)s'):
        logger = logging.getLogger(self.__log_file_path)
        # 设置日志格式
        format_str = logging.Formatter(fmt)
        logger.setLevel(self.level_relations.get(cfg.get(self.log_conf_section, "level")))  # 设置日志级别

        # 往文件里写入 指定间隔时间自动生成文件的处理器
        when = cfg.get(self.log_conf_section, "when")
        back_count = cfg.getInt(self.log_conf_section, "backupCount")
        encoding = cfg.get(self.log_conf_section, "encoding")
        th = handlers.TimedRotatingFileHandler(filename=self.__log_file_path, when=when, backupCount=back_count,
                                               encoding=encoding)
        
    
        th.setFormatter(format_str)  # 设置文件里写入的格式

        logger.addHandler(th)

        if self.__log2stream == 1:  # 根据配置，是否往控制台输出
            sh = logging.StreamHandler()  # 往屏幕上输出
            sh.setFormatter(format_str)  # 设置屏幕上显示的格式
            logger.addHandler(sh)  # 把对象加到logger里
        return logger
    
def get_logger():
    if not LogUtil.my_instance:
        LogUtil.my_instance = LogUtil().getLogger()
    return LogUtil.my_instance

logger = get_logger()
