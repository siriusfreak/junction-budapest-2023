import numpy as np
import io
from tensorflow.keras.models import load_model
import imutils
import matplotlib.pyplot as plt
import cv2
import numpy as np
from tensorflow.keras.preprocessing.image import img_to_array
from PIL import Image
import pathlib
import tensorflow
from fastapi import FastAPI, UploadFile, File, HTTPException

app = FastAPI()

class Model:
    def __init__(self, model_filepath):
        self.graph_def = tensorflow.compat.v1.GraphDef()
        self.graph_def.ParseFromString(model_filepath.read_bytes())

        input_names, self.output_names = self._get_graph_inout(self.graph_def)
        assert len(input_names) == 1 and len(self.output_names) == 3
        self.input_name = input_names[0]
        self.input_shape = self._get_input_shape(self.graph_def, self.input_name)

    def predict(self, image):
        image = image.resize(self.input_shape)
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
    PROB_THRESHOLD = 0.4 # Minimum probably to show results.

    image = gambar
    assert set(outputs.keys()) == set(['detected_boxes', 'detected_classes', 'detected_scores'])
    l, t, d = image.shape
    labelopen = open("labels.txt", 'r')

    labels = [line.split(',') for line in labelopen.readlines()]

    for box, class_id, score in zip(outputs['detected_boxes'][0], outputs['detected_classes'][0], outputs['detected_scores'][0]):
        if score > PROB_THRESHOLD:
            print(f"Label: {class_id}, Probability: {score:.5f}, box: ({box[0]:.5f}, {box[1]:.5f}) ({box[2]:.5f}, {box[3]:.5f})")
            x = box[0] * t
            y = box[1] * l
            h = box[2] * t
            w =  box[3] * l
            result_image = cv2.rectangle(image, (int(x), int(y)), (int(h), int(w)),  (255,215,0), 3)
            cv2.putText(result_image, labels[int(class_id)][0], (int(x), int(y)-10), fontFace = cv2.FONT_HERSHEY_SIMPLEX, fontScale = 0.5, color = (255,215,0), thickness = 2)
    
        return result_image

p = pathlib.Path("model.pb")
model = Model(p)

@app.post("/predict/")
async def predict(file: UploadFile = File(...)):
    if not file.content_type.startswith('image/'):
        raise HTTPException(status_code=400, detail="The uploaded file is not an image.")

    contents = await file.read()

    try:
        image_stream = io.BytesIO(contents)
        image = Image.open(image_stream)
    except IOError:
        raise HTTPException(status_code=400, detail="Invalid image file.")

    outputs = model.predict(image)
    
    return print_outputs(outputs, image)

