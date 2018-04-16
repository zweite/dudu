package proxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"dudu/commons/mymgo"
	"dudu/commons/util"
	"dudu/config"

	client "github.com/influxdata/influxdb/client/v2"
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
	Close() error
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

func (m *MongoPersistor) Close() error {
	return nil
}

type InfluxPersistor struct {
	httpClient client.Client
	db         string
}

func NewInfluxPersistor(cfg *config.DBConfig) (*InfluxPersistor, error) {
	httpClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port),
		Username: cfg.User,
		Password: cfg.Passwd,
		Timeout:  time.Second * 30, // 默认30s写超时
	})
	if err != nil {
		return nil, err
	}

	return &InfluxPersistor{
		httpClient: httpClient,
		db:         cfg.Db,
	}, nil
}

func (influx *InfluxPersistor) Save(endPoint, hostName, metric string,
	value interface{}, version int64) (err error) {

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influx.db,
		Precision: "ms",
	})

	tags := map[string]string{
		"endPoint": endPoint,
		"hostName": hostName,
		"metric":   metric,
	}

	fields := make(map[string]interface{})

	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		if u, ok := value.(uint64); ok {
			value = int64(u) // 这里是个坑，需要注意!!!
		}
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fields[metric] = value
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		v := reflect.ValueOf(value)
		for i := 0; i < v.Len(); i++ {
			if err := influx.Save(endPoint, hostName, metric, v.Index(i).Interface(), version); err != nil {
				continue
			}
		}
		return

	case reflect.Map:
		fallthrough
	case reflect.Struct:
		fallthrough
	case reflect.Ptr:
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &fields); err != nil {
			return err
		}
	default:
		return fmt.Errorf("can't support data type:%s", t.Kind().String())
	}

	// version is millisecond
	// 1000 = time.Second/time.Millisecond
	timeVersion := time.Unix(version/1000, version%1000*int64(time.Millisecond/time.Nanosecond))
	pt, err := client.NewPoint(metric, tags, fields, timeVersion)
	if err != nil {
		return
	}

	bp.AddPoint(pt)
	return influx.httpClient.Write(bp)
}

func (influx *InfluxPersistor) Close() error {
	return influx.httpClient.Close()
}
