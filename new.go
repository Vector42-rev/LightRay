package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
)

var tmpl = `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Ray Docker Controller</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f9;
            margin: 0;
            padding: 0;
        }

        .container {
            width: 80%;
            margin: auto;
            padding: 20px;
            background-color: white;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            border-radius: 10px;
        }

        h1 {
            text-align: center;
            color: #333;
        }

        .button-container {
            display: flex;
            flex-direction: column;
            gap: 15px;
            margin-top: 20px;
        }

        .button-container button {
            padding: 10px;
            font-size: 16px;
            color: white;
            background-color: #4CAF50;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            transition: background-color 0.3s ease;
        }

        .button-container button:hover {
            background-color: #45a049;
        }

        pre {
            background-color: #f0f0f0;
            padding: 15px;
            border-radius: 5px;
            font-family: monospace;
            white-space: pre-wrap;
            word-wrap: break-word;
        }

        .footer {
            text-align: center;
            margin-top: 30px;
            font-size: 14px;
            color: #777;
        }

        .footer a {
            color: #4CAF50;
            text-decoration: none;
        }
    </style>
</head>
<body>

    <div class="container">
        <h1>Ray Container Management</h1>

        <div class="button-container">
            <button onclick="send('/start')">Start Docker Container</button>
            <button onclick="send('/env')">Set Ray Environment Variables</button>
            <button onclick="send('/copy')">Copy Training File</button>
            <button onclick="send('/ray?cmd=ray start --head --port=6379 --dashboard-host=0.0.0.0')">Start Ray Head</button>
            <button onclick="send('/ray?cmd=ray status')">Check Ray Status</button>
            <button onclick="send('/train?cmd=python ray_train_5_gpu.py')">Start Training</button>
            <button onclick="send('/ray?cmd=ray stop')">Stop Ray Cluster</button>
            <button onclick="send('/stop')">Stop Docker Container</button>
        </div>

        <pre id="output"></pre>
    </div>

    <div class="footer">
        <p>&copy; 2025 Ray Docker Controller | <a href="https://github.com" target="_blank">GitHub Repository</a></p>
    </div>

    <script>
        function send(url) {
            fetch(url)
                .then(response => response.text())
                .then(data => {
                    document.getElementById('output').textContent = data;
                });
        }
    </script>

</body>
</html>
`

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/start", startContainer)
	http.HandleFunc("/env", setRayEnv)
	http.HandleFunc("/ray", runRayCommand)
	http.HandleFunc("/copy", fileCopy)
	http.HandleFunc("/train", runRayCommand)
	http.HandleFunc("/stop", stopContainer)
	fmt.Println("[+] Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("webpage").Parse(tmpl))
	t.Execute(w, nil)
}

func startContainer(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("docker", "run", "--gpus", "all", "--name", "ray-head", "--network=host", "rayproject/ray-ml:latest-py39-cu118", "sleep", "infinity")
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, string(output)+"\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Started Container:\n%s", output)
}

func setRayEnv(w http.ResponseWriter, r *http.Request) {
	envScript := `
LARGE_N=99999999999
export RAY_health_check_timeout_ms=$LARGE_N
export RAY_grpc_keepalive_time_ms=$LARGE_N
export RAY_grpc_client_keepalive_time_ms=$LARGE_N
export RAY_grpc_client_keepalive_timeout_ms=$LARGE_N
export RAY_health_check_initial_delay_ms=$LARGE_N
export RAY_health_check_period_ms=$LARGE_N
export RAY_health_check_failure_threshold=10
`
	cmd := exec.Command("docker", "exec", "ray-head", "bash", "-c", envScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, string(output)+"\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Set Ray Environment Variables:\n%s", output)
}

func runRayCommand(w http.ResponseWriter, r *http.Request) {
	cmdStr := r.URL.Query().Get("cmd")
	cmd := exec.Command("docker", "exec", "ray-head", "bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, string(output)+"\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Command Output:\n%s", output)
}

func stopContainer(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("docker", "rm", "-f", "ray-head")
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, string(output)+"\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Stopped and Removed Container:\n%s", output)
}

func fileCopy(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("docker", "cp", "/home/veera/ray_train_5_gpu.py", "ray-head:/home/ray")
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, string(output)+"\n"+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "File copied:\n%s", output)
}

