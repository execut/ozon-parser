package main

import (
    "github.com/stretchr/testify/assert"
    "io"
    "os"
    "testing"
)

func captureOutput(f func()) string {
    orig := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    f()
    os.Stdout = orig
    w.Close()
    out, _ := io.ReadAll(r)
    return string(out)
}

func TestKeywordsListFunc(t *testing.T) {
    os.Args = []string{"test", "-action", "keywordsList"}

    result := captureOutput(func() {
        main()
    })

    assert.Equal(t, "Test keyword 1\nTest keyword 2\n", result)
}

func TestAddKeywordFunc(t *testing.T) {
    expected := "Test keyword"
    os.Args = []string{"test", "-action", "addKeyword", expected}

    result := captureOutput(func() {
        main()
    })

    assert.Equal(t, "Keyword \""+expected+"\" added\n", result)
}
