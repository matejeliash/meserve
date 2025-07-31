package main

import "html/template"

func GetTemplate() *template.Template {

	return template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Index of {{ .Path }}</title>
	<style>
		body {
			font-family: sans-serif;
			margin: 2rem;
			text-align:teft;
		}
		a {
			text-decoration: none;
			color: #0366d6;
		}
		a:hover {
			text-decoration: underline;
		}
	</style>
</head>
<body>
<h2>Index of {{ .Path }}</h2>
{{if ne .Path "/"}}<p><a href="../">⬅️ Up</a></p>{{end}}

{{ if .DiskStatus }} <p><strong>Free space:</strong> {{ .DiskStatus }}</p>{{ end }}
<form id="uploadForm">
<p>Upload files to current directory:</p>
      <input type="file" id="fileInput" name="file" multiple required>
      <button type="submit">Upload</button>
  </form>

  <div id="progress">
      <p id="percent"></p>
      <p id="speed"></p>
      <p id="status"> nothing...</p>
  </div>



<pre>
{{range .Files}}
{{if .IsDir}}📁 Dir {{else}}📄 File{{end}} {{.HumanSize}} {{.FormattedModTime}}   <a href="{{.Path}}">{{.Name}}{{if .IsDir}}/{{end}}</a>{{end}}
</pre>
{{if ne .Path "/"}}<p><a href="../">⬅️ Up</a></p>{{end}}

<script>
    const form = document.getElementById('uploadForm');
    const fileInput = document.getElementById('fileInput');
    const percentDisplay = document.getElementById('percent');
    const speedDisplay = document.getElementById('speed');
    const statusDisplay = document.getElementById('status');

    form.addEventListener('submit', function (e) {
        e.preventDefault();

        statusDisplay.textContent="uploading ...";

        const files = fileInput.files;
        if (!files.length) return;

        const formData = new FormData();
        for (const file of files) {
            formData.append('file', file);
        }

        const xhr = new XMLHttpRequest();
        xhr.open('POST', window.location.pathname);

        const startTime = Date.now();

        xhr.upload.onprogress = function (e) {
            if (e.lengthComputable) {
                const percent = (e.loaded / e.total * 100).toFixed(2);
                percentDisplay.textContent = "Progress: " + percent + "%";

                const elapsed = (Date.now() - startTime) / 1000;
                const speed = (e.loaded / 1024 / 1024 / elapsed).toFixed(2);
                speedDisplay.textContent = "Speed:   " + speed + " MB/s";

            }
        };

        xhr.onload = function () {
            if (xhr.status === 200) {
                alert(xhr.responseText);
                statusDisplay.textContent="finished";
            } else {
                alert('Upload failed.');
            }
        };

        xhr.onerror = function () {
            alert('An error occurred while uploading.');
            console.error('Upload error:', xhr.status, xhr.statusText);

        };

        xhr.send(formData);
    });
</script>



</body>
</html>
				`))

}
