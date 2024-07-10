package main

import (
    "html/template"
    "log"
    "net/http"
    "sync"
	"strconv"
)

type Task struct {
    ID   int
    Name string
}

var (
    tasks   []Task
    taskID  int
    taskMu  sync.Mutex
    tpl     *template.Template
)

func init() {
    tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/tasks", tasksHandler)
    http.HandleFunc("/add", addTaskHandler)
    http.HandleFunc("/delete", deleteTaskHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    tpl.ExecuteTemplate(w, "index.html", nil)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
    tpl.ExecuteTemplate(w, "tasks.html", tasks)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        taskMu.Lock()
        defer taskMu.Unlock()
        taskID++
        task := Task{
            ID:   taskID,
            Name: r.FormValue("name"),
        }
        tasks = append(tasks, task)
        tpl.ExecuteTemplate(w, "task.html", task)
    }
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        taskMu.Lock()
        defer taskMu.Unlock()
        id := r.FormValue("id")
        for i, task := range tasks {
            if strconv.Itoa(task.ID) == id {
                tasks = append(tasks[:i], tasks[i+1:]...)
                break
            }
        }
        tpl.ExecuteTemplate(w, "tasks.html", tasks)
    }
}
