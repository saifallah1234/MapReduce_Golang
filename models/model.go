package models
type Task struct {
	Name     string `json:"name"`    
	Number   int    `json:"number"`  
	File     string `json:"file"`     
	NbReduce int    `json:"nbReduce"` }
type Status struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	WorkerId int       `json:"workerId"`}
type Worker struct {
	Id     int  `json:"id"`
	Status string `json:"status"` 
	Task   Task `json:"task"`}

