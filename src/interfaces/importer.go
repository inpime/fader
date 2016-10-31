package interfaces

import (
	"bufio"
	"errors"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"bytes"
	"fmt"
	"log"
	"os"

	"encoding/base64"
)

const (
	appname          = "FADER_ARCHIVE"
	archiveSeparator = "|"
)

type (
	partFile string
)

var (
	endbuckets               = "----end buckets----"
	infoPartFile    partFile = "partfile:firtsline"
	bucketsPartFile partFile = "partfile:buckets"
	filesPartFile   partFile = "partfile:files"
)

func NewImportManager(
	bm BucketManager,
	fm FileManager,
) *ImportManager {
	return &ImportManager{
		bm: bm,
		fm: fm,
		l:  log.New(os.Stdout, "[import manager]", -1),
	}
}

type ImportManager struct {
	bm BucketManager
	fm FileManager
	l  *log.Logger
}

func (m *ImportManager) Import(data []byte) (info ArchiveInfoLine, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var mode = infoPartFile

	for scanner.Scan() {
		line := scanner.Bytes()

		if mode == bucketsPartFile && string(line) == endbuckets {
			mode = filesPartFile
			continue
		}

		if mode == infoPartFile {
			// first line
			data, err := base64.StdEncoding.DecodeString(string(line))
			if err != nil {
				m.l.Println("[ERR] base64 decode first line data,", err)
				continue
			}
			info = ArchiveInfoLine(string(data))

			m.l.Println("[INFO] start import from:", info)

			if info.AppName() != appname {
				m.l.Println("[ERR] not match appname:", info.AppName(), appname)
				break
			}

			// TODO: more checks

			mode = bucketsPartFile

			continue
		}

		switch mode {
		case bucketsPartFile:
			data, err := base64.StdEncoding.DecodeString(string(line))
			if err != nil {
				m.l.Println("[ERR] base64 decode bucket data,", err)
				continue
			}
			bucket := NewBucket()
			if err := bucket.UnmarshalMsgpack(data); err != nil {
				m.l.Println("[ERR] unmarshal bucket,", err)
				continue
			}

			_existsBucket, err := m.bm.FindBucket(
				bucket.BucketID,
				PrimaryIDsData|PrimaryNamesData,
			)

			m.l.Println("[DEBUG] find bucket", bucket.BucketID, bucket.BucketName, err)

			if err != nil && err != ErrNotFound {
				m.l.Println(
					"[ERR] find bucket by ID",
					bucket.BucketID,
					":",
					err)
				continue
			}

			if uuid.Equal(_existsBucket.BucketID, bucket.BucketID) &&
				!uuid.Equal(uuid.Nil, bucket.BucketID) {

				m.l.Println(
					"[INFO] bucket",
					bucket.BucketID,
					"exists, upgrade bucket...")

				if err := m.bm.UpdateBucket(
					bucket,
					FullBucket,
				); err != nil {
					m.l.Println(
						"[ERR] updated bucket",
						bucket.BucketID,
						err)
				}
				continue
			}

			m.l.Println(
				"[INFO] bucket",
				bucket.BucketID,
				"not exists, create bucket...")

			if err := m.bm.CreateBucket(bucket); err != nil {
				m.l.Println(
					"[ERR] created bucket",
					bucket.BucketID,
					err)
				continue
			}
		case filesPartFile:
			data, err := base64.StdEncoding.DecodeString(string(line))
			if err != nil {
				m.l.Println("[ERR] base64 decode file data,", err)
				continue
			}
			file := NewFile()
			if err := file.UnmarshalMsgpack(data); err != nil {
				m.l.Println("[ERR] unmarshal file,", err)
				continue
			}

			_existsFile, err := m.fm.FindFile(
				file.FileID,
				PrimaryIDsData|PrimaryNamesData,
			)

			m.l.Println("[DEBUG] find file", file.FileID, file.FileName, err)

			if err != nil && err != ErrNotFound {
				m.l.Println(
					"[ERR] find file by ID",
					file.FileID,
					":",
					err)
				continue
			}

			if uuid.Equal(_existsFile.FileID, file.FileID) &&
				!uuid.Equal(uuid.Nil, file.FileID) {

				m.l.Println(
					"[INFO] file",
					file.FileID,
					"exists, upgrade file...")

				if err := m.fm.UpdateFileFrom(
					file,
					FullFile,
				); err != nil {
					m.l.Println(
						"[ERR] updated file",
						file.FileID,
						err)
				}
				continue
			}

			m.l.Println(
				"[INFO] file",
				file.FileID,
				"not exists, create file...")

			if err := m.fm.CreateFile(file); err != nil {
				m.l.Println(
					"[ERR] created file",
					file.FileID,
					err)
				continue
			}
		default:
			m.l.Println("unexpected mode,", mode)
		}
	}

	return info, scanner.Err()
}

