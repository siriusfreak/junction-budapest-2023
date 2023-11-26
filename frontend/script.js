document.getElementById('videoInput').addEventListener('change', function() {
    let fileName = this.files[0].name;
    let maxLength = 30; // Maximum length of file name to display

    if (fileName.length > maxLength) {
        fileName = fileName.substring(0, maxLength - 3) + '...'; // Truncate and add ellipsis
    }

    document.getElementById('fileName').innerText = fileName;
});

document.getElementById('customButton').addEventListener('click', function() {
    resetCheckResults();
    document.getElementById('videoInput').click();
});

document.getElementById('uploadButton').addEventListener('click', function() {
    var videoFile = document.getElementById('videoInput').files[0];
    if (videoFile) {
        setUploadControlsEnabled(false); // Disable buttons
        uploadVideo(videoFile); // Call the updated upload function
    } else {
        alert("Please select a video file first.");
    }
});

function uploadVideo(file) {
    resetCheckResults()

    const formData = new FormData();
    formData.append('video', file);

    var progressBarContainer = document.getElementById('progressBarContainer');
    var progressBar = document.getElementById('progressBar');

    var detailsContainer = document.getElementById('videoFakeCandidatDetails');
    detailsContainer.style.display = 'none';

    progressBarContainer.style.display = 'block';
    progressBar.style.width = '0%';
    progressBar.textContent = '0%';

    fetch('/proxy-upload', {
        method: 'POST',
        body: formData
    })
        .then(response => response.json())
        .then(data => {
            console.log("Upload successful. UID:", data.uid);
            checkUploadStatus(data.uid);
        })
        .catch(error => {
            console.error("Upload failed:", error);
            displayFail();
        });
}
function checkUploadStatus(uid) {
    const checkInterval = setInterval(function() {
        fetch(`/proxy-status/${uid}`)
            .then(response => response.json())
            .then(data => {
                updateProgressBar(data.completion_percentage);
                if (data.completion_percentage === 100) {
                    clearInterval(checkInterval);
                    displayDetails(data.VideoFakeCandidat);
                    if (data.confidence_level) {
                        displaySuccess();
                    } else {
                        displayFail();
                    }
                }
            })
            .catch(error => {
                console.error("Status check failed:", error);
                clearInterval(checkInterval);
                displayFail(); // Handle failure in case of an error
            });
    }, 1000); // Check every 1 second
}

function updateProgressBar(percentage) {
    const roundedPercentage = Math.round(percentage); // Rounds to the nearest whole number
    const progressBar = document.getElementById('progressBar');
    progressBar.style.width = roundedPercentage + '%';
    progressBar.textContent = roundedPercentage + '%';
}
function setUploadControlsEnabled(enabled) {
    document.getElementById('videoInput').disabled = !enabled;
    document.getElementById('customButton').disabled = !enabled;
    document.getElementById('uploadButton').disabled = !enabled;
}

function displaySuccess() {
    // Hide the progress bar
    var successImageContainer = document.getElementById('successImageContainer');
    successImageContainer.style.display = 'block';
    setUploadControlsEnabled(true);
}

function displayFail() {
    // Hide the progress bar
    var failImageContainer = document.getElementById('failImageContainer');
    failImageContainer.style.display = 'block';
    setUploadControlsEnabled(true);
}

function resetCheckResults() {
    document.getElementById('successImageContainer').style.display = 'none';
    document.getElementById('failImageContainer').style.display = 'none';
}

function displayDetails(videoFakeCandidat) {
    const detailsContainer = document.getElementById('videoFakeCandidatDetails');
    detailsContainer.innerHTML = `
        <p>Audio Fake Detection: ${videoFakeCandidat.AudioFakeDetectionResult}</p>
        <p>Deepfake Detection: ${videoFakeCandidat.DeepfakeDetectResult}</p>
        <p>Whisper Detection: ${videoFakeCandidat.WhisperLargeV3Result}</p>
        <p>One Person Detection: ${videoFakeCandidat.OnePersonDetectResult}</p>
    `;
    detailsContainer.style.display = 'block';
}