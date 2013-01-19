package utils

import (
    "fmt"
    "github.com/sdegutis/go.assert"
    "testing"
)

const (
    generateRndStringsCapacity   = 10000
    generateRndStringsBufferSize = 15
)

func TestGeneticKey(t *testing.T) {
    fmt.Println("test GeneticKey started...")
    defer fmt.Println("test GeneticKey done.")

    results := generateTestDataBy(t, func() (str string, err error) {
        str = GeneticKey()
        return
    })

    fmt.Printf("\t%v random strings by TestGeneticKey generated, start comparing...\n", generateRndStringsCapacity)

    assertNoDuplicatedValues(t, results)
}

func TestGenerateRandomString(t *testing.T) {
    fmt.Println("test GenerateRandomString started...")
    defer fmt.Println("test GenerateRandomString done.")

    results := generateTestDataBy(t, func() (str string, err error) {
        str, err = GenerateRandomString(generateRndStringsBufferSize)
        return
    })

    fmt.Printf("\t%v random strings by GenerateRandomString generated, start comparing...\n", generateRndStringsCapacity)

    assertNoDuplicatedValues(t, results)
}

func generateTestDataBy(t *testing.T, f func() (string, error)) (results []string) {
    results = make([]string, generateRndStringsCapacity)
    for i := 0; i < generateRndStringsCapacity; i++ {
        str, err := f()
        assert.Equals(t, nil, err)

        results[i] = str
    }
    return
}

func assertNoDuplicatedValues(t *testing.T, values []string) {
    capacity := len(values)
    for i := 0; i < capacity; i++ {
        str1 := values[i]
        for j := i + 1; j < capacity; j++ {
            str2 := values[j]

            assert.NotEquals(t, str2, str1)
        }

        if i > 0 && i%500 == 0 {
            fmt.Printf("\t%v values campared...\n", i)
        }
    }
}
