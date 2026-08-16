package identity

import("path/filepath";"testing")
func TestSaveAndLoad(t *testing.T){path:=filepath.Join(t.TempDir(),"credentials.json");want:=Credentials{AgentID:"agent-1",Token:"secret"};if err:=Save(path,want);err!=nil{t.Fatal(err)};got,err:=Load(path);if err!=nil{t.Fatal(err)};if got!=want{t.Fatalf("got %#v want %#v",got,want)}}
