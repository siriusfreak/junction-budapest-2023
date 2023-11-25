import json

import torch
import numpy as np
from models import image

import torchaudio
from fastapi import FastAPI, File, UploadFile
from fastapi.responses import PlainTextResponse
import shutil
import tempfile

# Define the audio_args dictionary
audio_args = {
    'nb_samp': 64600,
    'first_conv': 1024,
    'in_channels': 1,
    'filts': [20, [20, 20], [20, 128], [128, 128]],
    'blocks': [2, 4],
    'nb_fc_node': 1024,
    'gru_node': 1024,
    'nb_gru_layer': 3,
    'nb_classes': 2
}

def preprocess_audio(audio_file):
    waveform, _ = torchaudio.load(audio_file)

    # Add a batch dimension
    audio_pt = torch.unsqueeze(torch.Tensor(waveform[0]), dim=0)
    return audio_pt

def deepfakes_spec_predict(input_audio):
    x = input_audio
    audio = preprocess_audio(x)
    if torch.cuda.is_available():
        audio = audio.cuda()
    spec_grads = spec_model.forward(audio)
    spec_grads_inv = np.exp(spec_grads.cpu().detach().numpy().squeeze())

    # multimodal_grads = multimodal.spec_depth[0].forward(spec_grads)

    # out = nn.Softmax()(multimodal_grads)
    # max = torch.argmax(out, dim = -1) #Index of the max value in the tensor.
    # max_value = out[max] #Actual value of the tensor.
    max_value = np.argmax(spec_grads_inv)

    if max_value > 0.5:
        preds = max_value
        text2 = False

    else:
        preds = 1 - max_value
        text2 = True

    return {
        "confidence": float(preds),
        "fake": text2
    }

def load_spec_modality_model(args):
    spec_encoder = image.RawNet(args)
    ckpt = torch.load('checkpoints/model.pth', map_location = torch.device(args.device))
    spec_encoder.load_state_dict(ckpt['spec_encoder'], strict = True)
    if torch.cuda.is_available():
        spec_encoder = spec_encoder.cuda()
    spec_encoder.eval()
    return spec_encoder


args = {
        'device': "cuda:0" if torch.cuda.is_available() else "cpu",
        'pretrained_image_encoder': False,
        'freeze_image_encoder': False,
        'pretrained_audio_encoder': False,
        'freeze_audio_encoder': False,
        'in_channels': 1,
        'nb_fc_node': 1024,
        'gru_node': 1024,
        'nb_gru_layer': 3,
        'nb_classes': 2
}

class Dict2Class(object):
    def __init__(self, my_dict):
        for key in my_dict:
            setattr(self, key, my_dict[key])

print(args)
cargs = Dict2Class(args)

spec_model = load_spec_modality_model(cargs)


app = FastAPI()
@app.post("/predict")
async def predict(file: UploadFile = File(...)):
    filename = file.filename
    with tempfile.NamedTemporaryFile(suffix=filename) as temp_file:
        # Copy the contents of the uploaded file to the temporary file
        shutil.copyfileobj(file.file, temp_file)
        # Get the path of the temporary file
        temp_file_path = temp_file.name

        with tempfile.NamedTemporaryFile(suffix=filename) as out_file:
            res = deepfakes_spec_predict(temp_file_path)

            return PlainTextResponse(content=json.dumps(res))

