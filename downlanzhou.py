import requests
import sys

import argparse



def downLoadLanZhou(lanzhouurl):
    analyzeUrl = 'https://api.hanximeng.com/lanzou/?url='

    allUrl = analyzeUrl + lanzhouurl

    a = requests.get(allUrl)

    if a.status_code != 200:
        print(f"Failed to analyze URL: {a.status_code}")
        print(a.text)
        exit(1)

    allJson = a.json()

    trueUrl = allJson.get('downUrl')
    if not trueUrl:
        print("Failed to get download URL from JSON response")
        print(allJson)
        exit(1)

    print(trueUrl)
    file_name= allJson.get('name')
    print(file_name)
    try:
        b = requests.get(trueUrl, stream=True)
        b.raise_for_status()  # Raises an HTTPError for bad responses (4xx and 5xx)

        # Save the file to disk
        with open(file_name, 'wb') as f:
            for chunk in b.iter_content(chunk_size=8192):
                f.write(chunk)

        print(f"File downloaded successfully and saved as {file_name}")

    except requests.exceptions.HTTPError as http_err:
        print(f"HTTP error occurred: {http_err}")
    except Exception as err:
        print(f"Other error occurred: {err}")


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description="Download files from Lanzou Cloud.")
    parser.add_argument(
        "lanzhouurl", 
        nargs="?",  # 可选参数
        default="https://wwkd.lanzn.com/iRrqN2r641xe",  # 默认值
        help="The Lanzou Cloud URL to download from."
    )
    args = parser.parse_args()

    # 调用下载函数
    downLoadLanZhou(args.lanzhouurl)