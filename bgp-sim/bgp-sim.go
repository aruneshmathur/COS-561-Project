package main

import (
  "bufio"
  "fmt"
  "strings"
  "os"
  "strconv"
  "sync"
)

func readTopology(path string) (map[string][]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  edges := make(map[string][]string) //Map with keys AS 0 and values AS 1
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := strings.Split(scanner.Text(), "|")
    as0, as1 := line[0], line[1]
    edges[as1] = append(edges[as1], as0)
    edges[as0] = append(edges[as0], as1) 
  }
  return edges, scanner.Err()
}

//POLICY MESSAGE SPEC: "prefix|half-life|max-supress|reuse|suppress"

func readPolicy(path string) ([]int, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  line_int := []string{}
  for scanner.Scan() {
    line := strings.Split(scanner.Text(), "|")
    for _, i := range line {
        j, err := strconv.Atoi(i)
        if err != nil {
            return nil, err
        }
        line_int = append(line_int, j)
    }
  }
  return line_int, scanner.Err()
}

func readPolicies(edges map[string][]string, directory string) (map[string][]int, error) {
  	policies := make(map[string][]string) //Map with keys AS 0 and values AS 1
  	for as_name := range edges {
  		filename := directory + as_name + ".as"
  		policy, err := readPolicy()
  		if err != nil {
  			continue
  		}
  		policies[as_name] = policy
  	}
 	return policies, nil
}

type bgpEntry struct {
	prefix 	string
	pref 	int
	route	[] string
	active	bool
	rfd_penalty	int
	rfd_supress	bool
}

func duplicateBgpEntryCheck(bgpEntry bgpEntry, bgpEntryArray []bgpEntry) (bool) {
	for i := 0; i < len
}

func addBgpEntry(as string, bgpEntry bgpEntry, topology map[string][]string, bgp_table map[string][]bgpEntry) {
	if duplicateBgpEntryCheck(bgpEntry, bgp_table[as]) {
		bgp_table[as] = append(bgp_table[as], bgpEntry)
	} else {
		if false bgp_table[as].rfd_supress {
			bgp_table[as].active = true
		}
	}
	bgpEntry_copy = deepcopy.Copy(bgpEntry)
	bgpEntry_copy.route = append(bgpEntry_copy.route, as)
	for i := 0; i < len(topology[as]); i++ {
		addBgpEntry(topology[as][i], bgpEntry_copy, topology, bgp_table)
	}
}

func delBgpEntry(as string, topology map[string][]string, bgp_table map[string][][]string) {
	
}

func supressBgpEntry(as string, topology map[string][]string, bgp_table map[string][][]string) {
	
}

func initializeBgpTables() {
	bgp_table := make(map[string][]string)

}

func updateRFD() {
	
}

func updateBgpTables(bgp_traffic <-chan int, policies map[string][]int, topology map[string][]int,  log_message chan<- int) {
  	bgp_table := make(map[string][]string)
 	for {
 		bgp_packet := <-bgp_traffic
 		parsed_packet := strings.Split(bgp_packet, "|")
    	src_as_number,prefix,ann,del,rfd_damp,rfd_release := parsed_packet[0], parsed_packet[1], parsed_packet[2], parsed_packet[3], parsed_packet[4], parsed_packet[5]
    	if(ann) {
    		peers := topology[src_as_number]
    		for peer := range peers {
    			if _, ok := bgp_table[peer]; ok {

    				continue
    			}
    		}
    	}
 	}
}
//BGP MESSAGE SPEC: "src_as_number|prefix|announce|delete|rfd_damp|rfd_release"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("ERROR: To Few Args")
		panic(err)

	}

    edges, err := readTopology(os.Args[1])
    if err != nil {
    	fmt.Println("ERROR: Import Topology")
     	panic(err)
    }

    policies, err := readPolicies(edges, os.Args[2])
    if err != nil {
    	fmt.Println("ERROR: Import AS Routing Policies")
    	panic(err)
    }   

    bgp_traffic := make(chan string)
    log_message := make(chan string)

    fmt.Println(edges["24"])


    go bgp_tables()

    messages <- "ping" 
    messages <- "ping1" 

}
