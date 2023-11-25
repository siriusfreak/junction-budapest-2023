import requests
import os
import base64

# URL вашего FastAPI приложения
service_url = "http://localhost:8000/upload"

# Путь к папке с MP4 файлами
directory_path = "./videos/"
save_directory_path = "./frames/"
# Перебираем все MP4 файлы в папке
for filename in os.listdir(directory_path):
    if filename.endswith(".mp4"):
        video_path = os.path.join(directory_path, filename)
        # Отправляем файл на сервис
        with open(video_path, 'rb') as video_file:
            response = requests.post(
                service_url,
                files={"video": video_file},
                params={
                    "processed_percent": 50,
                    "confidence_threshold": 0.3,
                    "skip_milliseconds": 1000,
                }
            )

        # Проверяем ответ от сервиса
        if response.status_code == 200:
            frames = response.json()['frames']
            # Сохраняем каждый кадр в виде изображения
            for index, frame_base64 in enumerate(frames):
                frame = base64.b64decode(frame_base64)
                frame_path = os.path.join(save_directory_path, f"{filename}_frame_{index}.jpg")
                with open(frame_path, 'wb') as frame_file:
                    frame_file.write(frame)
            print(f"Кадры для {filename} сохранены. total_frames:", response.json()['total_frames'], ",frames:", len(frames), ",processed_frames:", response.json()['processed_frames'])
        else:
            print(f"Ошибка при запросе к сервису для файла {filename}: {response.text}")








































