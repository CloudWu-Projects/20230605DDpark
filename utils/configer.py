# -*- coding: utf-8 -*-

import json


# 配置文件固定在根目录的conf文件夹下
# 并且该文件夹地址，只有根目录的主程序运行时才能生效！！！
log_file_path = "conf/conf.json"
encoding = "utf-8"

class Configer:
    my_instance=None
    
    def __init__(self):
        self.load()        
        
    def load(self):
        try:
            with open(log_file_path,'r',encoding='utf-8') as config_file:
                self.__config = json.load(config_file)
        except Exception as e:
            print(e)

    def flush(self):
        with open(log_file_path,'w',encoding='utf-8') as config_file:
            json.dump(self.__config, config_file,indent=4,ensure_ascii=False)

    
    def get(self, section, option):
        if section in self.__config:
            if option in self.__config[section]:
                return self.__config[section][option]  
            raise Exception(f"配置文件中无{section} {option}配置")          
        raise Exception("配置文件中无【" + section + "】配置")
    
    def set(self,section,option,value):
        if section in self.__config:
            self.__config[section][option] = value
            self.flush()
        else:
            self.__config[section] = {option:value}
            self.flush()
        cfg.load()

    def getParkinfo(self,parkid):
        return cfg.get("parkinfo",parkid)

    def getConfig(self):
        return self.__config

    def getInt(self, section, option):
        return int(self.get(section, option))
    
    
def get_cfg():
    if not Configer.my_instance:
        Configer.my_instance = Configer()
    return Configer.my_instance

cfg = get_cfg()