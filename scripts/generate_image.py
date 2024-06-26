import json
import time
import sys

import requests
import base64

from sys import argv

script, datapath, text, style, userid, model_id, api, secret = argv

# cd C:\DEV\GO\ai_telegram_bot
# python C:/DEV/GO/ai_telegram_bot/scripts/generate_image.py C:/DEV/GO/ai_telegram_bot/data "Green forest near lake" "ANIME" "111" "4" "98CEC2F87ABC4AAC1F609FDFCBCB7ACD" "23E24C8712FEC1D4F0F571EB2BB1CFF4"

class Text2ImageAPI:

    def __init__(self, url, api_key, secret_key):
        self.URL = url
        self.AUTH_HEADERS = {
            'X-Key': f'Key {api_key}',
            'X-Secret': f'Secret {secret_key}',
        }

    def get_model(self):
        response = requests.get(self.URL + 'key/api/v1/models', headers=self.AUTH_HEADERS)
        data = response.json()
        return data[0]['id']

    def generate(self, prompt, model, images=1, width=1024, height=1024):
        params = {
            "type": "GENERATE",
            "style": style,
            "numImages": images,
            "width": width,
            "height": height,
            "generateParams": {
                "query": f"{prompt}"
            }
        }

        data = {
            'model_id': (None, model),
            'params': (None, json.dumps(params), 'application/json')
        }
        response = requests.post(self.URL + 'key/api/v1/text2image/run', headers=self.AUTH_HEADERS, files=data)
        data = response.json()
        #print(data)
        return data['uuid']

    def check_generation(self, request_id, attempts=13, delay=5):
        #a 13  12  11  10  9   8   7     6     5     4     3     2     1
        #t 15  20  26  33  41  50  1.00  1.11  1.23  1.36  1.50  2.05  2.21
        #d 5   6   7   8   9   10  11    12    13    14    15    16    0
        # меньше 10 секунд не бывает генераций, поэтому не пингуем в начале
        # t - в какой момент времени выполняется запрос, 0 это запуск скрипта
        # d - пауза после запроса
        # 4 sec. - столько необходимо для прочих функций (кроме check_generation)
        time.sleep(15)
        while attempts > 0:
            response = requests.get(self.URL + 'key/api/v1/text2image/status/' + request_id, headers=self.AUTH_HEADERS)
            data = response.json()
            if data['status'] == 'DONE':
                #print(data)
                return data['images']
            #else:
                #print(data)                    

            attempts -= 1
            
            # на последней итерации нет смысла вставать на паузу
            if attempts != 0 :
                time.sleep(delay)
                delay += 1


if __name__ == '__main__':
    
    api = Text2ImageAPI('https://api-key.fusionbrain.ai/', api, secret)
    #model_id = api.get_model()
    uuid = api.generate(text, model_id)
    images = api.check_generation(uuid)   
    
    if images != None:

        base64_img  = images[0]
        base64_img_bytes = base64_img.encode('utf-8')
            
        file_path = datapath+'/img_'+userid+'_0.png'
            
        with open(file_path, 'wb') as file_to_save:
            decoded_image_data = base64.decodebytes(base64_img_bytes)
            file_to_save.write(decoded_image_data)
            
        print(file_path)

