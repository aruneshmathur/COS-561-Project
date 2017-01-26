package main

import (
  "bufio"
  "fmt"
  "strings"
  "os"
  "strconv"
)

func readTopology(path string) (map[uint64]bool, [][]uint64, error) {
  file, err := os.Open(path)
  if err != nil {
	return nil, nil, err
  }
  defer file.Close()
  scanner := bufio.NewScanner(file)
  topology := make(map[uint64]bool) 
  edges := [][]uint64{}
  for scanner.Scan() {
	line := strings.Split(scanner.Text(), "|")




	as0, err := strconv.ParseInt(line[0], 10, 64)
	if err != nil {
		panic(err)
	}
	as1, err := strconv.ParseInt(line[1], 10, 64)

	as0_u := uint64(as0)
	as1_u := uint64(as1)
	entry := []uint64{as0_u, as1_u}
	if line[3] != "-1" {
		edges = append(edges, entry)

		topology[as0_u] = false
		topology[as1_u] = false

		if err != nil {
			panic(err)
		}
	}




/*
	reln_0_int, err := strconv.Atoi(line[3])
	reln_1_int := -1 * reln_0_int

	reln_0, err := strconv.ParseInt(line[0]+line[1], 10, 64)
	reln_1, err := strconv.ParseInt(line[1]+line[0], 10, 64)

	reln[reln_0] = reln_0_int
	reln[reln_1] = reln_1_int
	*/
  }

  return topology, edges, scanner.Err()
}
func readAttackers(path string) (map[uint64]uint64) {
  file, err := os.Open(path)
  if err != nil {
	panic(err)
}
  defer file.Close()
  scanner := bufio.NewScanner(file)
  topology := make(map[uint64]uint64) 
  for scanner.Scan() {
	line := strings.Split(scanner.Text(), ":")
	as0_int, err := strconv.ParseInt(line[0], 10, 64)
	value_int, err := strconv.ParseInt(line[1], 10, 64)
	if err != nil {
		panic(err)
	}
	topology[uint64(as0_int)] = uint64(value_int)
  	
	}

  return topology
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: To Few Args")
		panic(nil)
	}
	attackers := readAttackers(os.Args[2])


	for key := range(attackers) {
		if attackers[key] < 20 {
			delete(attackers, key)
		}
	}
	fmt.Println(len(attackers))
	for key := range(attackers) {
		wg.Add(1)
		go func() {
			path, _ := gr.Bfs(key)
			defenders := []uint64{}
			path_2, _ := gr_2.Bfs2(key)

			for j := range(path) {
				if _, ok := path_2[j]; ok {
					defenders = append(defenders, j)
					continue
				}
			}

			map.Set(key, defenders)
			defenders = append(defenders, key)
			channel <- defenders
			wg.Done()
		} ()
	}
	go func() {
		for {

			data := <-channel
			for i := 0; i < len(data) - 1; i++ {
				fmt.Println(data[len(data)-1],"|",data[i])
			}
		}

	} ()
	wg.Wait()
}
