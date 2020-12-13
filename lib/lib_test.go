package lib

import (
	"io/ioutil"
	"os"
	"testing"
)

var testFiles = []string{
	"test.key",
	"test.key.shard0",
	"test.key.shard1",
	"test.key.shard2",
	"test.key.shard3",
	"data.txt",
	"data.txt.shush",
}

func deleteTestFiles() {
	for _, f := range testFiles {
		os.Remove(f)
	}
}

func TestGen_Split_Merge(t *testing.T) {
	t.Cleanup(deleteTestFiles)
	deleteTestFiles()

	// generate a fresh key
	err := Gen("test.key")
	if err != nil {
		t.Fatal(err)
	}

	// remember its value for later
	original, err := ioutil.ReadFile("test.key")
	if err != nil {
		t.Fatal(err)
	}

	// generate shards
	err = Split("test.key", 4, 2)
	if err != nil {
		t.Fatal(err)
	}

	// delete the original key
	os.Remove("test.key")

	// merge the shards into a new key
	err = Merge([]string{"test.key.shard0", "test.key.shard1", "test.key.shard2"})
	if err != nil {
		t.Fatal(err)
	}

	recovered, err := ioutil.ReadFile("test.key")
	if err != nil {
		t.Fatal(err)
	}

	if string(recovered) != string(original) {
		t.Fatalf("%s (original) does not equal %s", original, recovered)
	}
}

const testData = "this is my test data!"

func TestEncrypt_Decrypt(t *testing.T) {
	t.Cleanup(deleteTestFiles)
	deleteTestFiles()

	// generate a fresh key
	err := Gen("test.key")
	if err != nil {
		t.Fatal(err)
	}

	// write dummy data
	err = ioutil.WriteFile("data.txt", []byte(testData), 0600)
	if err != nil {
		t.Fatal(err)
	}

	// encrypt dummy data
	err = Encrypt("test.key", "data.txt")
	if err != nil {
		t.Fatal(err)
	}

	// oops we "lost" the original data
	err = os.Remove("data.txt")
	if err != nil {
		t.Fatal(err)
	}

	// try to recover
	err = Decrypt("test.key", "data.txt.shush")
	if err != nil {
		t.Fatal(err)
	}

	result, err := ioutil.ReadFile("data.txt")
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != testData {
		t.Fatal("decrypted data doesn't match what we encrypted")
	}
}

func TestBase64Encode(t *testing.T) {
	result := base64encode([]byte("test"))
	if string(result) != "dGVzdA==" {
		t.Fatal("didn't encode correctly", string(result))
	}
}

func TestBase64Decode(t *testing.T) {
	result := base64decode([]byte("dGVzdA=="))
	if string(result) != "test" {
		t.Fatal("didn't decode correctly", string(result))
	}
}
