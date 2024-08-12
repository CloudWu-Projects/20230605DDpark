# -*- coding: utf-8 -*-
import sqlite3
from utils.configer import cfg
from utils.logutil import logger

conf_section = "sqlite"


class SqliteUtil:
    # 连接对象
    __conn = None

    def __init__(self):
        self.__do_init()
        
    def __do_init(self):
        self.__data_path = cfg.get(conf_section, "data_path")
        
    def __get_conn(self):
        try:
            return sqlite3.connect(self.__data_path)
        except Exception as e:
            logger.error("创建数据库连接失败:%s",e)
            raise Exception("创建数据库连接失败!")

    def getOneBySql(self, sql, param=None):
        conn = self.__get_conn()
        cursor =  conn.cursor()
        # 返回数据格式改成字典
        cursor.row_factory = sqlite3.Row

        if param is None:
            cursor.execute(sql)
        else:
            cursor.execute(sql, param)
        data = cursor.fetchone()
        if cursor is not None :
            cursor.close()
        if conn is not None :
            conn.close()
        return data

    def __query(self, sql, param=None):
        count = 0
        conn = self.__get_conn()
        cursor =  conn.cursor()
        try:
            if param is None:
                count = cursor.execute(sql)
            else:
                count = cursor.execute(sql, param)
            conn.commit()
        except Exception as e:
            logger.error("查询出现异常:%s",e)
            conn.rollback()
        finally :
            if cursor is not None :
                cursor.close()
            if conn is not None :
                conn.close()
        return count

    def updateBySql(self, sql, param=None):
        """
        @summary: 更新数据表记录
        @param sql: ＳＱＬ格式及条件，使用(?,?)
        @param param: 要更新的  值 tuple
        @return: count 受影响的行数
        """
        return self.__query(sql, param)
