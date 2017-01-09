package main

import (
  "bufio"
  "fmt"
  "strings"
  "os"
  "time"
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

type asPolicyEntry struct {
	prefix 			string
	rfd_penalty 	int
	half-life 		int
	max-supress 	int
	reuse			int
	supress 		int
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
	available bool
	rfd_penalty	int
	rfd_supress	bool
	rfd_time_reset bool
}

func searchBgpEntryPrefix(bgpEntry bgpEntry, bgpEntryArray []bgpEntry) (int, bool) {
	for i := 0; i < len(bgpEntryArray); i++ {
		if(bgpEntry.prefix == bgpEntryArray[i].prefix) {
			return index, true
		}
	}
	return 0, false
}

func searchBgpEntryActivePrefix(bgpEntry bgpEntry, bgpEntryArray []bgpEntry) (int, bool) {
	for i := 0; i < len(bgpEntryArray); i++ {
		if(bgpEntry.prefix == bgpEntryArray[i].prefix &&
			bgpEntryArray[i]) {
			return index, true
		}
	}
	return 0, false
}

func searchBgpEntry(bgpEntry bgpEntry, bgpEntryArray []bgpEntry) (int, bool) {
	for i := 0; i < len(bgpEntryArray); i++ {
		if bgpEntry.prefix == bgpEntryArray[i].prefix {
			if len(bgpEntryArray[i].route) == len(bgpEntry.route) {
				for j := 0; j < len(bgpEntry.route); j++ {
					if bgpEntryArray[i].route[j] != bgpEntry.route[j] {
						return 0, false
					}
				}
				return i, true
			}
		}
	}
	return 0, false
}

func duplicateBgpEntryCheck(bgpEntry bgpEntry, bgpEntryArray []bgpEntry) (bool) {
	index, duplicate := searchBgpEntry(bgpEntry, bgpEntryArray)
	return duplicate
}

func announceBgpEntry(as string, bgpEntry bgpEntry, topology map[string][]string, bgp_table map[string][]bgpEntry) {
	bgpEntry_copy = deepcopy.Copy(bgpEntry)
	bgpEntry_copy.route = append(bgpEntry_copy.route, as)
	for i := 0; i < len(topology[as]); i++ {
		addBgpEntry(topology[as][i], bgpEntry_copy, topology, bgp_table)
	}
}

func updateBgpEntries(bgpEntry bgpEntry, bgpEntryArray []bgpEntry) {
	highest_pref := -1 
	index := -1
	for i := 0; i < len(bgpEntryArray); i++ {
		if bgpEntry.prefix == bgpEntryArray[i].prefix &&
			bgpEntryArray[i].enabled && 
			bgpEntryArray[i].rfd_supress == false &&
			bgpEntryArray[i].pref > highest_pref{
			highest_pref = bgpEntryArray[i].pref
			index = i
		}
	}
	if index != -1 {
		if false bgpEntryArray[index].active {
			active_prefix_index, prefix_active := searchBgpEntryActivePrefix(bgpEntry, bgpEntryArray)
			if prefix_active {
				withdrawBgpEntry(  )
			} 
			announceBgpEntry( )
		}
	}
}

func addBgpEntry(as string, bgpEntry bgpEntry,policies map[string][]int, topology map[string][]string, bgp_table map[string][]bgpEntry) {

	if false duplicateBgpEntryCheck(bgpEntry, bgp_table[as]) {
		bgp_table[as] = append(bgp_table[as], bgpEntry)
		go manageRfd(bgp_table[as][len(bgp_table[as]-1)], policies[as])
	} 
	index, _ :=searchBgpEntry(bgpEntry, bgp_table[as])
	bgp_table[as][index].pref = bgpEntry.pref 
	bgp_table[as][index].enabled = true
	
	updateBgpEntries(bgpEntry, bgpEntryArray)
}

func withdrawBgpEntry(as string, policies map[string][]int, topology map[string][]string, bgp_table map[string][][]string, root_as bool) {
	index, found := searchBgpEntry()
	if found {
		bgp_table[as][index].enabled = false
		bgp_table[as][index].active = false
		bgp_table[as][index].rfd_penalty = bgp_table[as][index].rfd_penalty + policies[as].rfd_penalty
		if bgp_table[as][index].rfd_penalty > policies[as].suppress {
			bgp_table[as][index].rfd_supress = true
			bgp_table[as][index].rfd_time_reset = true
		}
		updateBgpEntries() //ACTIVATE OTHER ROUTE POTENTIALLY 
		bgpEntry_copy = deepcopy.Copy(bgp_table[as][index])
		bgpEntry_copy.route = bgpEntry_copy.route[:len(bgpEntry_copy.route)-1]
		for i := 0; i < len(topology[as]); i++ {
			withdrawBgpEntry(topology[as][i], bgpEntry_copy, topology, bgp_table)
		}
	}
}

func manageRfd(bgpEntry bgpEntry, as_policy asPolicyEntry) {
	timer := 0
	for {
		if bgpEntry.rfd_time_reset {
			timer = 0
			bgpEntry.rfd_time_reset = false
		}
		if timer%asPolicyEntry.half-life == 0 {
			bgpEntry.rfd_penalty = bgpEntry.rfd_penalty/2
			if bgpEntry.rfd_penalty <= asPolicyEntry.suppress {
				bgpEntry.rfd_supress = false
				updateBgpEntries( )
			}
		}
		if timer%asPolicyEntry.max-supress == 0 {
			timer = 0
			bgpEntry.rfd_supress = false
			updateBgpEntries( )
		}
		time.Sleep(time.Minute * 1)
		timer++
	}
}

func initializeBgpTables() {
	bgp_table := make(map[string][]string)

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
