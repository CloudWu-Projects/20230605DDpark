# -*- coding: utf-8 -*-
# flake8: noqa

from qiniu import Auth, put_file, etag,CdnManager
import qiniu.config

#需要填写你的 Access Key 和 Secret Key
access_key = 'z8OskufSjOkWbJt7j7asi-1uWp82_ed7l66MkJNT'
secret_key = 'DcXQn07x_qZU5B2h8CNopJwlv4ceosM3PTlep1gs'


def upload(localfile,key=None):
    #构建鉴权对象
    q = Auth(access_key, secret_key)

    #要上传的空间
    bucket_name = 'cloud-wu'

    #上传后保存的文件名
    if key is None:
        key = localfile.split('/')[-1]
    
   

    #生成上传 Token，可以指定过期时间等
    token = q.upload_token(bucket_name, key, 3600)

    #要上传文件的本地路径
    

    ret, info = put_file(token, key, localfile, version='v2')
    print(ret)
    print("----")
    print(info.url)
    assert ret['key'] == key
    assert ret['hash'] == etag(localfile)
    urls = ['http://7niu.hyman.store/'+key]
    cdn_manager = CdnManager(q)
    refresh_url_result = cdn_manager.refresh_urls(urls)
    print(refresh_url_result)
    # 刷新目录

import os ,re 

def extract_version_from_filename(file_path):
    try:
        # 检查文件路径是否存在
        

        # 获取文件名
        file_name = os.path.basename(file_path)
        print(file_name)
       # 使用正则表达式提取项目名称和版本号
        match = re.search(r'([a-zA-Z0-9_]+)_v(\d+\.\d+\.\d+)', file_name)
        if match:
            project_name = match.group(1)
            version = match.group(2)
            return project_name, version
        else:
            raise ValueError(f"文件名 {file_name} 中未找到版本号")
    except Exception as e:
        print(f"发生错误: {e}")
        return None
#

# read version from ../version.txt
def read_version_from_file(file_path):
    try:
        # 检查文件路径是否存在
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"文件 {file_path} 不存在")
        # 读取文件内容
        with open(file_path, "r") as file:
            version = file.read().strip()
            return version
    except Exception as e:
        print(f"发生错误: {e}")
        return None
    
import sys
# 获取命令行参数
print("脚本名:", sys.argv[0])  # 第一个参数是脚本名
print("传递的参数:", sys.argv[1:])  # 其余的参数是传递的参数

# 打印每个参数
for index, arg in enumerate(sys.argv):
    print(f"参数 {index}: {arg}")

localFolder="../"
if len(sys.argv) > 1:
    localFolder = sys.argv[1]

def upload_from_version():
    
        
    version = read_version_from_file(localFolder+"/version.txt")
    print(version)

    localfile = f'{localFolder}/ffmpeg_hls_go_v{version}.7z'
    filename,version = extract_version_from_filename(localfile)
    print(localfile)
    print(filename,version)
    upload(localfile,key=f"{filename}/{filename}.7z")

    with open(f"{localFolder}/fileversion.txt", "w") as file:
        file.write(f"{version}")

    upload(f"{localFolder}/fileversion.txt",key=f"{filename}/fileversion.txt")

import argparse
def upload_from_args():
    parser = argparse.ArgumentParser(description="Download files from Lanzou Cloud.")
    parser.add_argument(
        "localfile", 
        nargs="?",  # 可选参数
        default="",  # 默认值
        help="upload localfile to qiniu"
    )
    args = parser.parse_args()
    filename="yilingshequ_go"
    upload(args.localfile,key=f"{filename}/{filename}.zip")

if __name__ == '__main__':
    #upload_from_version()
    upload_from_args()