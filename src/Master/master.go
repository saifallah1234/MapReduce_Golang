package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"
	"MapReduceProject/models"
)



var nbR int = 4
var(
 allMapTasksDone bool
 content []string
 mapTasks chan models.Task
 reduceTasks chan models.Task
 nbMapTasks int
 nbReduceTasks int
 tasks map[string]models.Status
 workers map[int]string
 mu sync.Mutex)

type Work int

func (w *Work) Tasks(t bool, reply *map[string]models.Status) error {
	*reply = tasks
	return nil
}

/*func CheckTaskDone(tasks map[int][]Status, task Task, workerId int) bool {
	time.Sleep(time.Second * 10)
	statusList := tasks[workerId]
	for _, s := range statusList {
		if s.Task.Name == task.Name && s.Task.Number == task.Number {
			if s.End.IsZero() {
				mapTasks <- task
				return false
			}
		}
	}
	return true
}*/

func (w *Work) GetTask(taskRequest int, task *models.Task) error {
	time.Sleep(time.Second*2)
	mu.Lock()
	defer mu.Unlock()
	if !allMapTasksDone {
		select {
		case mapTask, ok := <-mapTasks:
			if ok {
				// Return map task
				*task = mapTask
				fmt.Println("hedhi okhtna task",task)
				fmt.Println("Map task given to", taskRequest)
				st := models.Status{Start : time.Now().Format("15:04:05.000"), WorkerId: taskRequest}
				tasks[task.Name+strconv.Itoa(task.Number)] = st
				//go CheckTaskDone(tasks, mapTask, taskRequest)
				workers[st.WorkerId] = "Working"
				return nil
			} else {
				println("map TASKS DONE")
				allMapTasksDone = true
				workers[taskRequest] = "Waiting"
			}
		case <-time.After(time.Second * 5):
			task = new(models.Task)
		}
	}

	// Try reduce task if map tasks are done
	if allMapTasksDone {
		select {
		case reduceTask, ok := <-reduceTasks:
			if ok {	
				*task = reduceTask
				fmt.Println("hedhi okhtna task reduce ",task)
				fmt.Println("Reduce task given to", taskRequest)
				st := models.Status{Start: time.Now().Format("15:04:05.000"), WorkerId: taskRequest}
				tasks[task.Name+strconv.Itoa(task.Number)] = st
				//go CheckTaskDone(tasks, reduceTask, taskRequest)
				return nil
			}else{
				workers[taskRequest] = "Waiting"
			}
		case <-time.After(time.Second * 5):
			task = new(models.Task)
		}
	}
	return nil

	/*select {
	case mapTask, ok := <-mapTasks:
		if ok {
			*task = mapTask
			fmt.Println("map Task given to", taskRequest)
			st := Status{Task: mapTask, Start: time.Now()}
			mu.Lock()
			tasks[taskRequest] = append(tasks[taskRequest], st)
			mu.Unlock()
			go CheckTaskDone(tasks, mapTask, taskRequest)
		} else {
			println("map TASKS Done")
			// inputFiles channel is closed, signal that map phase is done
			allMapTasksDone = true
		}
	case r := <-reduceTasks:
		if allMapTasksDone {
			fmt.Println("reduce Task given to", taskRequest)
			st := Status{Task: r, Start: time.Now()}
			mu.Lock()
			tasks[taskRequest] = append(tasks[taskRequest], st)
			mu.Unlock()
			go CheckTaskDone(tasks, r, taskRequest)
		}
	case <-time.After(time.Second * 5):
		allMapTasksDone = true
	default:
		close(mapTasks)
		task = new(Task)
		printTasks(tasks)
	}*/
}

func (w *Work) ReportTaskDone(worker models.Worker, _ *int) error {
	mu.Lock()
	task := worker.Task
	status := tasks[task.Name+strconv.Itoa(task.Number)]
	if status.Start != "" {  
		status.End = time.Now().Format("15:04:05.000")
	}
	
	mu.Unlock()

	tasks[task.Name+strconv.Itoa(task.Number)] = status
	return nil
}

/*func GetIntermidiatFiles(dir string, taskNumber int) []string {
	files := []string{}
	prefix := "mrtmp.map-"
	ex := string(taskNumber)
	entries, err := os.ReadDir(dir)
	checkError(err, "reading directory error")

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) && strings.HasSuffix(entry.Name(), ex) {
			fmt.Println(entry.Name())
			files = append(files, entry.Name())
		}
	}
	return files
}*/

func concatFiles(destination string, sources []string) error {
	// Créer ou ouvrir le fichier de destination
	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copier le contenu de chaque fichier source
	for _, src := range sources {
		srcFile, err := os.Open(src)
		if err != nil {
			return err
		}
		_, err = io.Copy(destFile, srcFile)
		srcFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

type KeyValue struct {
	Key   string
	Value string
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func main() {
	allMapTasksDone = false
	nbReduceTasks = 0
	tasks = make(map[string]models.Status)
	content = []string{"f1.txt", "f2.txt", "f3.txt", "f4.txt"}
	workers = make(map[int]string) 
	mapTasks = make(chan models.Task, len(content))
	reduceTasks = make(chan models.Task, nbR*len(content))
	nbMapTasks = 0
	for i, f := range content {
		mapTasks <- models.Task{Name: "map", Number: i, File: f, NbReduce: nbR}
		}	
	for i := 0; i < nbR; i++ {
		reduceTasks <- models.Task{Name: "reduce", Number: i, File: "", NbReduce: 0}
	}
	close(mapTasks)
	close(reduceTasks)
	work := new(Work)
	rpc.Register(work)
	go StartDashboardServer()

	listener, err := net.Listen("tcp", ":1234")
	checkError(err, "Listener Error")
	fmt.Println("serving on 1234")
	for {
		conn, err := listener.Accept()
		checkError(err, "Connexion Error")
		go rpc.ServeConn(conn)
	}
}
