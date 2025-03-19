import hashlib
import json

ukey = "7JPWIA1SGV9N17LE"
#ukey="dny75"
def generate_md5_sign(data):
    # 将 JSON 数据转换为字符串并追加 ukey
    sign_str = json.dumps(data, separators=(',', ':'), ensure_ascii=False) + f"key={ukey}"
    print("generate_md5_sign",sign_str)
    # 进行 MD5 加密
    md5_hash = hashlib.md5(sign_str.encode('utf-8')).hexdigest().upper()
    print("generate_md5_sign",md5_hash)
    return md5_hash

def makeSign(payload):
    dataS =json.dumps( payload , separators=(',', ':'),ensure_ascii=False)
    
    data = dataS + "key=" + ukey
    print("makeSign",data)
    newSign=hashlib.md5(data.encode('utf-8')).hexdigest().upper()
    print("makeSign",newSign)
    return newSign


# 示例数据
data = {   
    "query_time": 1672129283,    
    "car_number": "京A5566TT" 
}
generate_md5_sign(data)
makeSign(data)
data = {     
    "car_number": "京A5566TT",
    "query_time": 1672129283     
}
generate_md5_sign(data)
makeSign(data)
# # 生成 sign
# sign = generate_md5_sign(data, ukey)
# print("      "+sign)
# print("need: 9BBD511E3178A91E66AD4C6A043CEC55")



# def test():
#     a = '{"car_number":"浙AB19097","query_time":1672129283}'
#     b = makeSign(a)
#     print(b)

# sign=makeSign(data)
# print("      "+sign)

print("need: 9BBD511E3178A91E66AD4C6A043CEC55")
