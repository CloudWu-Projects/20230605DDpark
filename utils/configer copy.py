# -*- coding: utf-8 -*-

from configparser import RawConfigParser

# 配置文件固定在根目录的conf文件夹下
# 并且该文件夹地址，只有根目录的主程序运行时才能生效！！！
log_file_path = "conf/conf.conf"
encoding = "utf-8"

class Configer:
    my_instance=None
    
    def __init__(self):
        self.__config = MyConfigParser()
        self.__config.read(log_file_path, encoding=encoding)
        if (len(self.__config.sections()) == 0):  # 配置文件为空时，抛出异常
            raise Exception("配置文件异常,请检查文件位置及内部配置!")

    def get(self, section, option):
        self.__verdictOption(section, option)
        return self.__config.get(section, option)

    def getInt(self, section, option):
        return int(self.get(section, option))
    
    def getWithDefault(self, section, option,default):
        try:
            v = self.get(section, option)
            return v if v else default
        except Exception as e:
            return default
    
    def getOptions(self, section):
        self.__verdictSection(section)
        return self.__config.options(section)

    def hasOption(self, section,option):
        return self.__config.has_option(section, option)

    def __verdictSection(self, section):
        if (self.__config.has_section(section) == False):
            raise Exception("配置文件中无【" + section + "】配置")

    def __verdictOption(self, section, option):
        self.__verdictSection(section)
        if (self.__config.has_option(section, option) == False):
            raise Exception("配置文件【" + section + "】节点中无【" + option + "】配置")


# 解决option转换为小写的问题，重写optionxform方法
class MyConfigParser(RawConfigParser):
    def optionxform(self, optionstr):
        return optionstr
    
def get_cfg():
    if not Configer.my_instance:
        Configer.my_instance = Configer()
    return Configer.my_instance

cfg = get_cfg()