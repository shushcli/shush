package shards

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestCreate(t *testing.T) {
	err := Create("test.txt", 4, 2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRecover(t *testing.T) {
	err := Recover("test.txt.shard*")
	if err != nil {
		t.Fatal(err)
	}

	cwd, err := os.Getwd()
	bytes, err := ioutil.ReadFile(path.Join(cwd, "test-recovered.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != "This is an example secret file. Yay!" {
		t.Fatal("didn't recover as expected")
	}
}

func TestBase64Encode(t *testing.T) {
	result := base64encode([]byte("test"))
	if string(result) != "dGVzdA==" {
		t.Fatal("Didn't encode correctly", string(result))
	}
}

func TestBase64Decode(t *testing.T) {
	result := base64decode([]byte("dGVzdA=="))
	if string(result) != "test" {
		t.Fatal("Didn't decode correctly", string(result))
	}
}
