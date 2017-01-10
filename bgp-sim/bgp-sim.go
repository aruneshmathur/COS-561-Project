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

///### TODO: PSUDO FUNCTION CALLS USED IN SOME PLACES CURRENTLY ###
///### TODO: LOGGING TO CONSOLE ###

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

func readPolicy(path string) (asPolicyEntry, error) {
  file, err := os.Open(path)
  if err != nil {
  return nil, err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  line := []string{}
  line_int := []int{}

  for scanner.Scan() {
    line = strings.Split(scanner.Text(), "|")
  }
  prefix := line[len(line) - 1]
  line = line[:len(line) - 1]
  for _, i := range line {
    j, err := strconv.Atoi(i)
    if err != nil {
      return nil, err
    }
    line_int = append(line_int, j)
  }
  return asPolicyEntry{prefix: prefix, rfd_penalty: line_int[0], half-life: line_int[1], max-supress: line_int[2], reuse: line_int[3], supress: line_int[4]}, scanner.Err()
}

func readPolicies(edges map[string][]string, directory string) (map[string]asPolicyEntry, error) {
	policies := make(map[string]asPolicyEntry) //Map with keys AS 0 and values AS 1
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

func announceBgpEntry(as string, bgpEntry bgpEntry, policies map[string]asPolicyEntry, topology map[string][]string, bgp_table map[string][]bgpEntry) {
	bgpEntry_copy = deepcopy.Copy(bgpEntry)
	bgpEntry_copy.route = append(bgpEntry_copy.route, as)
	for i := 0; i < len(topology[as]); i++ {
		addBgpEntry(topology[as][i], bgpEntry_copy, policies, topology, bgp_table)
	}
}

func updateBgpEntries(as string, bgpEntry bgpEntry, policies map[string]asPolicyEntry, topology map[string][]string, bgp_table map[string][]bgpEntry) {
	highest_pref := -1 
	index := -1
    bgpEntryArray := bgp_table[as];
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
				withdrawBgpEntry(as, bgpEntry, policies, topology, bgp_table)
			} 
			announceBgpEntry(as, bgpEntry, policies, topology, bgp_table)
		}
	}
}

func addBgpEntry(as string, bgpEntry bgpEntry, policies map[string]asPolicyEntry, topology map[string][]string, bgp_table map[string][]bgpEntry) {
	if false duplicateBgpEntryCheck(bgpEntry, bgp_table[as]) {
		bgp_table[as] = append(bgp_table[as], bgpEntry)
		go manageRfd(as, bgp_table[as][len(bgp_table[as]-1)], policies, topology, bgp_table)
	} 
	index, _ :=searchBgpEntry(bgpEntry, bgp_table[as])
	bgp_table[as][index].pref = bgpEntry.pref 
	bgp_table[as][index].enabled = true
	
	updateBgpEntries(as, bgpEntry, policies, topology, bgp_table)
}

func withdrawBgpEntry(as string, bgpEntry bgpEntry, policies map[string]asPolicyEntry, topology map[string][]string, bgp_table map[string][]bgpEntry) {
	index, found := searchBgpEntry(bgpEntry, bgp_table[as])
	if found {
		bgp_table[as][index].enabled = false
		bgp_table[as][index].active = false
		bgp_table[as][index].rfd_penalty = bgp_table[as][index].rfd_penalty + policies[as].rfd_penalty
		if bgp_table[as][index].rfd_penalty > policies[as].suppress {
			bgp_table[as][index].rfd_supress = true
			bgp_table[as][index].rfd_time_reset = true
		}
		updateBgpEntries(as, bgpEntry, policies, topology, bgp_table) //ACTIVATE OTHER ROUTE POTENTIALLY 
		bgpEntry_copy = deepcopy.Copy(bgp_table[as][index])
		bgpEntry_copy.route = bgpEntry_copy.route[:len(bgpEntry_copy.route)-1]
		for i := 0; i < len(topology[as]); i++ {
			withdrawBgpEntry(topology[as][i], bgpEntry_copy, policies, topology, bgp_table)
		}
	}
}

func manageRfd(as string, bgpEntry bgpEntry, policies map[string]asPolicyEntry, topology map[string][]string, bgp_table map[string][]bgpEntry) {
	timer := 0
    asPolicyEntry := policies[as]
	for {
		if bgpEntry.rfd_time_reset {
			timer = 0
			bgpEntry.rfd_time_reset = false
		}
		if timer%asPolicyEntry.half-life == 0 {
			bgpEntry.rfd_penalty = bgpEntry.rfd_penalty/2
			if bgpEntry.rfd_penalty <= asPolicyEntry.suppress {
				bgpEntry.rfd_supress = false
				updateBgpEntries(as, bgpEntry, policies, topology, bgp_table)
			}
		}
		if timer%asPolicyEntry.max-supress == 0 {
			timer = 0
			bgpEntry.rfd_supress = false
			updateBgpEntries(as, bgpEntry, policies, topology, bgp_table)
		}
		time.Sleep(time.Minute * 1)
		timer++
	}
}

func initializeBgpTables(policies map[string]asPolicyEntry, topology map[string][]string, bgp_table map[string][]bgpEntry) {
	for i := range topology {
		newBgpEntry := bgpEntry{prefix: i, pref: 0, route: {i}, active: false, available: true, rfd_penalty: 0, rfd_supress: false, rfd_time_reset: false}
		addBgpEntry(newBgpEntry)
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

	bgp_table := make(map[string][]string)

	initializeBgpTables()

}
