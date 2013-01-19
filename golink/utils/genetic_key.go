package utils

import (
    //"fmt";
    cryptoRand "crypto/rand"
    "encoding/base64"
    "github.com/QLeelulu/ohlala/golink"
    "io"
    "math/rand"
    "time"
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
    candidate += geneSet[geneIndex : 1+geneIndex]
    if parentIndex+1 < len(parent) {
        candidate += parent[parentIndex+1:]
    }
    return candidate
}

func generateParent(geneSet string, length int) string {
    s := ""
    for i := 0; i < length; i++ {
        index := rand.Intn(len(geneSet))
        s += geneSet[index : 1+index]
    }
    return s
}

// a more randomized method to generate random string
func GenerateRandomString(bufferSize uint32) (string, error) {
    buf := make([]byte, bufferSize)
    n, err := io.ReadFull(cryptoRand.Reader, buf)
    if n != len(buf) || err != nil {
        return "", err
    }

    str := ConvertByteArrayToBase64String(buf)
    return str, nil
}

func ConvertByteArrayToBase64String(buf []byte) string {
    encoding := base64.StdEncoding
    str := encoding.EncodeToString(buf)
    return str
}
