import json
import time
import sys

import requests
import base64
import random

from sys import argv

script, datapath, text = argv

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
        return data['uuid']

    def check_generation(self, request_id, attempts=10, delay=10):
        while attempts > 0:
            response = requests.get(self.URL + 'key/api/v1/text2image/status/' + request_id, headers=self.AUTH_HEADERS)
            data = response.json()
            if data['status'] == 'DONE':
                return data['images']

            attempts -= 1
            time.sleep(delay)


if __name__ == '__main__':
    
    api = Text2ImageAPI('https://api-key.fusionbrain.ai/', '1B189E2CFA69FFD2130FC56294B96DA9', '7C0B7C7FE4FFA4F6EE9DF0CFA257167C')
    uuid = api.generate(text, 4)
    images = api.check_generation(uuid)   
    
    base64_img  = images[0]
    base64_img_bytes = base64_img.encode('utf-8')
    
    num = str(int(random.uniform(1, 500)))
#    file_path = 'C:/DEV/GO/telegram_bot_1/data/image_'+num+'.jpg'
    file_path = datapath+'/image_'+num+'.jpg'
     
    with open(file_path, 'wb') as file_to_save:
        decoded_image_data = base64.decodebytes(base64_img_bytes)
        file_to_save.write(decoded_image_data)
        
    print(file_path)    