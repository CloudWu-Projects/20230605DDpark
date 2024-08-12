# -*- coding: utf-8 -*-
#各省车牌简称与数据库表名称的对照
carplate_tabel={
"川":"oss_sc",
"贵":"oss_gz",
"云":"oss_yn",
"陕":"oss_sx_xian",
"甘":"oss_gs",
"晋":"oss_sx_taiyuan",
"冀":"oss_hb",
"豫":"oss_henan",
"津":"oss_tj",
"京":"oss_bj",
"渝":"oss_cq",
"沪":"oss_sh",
"藏":"oss_xz",
"新":"oss_xj",
"蒙":"oss_nm",
"宁":"oss_nx",
"粤":"oss_gd",
"苏":"oss_js",
"鲁":"oss_sd",
"湘":"oss_hunan",
"鄂":"oss_hubei",
"皖":"oss_anhui",
"黑":"oss_heilongjiang",
"辽":"oss_liaoning",
"吉":"oss_jilin",
"浙":"oss_zhejiang",
"桂":"oss_guangxi",
"青":"oss_qinghai",
"琼":"oss_hainan",
"闽":"oss_fujian",
"赣":"oss_jiangxi"
}
# 排放标准中文对照
paifang ={
"国四":"国4",
"国五":"国5",
"国六":"国6"
}
# 查询停车场数据
SELECT_PARK_SQL = "select id FROM park where park_id=?"
# 查询数据
SELECT_CAR_INFO_SQL = "select oss,pfjd FROM TABLE_NAME where vin=?"
# 新增数据
INSERT_CAR_INFO_SQL = "INSERT INTO TABLE_NAME ('carnumber','vin','motor','oss','pfjd','url') VALUES (?,?,?,?,?,?)"

ERROR_CODE = 0
SUCCESS_CODE = 1

def make_success_data(data):
    return {"state":SUCCESS_CODE,"data":data,"errmsg":"成功"}

def make_error_data(errmsg):
    return {"state":ERROR_CODE,"errmsg":errmsg}


SYSTEM = "LINUX"
