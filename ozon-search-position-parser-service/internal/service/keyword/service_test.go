package keyword

import "testing"

func TestList(t *testing.T) {
    service := NewService()
    list := service.List()
    if len(list) == 0 {
        t.Error("Service return nothing")
    }

    firstKeyword := list[0]
    expected := "Keyword 1"
    if firstKeyword.Name != expected {
        t.Errorf("Got name %q, expected %q ", firstKeyword.Name, expected)
    }
}
