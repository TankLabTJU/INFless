#import base64
import requests
import numpy as np
def handle(req):
    """handle a request to the function
    Args:
        req (str): request body
    """
    #SERVER_URL = 'http://127.0.0.1:8501/v1/models/resnet:predict'
    #dl_request = requests.get("https://note.youdao.com/yws/public/resource/ac12306f6711e44a7b6ba15a26a156b1/xmlnote/9A8F6AB2B96549F0B648CFE6F12C45F9/17701", stream=True)
    #dl_request.raise_for_status()
    #jpeg_bytes = base64.b64encode(IMAGE_PATH)
    # Compose a JSON Predict request (send JPEG image in base64)
    #for i in range(1,100):
     #   jpeg_bytes = base64.b64encode(dl_request.content).decode('utf-8')
      #  predict_request = '{"instances" : [{"b64": "%s"}]}' % jpeg_bytes
       # response = requests.post(SERVER_URL, data=predict_request)
       # response.raise_for_status()

    # prediction = response.json()['predictions'][0]
       # print(response.json())
        #print(response.elapsed.total_seconds() * 1000)
    predict_request='{"signature_name": "serving_default", "instances":[{"inputs":%s}] }' %np.random.rand(224,224,3).tolist()
    response = requests.post(url="http://127.0.0.1:8501/v1/models/resnet:predict", data=predict_request)
    response.raise_for_status()
    print(response.elapsed.total_seconds() * 1000)

    return req
