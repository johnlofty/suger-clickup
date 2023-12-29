package models

import "testing"

func TestExtractUsername(t *testing.T) {
	text := "(alvin)[1019234] says: content"
	username, ok := ExtractComment(text)
	t.Errorf("Extracted Username:%+v ok:%v", username, ok)
	text2 := "(alvin@gmail.com)[1019234] says: content"
	username, ok = ExtractComment(text2)
	t.Errorf("Extracted Username:%+v ok:%v", username, ok)
}
