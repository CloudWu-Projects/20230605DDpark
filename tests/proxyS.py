responseJson={
	"code": 0,
	"data": [{
			"mchTradeId": "DFG678932NADASD",
			"plateId": "冀ATTTT1",
			"tradeType": "0",
			"successTime": "2025-07-24 00:00:00",
			"tradeAmount": 100
		},
		{
			"mchTradeId": "DHFHFDUSFJKS134324",
			"plateId": "冀ATTTT2",
			"tradeType": "1",
			"successTime": "2025-07-24 00:00:00",
			"tradeAmount": 100
		}
	],
	"msg": ""
}


import http.server
import json

class MyHandler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        print(self.path)
        if self.path == '/proxy':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(bytes(json.dumps(responseJson), 'utf-8'))
        else:
            self.send_error(404)
    
    def do_POST(self):
        print(self.path)
        if self.path == '/proxy':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(bytes(json.dumps(responseJson), 'utf-8'))
        else:
            self.send_error(404)

with http.server.HTTPServer(('localhost', 28080), MyHandler) as httpd:
    print('Serving at port', 28080)
    httpd.serve_forever()
#http.server.test(HandlerClass=MyHandler, port=28080)
#http.server.test(HandlerClass=http.server.SimpleHTTPRequestHandler, port=8080)