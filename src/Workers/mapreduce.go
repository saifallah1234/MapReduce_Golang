package main

import (
	"encoding/json"
	"hash/fnv"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const prefix = "mrtmp."

// KeyValue is a type used to hold the key/value pairs passed to the map and
// reduce functions.
type KeyValue struct {
	Key   string
	Value string
}

func MergeAllResults(jobName string, nReduce, nMap int) {
	finalFile, err := os.Create(AnsName(jobName))
	if err != nil {
		log.Fatalf("erraur: %v", err)
	}
	defer finalFile.Close()

	encoder := json.NewEncoder(finalFile)

	for i := 0; i < nReduce; i++ {
		fileName := MergeName(jobName, i)
		f, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("erreur: %v", err)
		}

		decoder := json.NewDecoder(f)
		var kv KeyValue
		for decoder.Decode(&kv) == nil {
			err := encoder.Encode(&kv)
			if err != nil {
				log.Fatalf("erreur: %v", err)
			}
		}
		f.Close()
	}
	CleanIntermediary("map", nMap, nReduce)
	CleanIntermediary(jobName, nMap, nReduce)
}

// reduceName constructs the name of the intermediate file which map task
// <mapTask> produces for reduce task <reduceTask>.
func ReduceName(jobName string, mapTask int, reduceTask int) string {
	return prefix + jobName + "-" + strconv.Itoa(mapTask) + "-" + strconv.Itoa(reduceTask)
}

// mergeName constructs the name of the output file of reduce task <reduceTask>
func MergeName(jobName string, reduceTask int) string {
	return prefix + jobName + "-res-" + strconv.Itoa(reduceTask)
}

// ansName constructs the name of the output file of the final answer
func AnsName(jobName string) string {
	return prefix + jobName
}

// clean all intermediary files generated for a job
func CleanIntermediary(jobName string, nMap, nReduce int) {
	// Supprimer les fichiers intermédiaires produits les tâches map
	for reduceTNbr := 0; reduceTNbr < nReduce; reduceTNbr++ {
		for mapTNbr := 0; mapTNbr < nMap; mapTNbr++ {
			os.Remove(ReduceName(jobName, mapTNbr, reduceTNbr))
		}
		os.Remove(MergeName(jobName, reduceTNbr))
	}
}

// Is used to associate to each key a unique reduce file
func ihash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
func mapF(contents string) (kvs []KeyValue) {
	// Exemple simple de split et comptage
	words := strings.Split(contents, " ")
	kvs = []KeyValue{}
	for _, w := range words {
		w1 := strings.ToLower(w)
		kvs = append(kvs, KeyValue{Key: w1, Value: "1"})

	}
	return
}

// doMap applique la fonction mapF, et sauvegarde les résultats.
// A COMPLETER
func DoMap(
	jobName string,
	mapTaskNumber int,
	inFile string,
	nReduce int,
	mapF func(contents string) []KeyValue,
) {
	// Lire le contenu du fichier d'entrée
	content, err := os.ReadFile(inFile)
	if err != nil {
		log.Fatalf(inFile, err)
	}
	// Appliquer la fonction mapF pour obtenir des paires clé-valeur
	kvs := mapF(string(content))

	// Créer nReduce fichiers, un pour chaque tâche de réduction
	// utiliser reduceName
	files := make([]*os.File, nReduce)
	encoders := make([]*json.Encoder, nReduce)

	for i := 0; i < nReduce; i++ {
		fileName := ReduceName(jobName, mapTaskNumber, i)
		f, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("erreur %v", err)
		}
		files[i] = f
		encoders[i] = json.NewEncoder(f)
		defer f.Close()

	}
	// Partitionner les paires clé-valeur en fonction du hachage de la clé
	for _, kv := range kvs {
		r := int(ihash(kv.Key)) % nReduce
		err := encoders[r].Encode(&kv)
		if err != nil {
			log.Fatalf("erreur %v", err)
		}
	}
	// Calculer la tâche de réduction associée à chaque clé
	// Écrire la paire clé-valeur dans le fichier approprié

}
func reduceF(key string, values []string) string {
	count := 0
	for _, val := range values {
		n, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("erreur")
		}
		count += n
	}
	return strconv.Itoa(count)
}

// doReduce effectue une tâche de réduction en lisant les fichiers
// intermédiaires, en regroupant les valeurs par clé, et en appliquant
// la fonction reduceF.
// A COMPLETER

func DoReduce(
	jobName string,
	reduceTaskNumber int,
	nMap int,
	reduceF func(key string, values []string) string,
) {
	time.Sleep(time.Second*3)
	// Map temporaire pour regrouper les valeurs par clé
	intermediate := make(map[string][]string)

	// Lire les fichiers intermédiaires générés par chaque map
	for i := 0; i < nMap; i++ {
		fileName := ReduceName("map", i, reduceTaskNumber)
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatalf(fileName, err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		var kv KeyValue
		for decoder.Decode(&kv) == nil {
			intermediate[kv.Key] = append(intermediate[kv.Key], kv.Value)
			
		}
	}

	// Créer le fichier de sortie
	outputFileName := MergeName(jobName, reduceTaskNumber)
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("erreur %v", err)
	}
	defer outputFile.Close()
	enc := json.NewEncoder(outputFile)

	keys := make([]string, 0, len(intermediate))
	for k := range intermediate {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Appliquer reduceF à chaque clé, puis écrire le résultat
	for _, key := range keys {
		reducedValue := reduceF(key, intermediate[key])

		err := enc.Encode(&KeyValue{Key: key, Value: reducedValue})
		if err != nil {
			log.Fatalf("erreur %v", err)
		}
	}
}
