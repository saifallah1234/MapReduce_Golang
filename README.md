# MapReduce_Golang
MapReduce_Golang is a distributed computing framework built in Go that implements the MapReduce programming paradigm. This system allows for the efficient processing of large data sets across multiple nodes by dividing tasks into map and reduce operations. The project also hints at a frontend interface, suggesting future user interactivity.

### Tools

<p> <img alt="Go" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white"/> 
<img alt="RPC" src="https://img.shields.io/badge/RPC-3E4E88?style=for-the-badge"/>
<img alt="JavaScript" src="https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black"/> 
<img alt="Tailwind CSS" src="https://img.shields.io/badge/TailwindCSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white"/> </p>

## Features

### Distributed Computing :
Implements MapReduce for large-scale data processing across multiple Go routines/nodes.

### Concurrency and Parallelism : 
Utilizes Go's goroutines and channels for efficient task distribution and parallel processing.

### Modular Architecture :
Designed with clear separation between the map, reduce, and coordination logic for maintainability.

### Frontend Ready :
Includes directories and setup for the Dashboard using Tailwind CSS and JavaScript 

## Data Flow
### Map Phase:
Each worker applies a map function to input data chunks.

### Shuffle Phase:
Intermediate results are grouped by key.

### Reduce Phase:
The reduce function is applied to each group to aggregate results.

## Project structure 

MapReduce_Golang/

├── go.mod

├──models/

├── src/

│   ├── master/  

│   ├── workers/      

│   └── shared/      

├── TailWindCSS_test/


## Setup

Clone repo then run this commands 
cd src/master
go run .

cd src/workers
go run .
