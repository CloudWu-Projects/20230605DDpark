import requests





aurl=[
'http://7niu.hyman.store/ffmpeg_hls_go/fileversion.txt',
'http://7niu.hyman.store/ffmpeg_hls_go/ffmpeg_hls_go.7z'
]
def prepareURL(baseUrl):
    #为下载 URL 加上过期时间 e 参数，Unix时间戳：

    #DownloadUrl = 'http://78re52.com1.z0.glb.clouddn.com/resource/flower.jpg?e=1451491200'
    #获取当前时间戳：

    import time
    expire = int(time.time())+3600
    
    #拼接下载 URL：

    DownloadUrl = baseUrl + '?e=' + str(expire)

    
    access_key = 'z8OskufSjOkWbJt7j7asi-1uWp82_ed7l66MkJNT'
    secret_key = 'DcXQn07x_qZU5B2h8CNopJwlv4ceosM3PTlep1gs'
   # 3.对上一步得到的 URL 字符串计算HMAC-SHA1签名（假设访问密钥（AK/SK）是 MY_SECRET_KEY），并对结果做URL 安全的 Base64 编码：

#Sign = hmac_sha1(DownloadUrl, 'MY_SECRET_KEY')
#EncodedSign = urlsafe_base64_encode(Sign)
    import hmac
    import hashlib
    import base64
    def hmac_sha1(key, message):
        return hmac.new(key.encode('utf-8'), message.encode('utf-8'), hashlib.sha1).digest()
    def urlsafe_base64_encode(s):
        return base64.urlsafe_b64encode(s).decode('utf-8').replace('=', '')
    Sign = hmac_sha1(secret_key, DownloadUrl)
    EncodedSign = urlsafe_base64_encode(Sign)
#    4.将访问密钥（AK/SK）（假设是 MY_ACCESS_KEY）与上一步计算得到的结果用英文符号 : 连接起来：

#Token = 'MY_ACCESS_KEY:438dd8pXocjYuF-6dTcKMtETB2g='
#5.将上述 Token 拼接到含过期时间参数 e 的 DownloadUrl 之后，作为最后的下载 URL：

#RealDownloadUrl = 'http://78re52.com1.z0.glb.clouddn.com/resource/flower.jpg?e=1451491200&token=MY_ACCESS_KEY:438dd8pXocjYuF-6dTcKMtETB2g='
    
    DownloadUrl = DownloadUrl + '&token=' + access_key+':' + EncodedSign
    print(DownloadUrl)
    
    return DownloadUrl


def download(url,targetPath=None):
    print(f'正在下载... {url}')
    if targetPath is None:
        targetPath= url.split('/')[-1]
    headers={'Referer':'http://7niu.hyman.store/'}

    headers={'Referer':'http://1qazxsw2.com'}
    r = requests.get(url,headers=headers)



    if r.status_code == 200:
        #创建目录
        import os
        if not os.path.exists(os.path.dirname(targetPath)):
            os.makedirs(os.path.dirname(targetPath))
        with open(targetPath, 'wb') as f:
            f.write(r.content)
        print('下载完成')
    else:
        print('下载失败')
        print(r.status_code)
        print(r.text)

def downloadMutile():
    for url in aurl:
        #download(url)
        download(prepareURL(url))

def download_from_args():
    import argparse
    parser = argparse.ArgumentParser(description="Download files from Lanzou Cloud.")
    parser.add_argument(
        "bucketPath", 
        nargs="?",  # 可选参数
        default="",  # 默认值
        help="download remote file from qiniu"
    )
    args = parser.parse_args()
    print(args.bucketPath)
    #'http://7niu.hyman.store/ffmpeg_hls_go/ffmpeg_hls_go.7z'
    downloadUrl='http://7niu.hyman.store/'+args.bucketPath
    download(prepareURL(downloadUrl),args.bucketPath)

if __name__ == '__main__':
    #downloadMutile()
    download_from_args()