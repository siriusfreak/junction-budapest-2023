import base64
import io
import numpy as np
import cv2
import numpy as np
from PIL import Image
import pathlib
import tensorflow
from fastapi import FastAPI, UploadFile, File, HTTPException
from tempfile import TemporaryDirectory
import os
import shutil

app = FastAPI()

PROB_THRESHOLD = 0.4 # Minimum probably to show results.
class Model:
    def __init__(self, model_filepath):
        self.graph_def = tensorflow.compat.v1.GraphDef()
        self.graph_def.ParseFromString(model_filepath.read_bytes())

        input_names, self.output_names = self._get_graph_inout(self.graph_def)
        assert len(input_names) == 1 and len(self.output_names) == 3
        self.input_name = input_names[0]
        self.input_shape = self._get_input_shape(self.graph_def, self.input_name)

    def predict(self, image_filepath):
        image = Image.fromarray(image_filepath).resize(self.input_shape)
        input_array = np.array(image, dtype=np.float32)[np.newaxis, :, :, :]

        with tensorflow.compat.v1.Session() as sess:
            tensorflow.import_graph_def(self.graph_def, name='')
            out_tensors = [sess.graph.get_tensor_by_name(o + ':0') for o in self.output_names]
            outputs = sess.run(out_tensors, {self.input_name + ':0': input_array})
            return {name: outputs[i][np.newaxis, ...] for i, name in enumerate(self.output_names)}

    @staticmethod
    def _get_graph_inout(graph_def):
        input_names = []
        inputs_set = set()
        outputs_set = set()

        for node in graph_def.node:
            if node.op == 'Placeholder':
                input_names.append(node.name)

            for i in node.input:
                inputs_set.add(i.split(':')[0])
            outputs_set.add(node.name)

        output_names = list(outputs_set - inputs_set)
        return input_names, output_names

    @staticmethod
    def _get_input_shape(graph_def, input_name):
        for node in graph_def.node:
            if node.name == input_name:
                return [dim.size for dim in node.attr['shape'].shape.dim][1:3]

def print_outputs(outputs, gambar):
  image = gambar
  assert set(outputs.keys()) == set(['detected_boxes', 'detected_classes', 'detected_scores'])
  l, t, d = image.shape
  labelopen = open("labels.txt", 'r')
  labels = [line.split(',') for line in labelopen.readlines()]

  eyes = []
  for box, class_id, score in zip(outputs['detected_boxes'][0], outputs['detected_classes'][0], outputs['detected_scores'][0]):
    if score > PROB_THRESHOLD:
      if class_id == 0:
        eyes.append("open")
      else:
        eyes.append("closed")
          

  return eyes

p = pathlib.Path("model.pb")
model = Model(p)

@app.post("/predict")
async def predict(video: UploadFile = File(...)):
    with TemporaryDirectory() as temp_dir:
        tmp_path = os.path.join(temp_dir, 'temp_video.mp4')
        with open(tmp_path, 'wb') as tmp_file:
            shutil.copyfileobj(video.file, tmp_file)
        cap = cv2.VideoCapture(tmp_path)

    try:
        if not cap.isOpened():
            raise HTTPException(status_code=400, detail="Unable to read video file.")
            
        fps = cap.get(cv2.CAP_PROP_FPS)
        if fps == 0:
            raise HTTPException(status_code=400, detail="FPS of video is zero, which indicates a problem with the video file.")

        total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
        video_length = total_frames / fps

        if video_length > 20:
            raise HTTPException(status_code=400, detail=f"Video is too long. Maximum length allowed is {150} seconds.")

        if fps > 60:
            raise HTTPException(status_code=400, detail=f"Video FPS is too high. Maximum FPS allowed is {60}.")

        frame_count = 0
        frames = 0
        processed_count = 0
        fake = 0

        while True:
            percent = (frames/total_frames * 100)
            print(f"{percent:.2f}")

            ret, frame = cap.read()
            if not ret:
                break

            frames += 1
            if frames % 2 == 0:
                continue
            
            frame = cv2.resize(frame, (500,500), interpolation = cv2.INTER_AREA)

            outputs = model.predict(frame)
            eyes = print_outputs(outputs, frame)

            if len(eyes) != 2 or eyes[0] != 1 or eyes[1] != 1:
                fake += 1

            processed_count += 1
    finally:
        cap.release()

    return {"processed_count":processed_count, "fake_eyes": fake}