func (m *ImportManager) Export(
	verstion, author, description string,
) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	writer := bufio.NewWriter(buf)
	info := NewArchiveInfoLine(
		verstion,
		author,
		description,
		time.Now(),
	)

	// first line
	fmt.Fprintln(writer, base64.StdEncoding.EncodeToString([]byte(info)))

	// buckets

	err := m.bm.(BucketImportManager).EachBucket(func(bucket *Bucket) error {
		data, err := bucket.MarshalMsgpack()
		if err != nil {
			m.l.Println("[ERR] marshal bucket ID", bucket.BucketID, err)
			return errors.New("error marshaler bucket")
		}

		m.l.Println("[DEBUG] export bucket", bucket.BucketID, bucket.BucketName, err)

		_, err = fmt.Fprintln(writer, base64.StdEncoding.EncodeToString(data))
		return err
	})

	if err != nil {
		m.l.Println("[ERR] error iterating buckets", err)
		return []byte{}, err
	}

	// end buckets - separator
	fmt.Fprintln(writer, endbuckets)

	// files

	err = m.fm.(FileImportManager).EachFile(func(file *File) error {
		data, err := file.MarshalMsgpack()
		if err != nil {
			m.l.Println("[ERR] marshal file ID", file.FileID, err)
			return errors.New("error marshaler file")
		}

		m.l.Println("[DEBUG] export file", file.FileID, file.FileName, err)

		_, err = fmt.Fprintln(writer, base64.StdEncoding.EncodeToString(data))
		return err
	})

	if err != nil {
		m.l.Println("[ERR] error iterating files", err)
		return []byte{}, err
	}

	err = writer.Flush()

	return buf.Bytes(), err
}

//

func NewArchiveInfoLine(
	version, author, description string,
	datetime time.Time,
) ArchiveInfoLine {
	args := []string{
		appname,
		version,
		author,
		datetime.Format(time.RFC3339),
		description,
	}

	return ArchiveInfoLine(strings.Join(args, archiveSeparator))
}

type ArchiveInfoLine string

func (i ArchiveInfoLine) valid() bool {
	args := strings.Split(string(i), archiveSeparator)

	if len(args) != 5 {
		/*
		   1. appname,
		   2. version,
		   3. author,
		   4. datetime,
		   5. description,
		*/
		log.Println("[ArchiveInfoLine] invalid, unexpected number of arguments", i)
		return false
	}

	return true
}

func (i ArchiveInfoLine) AppName() string {
	if !i.valid() {
		return ""
	}
	return strings.Split(string(i), archiveSeparator)[0]
}

func (i ArchiveInfoLine) Version() string {
	if !i.valid() {
		return ""
	}
	return strings.Split(string(i), archiveSeparator)[1]
}

func (i ArchiveInfoLine) Author() string {
	if !i.valid() {
		return ""
	}
	return strings.Split(string(i), archiveSeparator)[2]
}

// DateTime RFC3339
func (i ArchiveInfoLine) DateTime() string {
	if !i.valid() {
		return ""
	}
	return strings.Split(string(i), archiveSeparator)[3]
}

func (i ArchiveInfoLine) Description() string {
	if !i.valid() {
		return ""
	}
	return strings.Split(string(i), archiveSeparator)[4]
}
