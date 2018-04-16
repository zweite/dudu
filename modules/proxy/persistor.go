package proxy

import (
	"bufio"
	"bytes"
	"dudu/commons/mymgo"
	"dudu/commons/util"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

// 预写持久器
type Persistor interface {
	Write(endpoint, hostname, metric, value, errMsg string, timestamp int64) (int, error)
	Flush() error
	Close() error
}

// 信息存储
type InfoPersistor interface {
	Save(endpoint, hostname, metric string, value interface{}, version int64) error
}

type FilePersistor struct {
	path string
	fi   *os.File
	buf  *bufio.Writer
}

func NewFilePersistor(filePath string) (*FilePersistor, error) {
	dir := filepath.Dir(filePath)
	if err := util.EnsureDir(dir, 0755); err != nil {
		return nil, err
	}

	fi, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FilePersistor{
		path: filePath,
		fi:   fi,
		buf:  bufio.NewWriter(fi),
	}, nil
}

func (f *FilePersistor) Write(endpoint, hostname, metric, value, errMsg string, timestamp int64) (i int, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, len(endpoint)+len(hostname)+len(metric)+len(value)+len(errMsg)+20))
	if _, err = buf.WriteString(endpoint); err != nil {
		return
	}

	if err = buf.WriteByte('\t'); err != nil {
		return
	}

	if _, err = buf.WriteString(hostname); err != nil {
		return
	}

	if err = buf.WriteByte('\t'); err != nil {
		return
	}

	if _, err = buf.WriteString(metric); err != nil {
		return
	}

	if err = buf.WriteByte('\t'); err != nil {
		return
	}

	if _, err = buf.WriteString(value); err != nil {
		return
	}

	if err = buf.WriteByte('\t'); err != nil {
		return
	}

	if _, err = buf.WriteString(errMsg); err != nil {
		return
	}

	if err = buf.WriteByte('\t'); err != nil {
		return
	}

	if _, err = buf.WriteString(strconv.FormatInt(timestamp, 10)); err != nil {
		return
	}

	if err = buf.WriteByte('\n'); err != nil {
		return
	}
	return f.buf.Write(buf.Bytes())
}

func (f *FilePersistor) Flush() error {
	return f.buf.Flush()
}

func (f *FilePersistor) Close() error {
	f.buf.Flush()
	return f.fi.Close()
}

type MongoPersistor struct {
	sess *mymgo.MdbSession
}

func NewMongoPersistor(cfg string) (*MongoPersistor, error) {
	sess, err := mymgo.Open(cfg)
	if err != nil {
		return nil, err
	}

	return &MongoPersistor{
		sess: sess,
	}, nil
}

func (m *MongoPersistor) Save(endPoint, hostName, metric string, value interface{}, version int64) error {
	return m.sess.Insert(metric, struct {
		Id       string      `bson:"_id"`
		EndPoint string      `bson:"endPoint"`
		HostName string      `bson:"hostName"`
		Metric   string      `bson:"metric"`
		Value    interface{} `bson:"value"`
		Version  int64       `bson:"version"`
	}{
		Id:       bson.NewObjectId().Hex(),
		EndPoint: endPoint,
		HostName: hostName,
		Metric:   metric,
		Value:    value,
		Version:  version,
	})
}
