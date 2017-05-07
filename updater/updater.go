package updater

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"encoding/hex"

	"sync"

	"encoding/json"

	"github.com/UnnoTed/fileb0x/file"
	"github.com/cheggaaa/pb"
)

// Auth holds authentication for the http basic auth
type Auth struct {
	Username string
	Password string
}

// ResponseInit holds a list of hashes from the server
// to be sent to the client so it can check if there
// is a new file or a changed file
type ResponseInit struct {
	Success bool
	Hashes  map[string]string
}

// Updater sends files that should be update to the b0x server
type Updater struct {
	Server string
	Auth   Auth

	RemoteHashes map[string]string
	LocalHashes  map[string]string
	ToUpdate     []string
}

// Init gets the list of file hash from the server
func (up *Updater) Init() error {
	return up.Get()
}

// Get gets the list of file hash from the server
func (up *Updater) Get() error {
	log.Println("Creating hash list request...")
	req, err := http.NewRequest("GET", up.Server, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(up.Auth.Username, up.Auth.Password)

	log.Println("Sending hash list request...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("Error Unautorized")
	}

	log.Println("Reading hash list response's body...")
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	log.Println("Parsing hash list response's body...")
	ri := &ResponseInit{}
	err = json.Unmarshal(buf.Bytes(), &ri)
	if err != nil {
		log.Println("Body is", buf.Bytes())
		return err
	}
	resp.Body.Close()

	// copy hash list
	if ri.Success {
		log.Println("Copying hash list...")
		up.RemoteHashes = ri.Hashes
		up.LocalHashes = map[string]string{}
		log.Println("Done")
	}

	return nil
}

// Updatable checks if there is any file that should be updaTed
func (up *Updater) Updatable(files map[string]*file.File) (bool, error) {
	hasUpdates := !up.EqualHashes(files)

	if hasUpdates {
		log.Println("----------------------------------------")
		log.Println("-- Found files that should be updated --")
		log.Println("----------------------------------------")
	} else {
		log.Println("-----------------------")
		log.Println("-- Nothing to update --")
		log.Println("-----------------------")
	}

	return hasUpdates, nil
}

// EqualHash checks if a local file hash equals a remote file hash
// it returns false when a remote file hash isn't found (new files)
func (up *Updater) EqualHash(name string) bool {
	hash, existsLocally := up.LocalHashes[name]
	_, existsRemotely := up.RemoteHashes[name]
	if !existsRemotely || !existsLocally || hash != up.RemoteHashes[name] {
		if hash != up.RemoteHashes[name] {
			log.Println("Found changes in file: ", name)

		} else if !existsRemotely && existsLocally {
			log.Println("Found new file: ", name)
		}

		return false
	}

	return true
}

// EqualHashes builds the list of local hashes before
// checking if there is any that should be updated
func (up *Updater) EqualHashes(files map[string]*file.File) bool {
	for _, f := range files {
		log.Println("Checking file for changes:", f.Path)

		if len(f.Bytes) == 0 && !f.ReplacedText {
			data, err := ioutil.ReadFile(f.OriginalPath)
			if err != nil {
				log.Fatal(err)
			}

			f.Bytes = data

			// removes the []byte("") from the string
			// when the data isn't in the Bytes variable
		} else if len(f.Bytes) == 0 && f.ReplacedText && len(f.Data) > 0 {
			f.Data = strings.TrimPrefix(f.Data, `[]byte("`)
			f.Data = strings.TrimSuffix(f.Data, `")`)
			f.Data = strings.Replace(f.Data, "\\x", "", -1)

			var err error
			f.Bytes, err = hex.DecodeString(f.Data)
			if err != nil {
				log.Println("SHIT", err)
				return false
			}

			f.Data = ""
		}

		sha := sha256.New()
		if _, err := sha.Write(f.Bytes); err != nil {
			log.Fatal(err)
			return false
		}

		up.LocalHashes[f.Path] = hex.EncodeToString(sha.Sum(nil))
	}

	// check if there is any file to update
	update := false
	for k := range up.LocalHashes {
		if !up.EqualHash(k) {
			up.ToUpdate = append(up.ToUpdate, k)
			update = true
		}
	}

	return !update
}

// ProgressReader implements a io.Reader with a Read
// function that lets a callback report how much
// of the file was read
type ProgressReader struct {
	io.Reader
	Reporter func(r int64)
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Reporter(int64(n))
	return
}

// UpdateFiles sends all files that should be updated to the server
// the limit is 3 concurrent files at once
func (up *Updater) UpdateFiles(files map[string]*file.File) error {
	updatable, err := up.Updatable(files)
	if err != nil {
		return err
	}

	if !updatable {
		return nil
	}

	var wg sync.WaitGroup

	// progressbar pool
	bpool, err := pb.StartPool()
	if err != nil {
		panic(err)
	}

	defer bpool.Stop()
	var current int64

	// Let's handle the clients asynchronously
	for _, name := range up.ToUpdate {
		for current == 3 {
		}

		current++
		wg.Add(1)

		go func(name string, wg *sync.WaitGroup) {
			defer wg.Done()
			defer func() {
				current--
			}()

			f := files[name]
			fr := bytes.NewReader(f.Bytes)

			bar := pb.New64(fr.Size()).SetUnits(pb.U_BYTES)
			bar.ShowSpeed = true
			bar.ShowTimeLeft = true
			bar.Prefix(f.Path)
			defer bar.Finish()

			// insert the bar into the progressbar pool
			bpool.Add(bar)

			// updates the progressbar
			pr := &ProgressReader{fr, func(r int64) {
				bar.Add64(r)
			}}

			r, w := io.Pipe()
			writer := multipart.NewWriter(w)

			// copy the file into the form
			go func(pr *ProgressReader) {
				defer w.Close()
				part, err := writer.CreateFormFile("file", f.Path)
				if err != nil {
					log.Fatal(err)
				}

				_, err = io.Copy(part, pr)
				if err != nil {
					log.Fatal(err)
				}

				err = writer.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(pr)

			// create a post request with basic auth
			// and the file included in a form
			req, err := http.NewRequest("POST", up.Server, r)
			if err != nil {
				log.Fatal(err)
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.SetBasicAuth(up.Auth.Username, up.Auth.Password)

			// sends the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)

			} else {
				body := &bytes.Buffer{}
				_, err := body.ReadFrom(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				resp.Body.Close()

				if body.String() != "ok" {
					log.Fatal(body.String())
				}
			}

		}(name, &wg)
	}

	wg.Wait()
	return nil
}
