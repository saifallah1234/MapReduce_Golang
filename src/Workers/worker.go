package main

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"MapReduceProject/models"

)


const nbWorkers = 3
var rep map[string]models.Status



type Worker struct {
	Id     int
	Task   models.Task
	Status string
}

func printStatusMap(statusMap map[string]models.Status) {
	fmt.Println("-----------------------------------")
	for key, status := range statusMap {

		fmt.Printf("Key: %s\n", key)
		fmt.Printf("  Start   : %s\n", status.Start)
		fmt.Printf("  End     : %s\n", status.End)
		fmt.Printf("  WorkerId: %d\n", status.WorkerId)

	}
	fmt.Println("-----------------------------------")
}

func (w *Worker) Work(wg *sync.WaitGroup, client *rpc.Client) error {
	defer wg.Done()
	for {
		var t models.Task
		err := client.Call("Work.GetTask", w.Id, &t)
		checkError(err, "client error")
		w.Task = t
		fmt.Printf("worker %d : %+v \n", w.Id, w.Task)
		if w.Task.Name == "map" && w.Task.File != "" {
			DoMap(w.Task.Name, w.Task.Number, w.Task.File, w.Task.NbReduce, mapF)

			var rep int
			err = client.Call("Work.ReportTaskDone", w, &rep)
			checkError(err, "client error")
			continue
		} else {
			if w.Task.Name == "reduce" && w.Task.File == "" {

				DoReduce(w.Task.Name, w.Task.Number, 4, reduceF)

				println("doreduce done")
				var rep int
				err = client.Call("Work.ReportTaskDone", w, &rep)
				checkError(err, "client error")
				continue
			} else {
				err = client.Call("Work.Tasks", true, &rep)
				checkError(err, "client error")
				println("BREAK DONE")
				//printTasks(rep)
				break
			}
		}
	}
	return nil
}

func checkError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %v", msg, err)
    }
}


func main() {
	var waitGroup sync.WaitGroup
	client, err := rpc.Dial("tcp", "localhost:1234")
	checkError(err, "connexion error")
	workers := []Worker{}
	for i := 0; i < nbWorkers; i++ {
		workers = append(workers, Worker{Id: i, Task: models.Task{}})
	}

	for _, w := range workers {
		waitGroup.Add(1)
		go w.Work(&waitGroup, client)
	}
	waitGroup.Wait()
	printStatusMap(rep)
	MergeAllResults("reduce", 4, 4)

}
