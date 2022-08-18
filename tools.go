package toolkit

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/*
This package in the toolkit folder will hold code that's
used when working on modules while developing.
*/
/*
Tools is a type to instanciate this module.
Any variable of this type will have access to all methods
trough receiver *Tools.
@MaxFileSize
		The maximum allowed size for an uploaded file. A value
		set by the developer using this module.
@ValidTypes
		Allowed uploadable file types to be specified by the
		developer using this module.
*/
type Tools struct {
	MaxFileSize int
	ValidTypes  []string
}

/*
Generate a string of random characters of length N.
===================================================
filenames < 100 characters.
forbidden characters are: <>:/\|?* and ASCI 1 - 31
@sourceStr
  - The source for generated random characters in the
    returned output string.
*/
const sourceStr = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890+&$%!"

func (t *Tools) RandomString(n int) string {
	if n < 100 {
		res, src := make([]rune, n), []rune(sourceStr)
		for i := range res {
			p, _ := rand.Prime(rand.Reader, len(src))
			x, y := p.Uint64(), uint64(len(src))
			res[i] = src[x%y]
		}
		return string(res)
	} else {
		return "N must be less than 100!"
	}
}

/*
Uploading one or more files to the server
=========================================
*
To send information back to whoever called the function
that handles the files uploading, we need a type.
*
The uploaded file originally has a name, but we don't
want to use that name; for one, it may be a duplicate
name and so overwrites an existing file. But we still
want to keep track of the original name. => WHY?
*
To retain the original filename seems the only reaison
d'etre for this type. => What do we want to do with it?
*
Actually this about feedBackData
*/
type FileData struct {
	NewFileName string
	OrgFileName string
	FileSize    int64
}

/*
The function should be available to anyone that creates
type Tools from the package toolkit, so we have the
receiver.
@params
r	*http.Request

	The POST made by the user comes with a request and
	we want to have access to that request.

uploadDir

	The directory where we store the uploaded file.

rename

	The user of our package sometimes may not want to
	rename their files by default

@returns

	A slice of pointers to file data. Needed(?) when
	the user uploads more then one file.
	error
*/
func (t *Tools) UploadFiles(r *http.Request, uploadDir string,
	rename ...bool) ([]*FileData, error) {
	// We can have one or more bools in rename, or notthing at all
	// We rename files by default.
	renameFile := true
	// But, there might be one or more value's in rename
	if len(rename) > 0 {
		// So we set rename to the first value in variatic
		renameFile = rename[0]
	}

	// Declare a return variable to store data about files;
	// a slice of pointers to struct FileData.
	var filesData []*FileData

	// t.MaxFileSize by default is zero, so we check to see
	// if set by the developer
	// if not, give a sensible default value.
	if t.MaxFileSize == 0 {
		// a reasonable size for a compressed image is 4 Mb.
		t.MaxFileSize = 1024 * 1024 * 4
	}

	// parses a request body as multipart/form-data
	err := r.ParseMultipartForm(int64(t.MaxFileSize))
	if err != nil {
		return nil, errors.New("uploaded file is too big")
	}

	// now go through the request and look for stored files
	for _, fHeaders := range r.MultipartForm.File {
		// next go through the headers of all uploaded files.
		for _, hdr := range fHeaders {
			// try to get a single file out of the request; the
			// first time the loop is executed, filesData = empty
			filesData, err = func(filesData []*FileData) ([]*FileData, error) {
				// create a place to store data about the file when 
				// we take it out of the request.
				var fileData FileData
				// can we actually open the uploaded file 'hdr'
				inFile, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				// OK, inFile contains a file; let's prevent a leak.
				defer inFile.Close()

				// Determine what kind of file we have. We need only
				// the first 512 bytes to take a peek into the header
				buff := make([]byte, 512)
				_, err = inFile.Read(buff)
				if err != nil {
					return nil, err
				}

				// check for permitted MIME type, to prevent scripts etc.
				allowed := false
				fileType := http.DetectContentType(buff) // image/png oid

				fmt.Println("fileType: ", fileType)			// fileType:  image/jpeg

				// validTypes are provided by the developer through Tools.
				// validTypes := []string{"image/jpeg image/png image/gif"}
				if len(t.ValidTypes) > 0 {
					// the developer HAS specified allowed file types
					fmt.Println("Types: ", t.ValidTypes)	// Types:  [image/jpeg, image/png, image/gif]
					for _, x := range t.ValidTypes {
						fmt.Println("x: ", x)								// x:  image/jpeg, image/png, image/gif
						if strings.EqualFold(fileType, x) {
							// the fileType of the actual file at hand is allowed
							allowed = true
						}
					} // for-loop
				} else {
					// the developer has NOT specified allowed file types.
					// tricky, but suit your selve.
					allowed = true
				}

				// check for a file type NOT allowed
				if !allowed {
					return nil, errors.New("uploaded file type NOT allowed")
				}

				// we've read the file, so we NEED to reset the cursor.
				_, err = inFile.Seek(0, 0)
				if err != nil {
					return nil, err
				}

				// handle renaming files (default setting)
				if renameFile {
					// generate a new filename (might NOT be unique)
					fileData.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(hdr.Filename))
				} else {
					fileData.NewFileName = hdr.Filename
				}

				// save the file to disk in the parameterized directory: THIS
				// is actually the core business of this function!
				var outFile *os.File
				defer outFile.Close()
				if outFile, err = os.Create(filepath.Join(uploadDir, fileData.NewFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outFile, inFile)
					if err != nil {
						return nil, err
					}
					fileData.FileSize = fileSize
				}

				// now we can append the file to the slice
				filesData = append(filesData, &fileData)

				return filesData, nil
			}(filesData)
			if err != nil {
				return filesData, err
			}
		} // for-loop
	} // for-loop
	return filesData, nil
}
