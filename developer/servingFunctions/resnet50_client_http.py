import requests
import numpy as np
predict_request='{"signature_name": "serving_default", "instances":[{"inputs":%s}] }' %np.random.rand(224,224,3).tolist()
response = requests.post(url="http://192.168.1.120:31212/function/resnet-50", data=predict_request)


response.raise_for_status()
print(response.elapsed.total_seconds() * 1000)
