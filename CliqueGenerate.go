package main

// input, edge between v1 and v2
// output, maximal clique
import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"git.garena.com/common/gocommon/goio"
	"git.garena.com/common/gocommon/goutil"
	"github.com/bradfitz/slice"
)

type array struct {
	clique []string
	length int
}

var (
	MinCliqueSize  = flag.Int("MinCliqueSize", 3, "min size for each clique")
	Neighbor       = make(map[string]([]string))
	DegneracyOrder = []string{}
	CliqueRecord   = []array{}
)

func union(cluster1 []string, cluster2 []string) []string {
	cluster := append(cluster1, cluster2...)
	mmap := goutil.ListToMapString(cluster)
	cluster = goutil.MapToListString(mmap)
	return cluster
}

func intersection(cluster1 []string, cluster2 []string) []string {
	intersec := []string{}
	map2 := goutil.ListToMapString(cluster2)
	for _, ele := range cluster1 {
		if _, find := map2[ele]; find {
			intersec = append(intersec, ele)
		}
	}
	return intersec
}

func subtraction(leftCluster []string, rightCluster []string) []string {
	// sub rightCluster from leftCluster
	map1 := goutil.ListToMapString(leftCluster)
	for _, ele := range rightCluster {
		if _, find := map1[ele]; find {
			delete(map1, ele)
		}
	}
	leftCluster = goutil.MapToListString(map1)
	return leftCluster
}

func degeneracyTrase() {
	for ind, node := range DegneracyOrder {
		neighborNode := Neighbor[node]
		candidates := intersection(neighborNode, DegneracyOrder)
		exclued := intersection(neighborNode, DegneracyOrder[:ind])
		clique := []string{node}
		bronkerKerbosch(clique, candidates, exclued)
	}
}

func randomPick(cluster []string) string {
	length := len(cluster)
	if length < 1 {
		return ""
	}
	ind := rand.Intn(length)
	return cluster[ind]
}

func formatPrint(clique []string) {
	length := len(clique)
	CliqueRecord = append(CliqueRecord, array{clique, length})
}

func bronkerKerbosch(clique []string, candidates []string, exclued []string) {
	// reported size ++
	if len(candidates) == 0 && len(exclued) == 0 {
		if len(clique) > *MinCliqueSize {
			formatPrint(clique)
		}
	}
	if len(union(candidates, exclued)) == 0 {
		return
	}
	node := randomPick(union(candidates, exclued))
	neighborNode := Neighbor[node]
	for _, vertice := range subtraction(candidates, neighborNode) {
		neighborVertice := Neighbor[vertice]
		new_candidate := intersection(candidates, neighborVertice)
		new_exclued := intersection(exclued, neighborVertice)
		new_clique := union(clique, []string{vertice})
		bronkerKerbosch(new_clique, new_candidate, new_exclued)
		candidates = subtraction(candidates, []string{vertice})
		exclued = union(exclued, []string{vertice})
	}
}

func processEdge(vertice1 string, vertice2 string, weight int) {
	/*
		if _, find := Edge[vertice1]; !find {
			Edge[vertice1] = make(map[string]int)
		}
		Edge[vertice1][vertice2] = weight
	*/
	if _, find := Neighbor[vertice1]; !find {
		Neighbor[vertice1] = []string{}
	}
	if _, find := Neighbor[vertice2]; !find {
		Neighbor[vertice2] = []string{}
	}
	Neighbor[vertice1] = append(Neighbor[vertice1], vertice2)
	Neighbor[vertice2] = append(Neighbor[vertice2], vertice1)
}

func postProcess() {
	// 1-step remove duplicate
	// 2-step get degeneracy ordering of vertice
	deg := make(map[string]int)
	for node, neighbor := range Neighbor {
		neighbor = goutil.RemoveDuplicateString(neighbor)
		Neighbor[node] = neighbor
		deg[node] = len(neighbor)
	}
	// get degeneracy ordering sequence

	for len(deg) > 0 {
		minn := 99999
		node := ""
		for vertice, outDegree := range deg { // find min outDegree vertice
			if outDegree < minn {
				minn = outDegree
				node = vertice
			}
		}

		// del the vertice and connect edge with this vertice

		DegneracyOrder = append(DegneracyOrder, node)
		delete(deg, node)

		for _, vertice := range Neighbor[node] {

			if _, find := deg[vertice]; find {
				deg[vertice] = deg[vertice] - 1
			}
		}
	}

	fmt.Println("this is the generacy order of current node")
	fmt.Println(DegneracyOrder)
}

func main() {

	flag.Parse()
	for input := range goio.NewInput(os.Stdin) {
		tks := strings.Split(input, "\t")
		if len(tks) < 3 {
			continue
		}
		num, _ := strconv.Atoi(tks[2])
		processEdge(tks[0], tks[1], num)
	}
	postProcess()
	degeneracyTrase()

	slice.Sort(CliqueRecord[:], func(i, j int) bool {
		return CliqueRecord[i].length > CliqueRecord[j].length
	})
	NewCliqueRecord := []array{}
	for ind, arr := range CliqueRecord {
		cliqueS := arr.clique
		mark := false
		for i := 0; i < ind; i++ {
			cliqueL := CliqueRecord[i].clique
			interc := intersection(cliqueS, cliqueL)
			if goutil.SetEqual(interc, cliqueS) {
				mark = true
				break
			}
		}
		if mark == false {
			NewCliqueRecord = append(NewCliqueRecord, arr)
		}
	}

	for _, arr := range NewCliqueRecord {
		fmt.Println(arr.clique)
	}
}
