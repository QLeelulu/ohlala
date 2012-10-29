package utils

import (
	//"fmt";
	"math/rand";
	"time"
    "github.com/QLeelulu/ohlala/golink"
)

func GeneticKey() string {
	rand.Seed(time.Now().UnixNano())
	geneSet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bestParent := generateParent(geneSet, golink.Genetic_Key_Len)

	return mutateParent(bestParent, geneSet)
}

func mutateParent(parent, geneSet string) string {
	geneIndex := rand.Intn(len(geneSet))
	parentIndex := rand.Intn(len(parent))
	candidate := ""
	if parentIndex > 0 {
		candidate += parent[:parentIndex]
	}
	candidate += geneSet[geneIndex:1+geneIndex]
	if parentIndex+1 < len(parent) {
		candidate += parent[parentIndex+1:]
	}
	return candidate
}

func generateParent(geneSet string, length int) string {
	s := ""
	for i := 0; i < length; i++ {
		index := rand.Intn(len(geneSet))
		s += geneSet[index:1+index]
	}
	return s
}






