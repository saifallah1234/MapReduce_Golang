package main

import (
    "encoding/json"
    "log"
    "net/http"
    "MapReduceProject/models"
    "bufio"
    "path/filepath"
    "runtime"
    "os"
)
type Data struct {
    TasksMap    map[string]models.Status `json:"tasks"`
    Workersmap map[int]string           `json:"users"`
}





func dataHandler(w http.ResponseWriter, r *http.Request) {
    mu.Lock()           
    defer mu.Unlock()
    response := Data{
        TasksMap:    tasks,
        Workersmap: workers,
    }
    
    w.Header().Set("Content-Type", "application/json")
    err := json.NewEncoder(w).Encode(response)
    // use the same tasks map from main.go
    if err != nil {
        http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
        return
    }
}
func reduceFileHandler(w http.ResponseWriter, r *http.Request) {
    file, err := os.Open("../Workers/mrtmp.reduce")
    if err != nil {
        if os.IsNotExist(err) {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode([]string{})
            return
        }
        // For other errors, respond with 500 error
        http.Error(w, "Could not open reduce file", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lines := []string{}
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        http.Error(w, "Error reading file", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(lines)
}

/*func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "C:/Users/saif allah/Desktop/MapReduceProject/TailWindCSS_test/src/index.html")
}*/

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    _, filename, _, _ := runtime.Caller(0)
    projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
    htmlPath := filepath.Join(projectRoot, "TailWindCSS_test", "src", "index.html")
    http.ServeFile(w, r, htmlPath)
}

func StartDashboardServer() {
    http.HandleFunc("/data", dataHandler)
    http.HandleFunc("/reducefile", reduceFileHandler) 
    http.HandleFunc("/", dashboardHandler)
    log.Println("Dashboard listening on http://localhost:8080/")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatalf("Dashboard server failed: %v", err)
    }
    }


