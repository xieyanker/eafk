package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"

	"k8s.io/api/core/v1"
	"k8s.io/kubectl/pkg/scheme"

	"k8s.io/apimachinery/pkg/runtime/serializer/protobuf"

	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"io/ioutil"
)

func main() {
	var endpoint, keyFile, certFile, caFile string
	flag.StringVar(&endpoint, "endpoint", "https://127.0.0.1:2379", "etcd endpoint.")
	flag.StringVar(&keyFile, "key", "", "TLS client key.")
	flag.StringVar(&certFile, "cert", "", "TLS client certificate.")
	flag.StringVar(&caFile, "cacert", "", "Server TLS CA certificate.")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprint(os.Stderr, "ERROR: you need to specify action: dump or ls [<key>] or get <key> or apply key <key> -f <file> or delete <key>\n")
		os.Exit(1)
	}
	if flag.Arg(0) == "get" && flag.NArg() == 1 {
		fmt.Fprint(os.Stderr, "ERROR: you need to specify <key> for get operation\n")
		os.Exit(1)
	}
	if flag.Arg(0) == "dump" && flag.NArg() != 1 {
		fmt.Fprint(os.Stderr, "ERROR: you cannot specify positional arguments with dump\n")
		os.Exit(1)
	}
	if flag.Arg(0) == "change-monitors-list" {
		if flag.Arg(1) == "" || flag.Arg(2) == "" {
			fmt.Fprint(os.Stderr, "ERROR: you have to specify both: PV name and list of comma-separated monitor IP-addresses\n")
			os.Exit(1)
		}
		if !regexp.MustCompile(`^pvc-[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`).MatchString(flag.Arg(1)) {
			fmt.Fprint(os.Stderr, "ERROR: invalid PV name. Ex: pvc-dd3afe18-1bae-411c-9a1d-df129847cb62")
			os.Exit(1)
		}
		if !regexp.MustCompile(`^(((\d+\.){3}\d+):\d+,?)+$`).MatchString(flag.Arg(2)) {
			fmt.Fprint(os.Stderr, "ERROR: invalid IP-address list. Ex: 1.1.1.1:6789,2.2.2.2:6789,3.3.3.3:6789\n")
			os.Exit(1)
		}
	}

	//TODO
	if flag.Arg(0) == "apply" {
		if flag.Arg(1) != "key" || flag.Arg(2) == "" || flag.Arg(3) != "-f" || flag.Arg(4) == "" {
			fmt.Fprint(os.Stderr, "ERROR: you need to specify <apply key xxx -f file> for apply operation\n")
			os.Exit(1)
		}
	}
	if flag.Arg(0) == "delete" && flag.NArg() == 1 {
		fmt.Fprint(os.Stderr, "ERROR: you need to specify <key> for delete operation\n")
		os.Exit(1)
	}

	action := flag.Arg(0)
	key := ""
	if flag.NArg() > 1 {
		key = flag.Arg(1)
	}

	var tlsConfig *tls.Config
	if len(certFile) != 0 || len(keyFile) != 0 || len(caFile) != 0 {
		tlsInfo := transport.TLSInfo{
			CertFile:      certFile,
			KeyFile:       keyFile,
			TrustedCAFile: caFile,
		}
		var err error
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: unable to create client config: %v\n", err)
			os.Exit(1)
		}
	}

	config := clientv3.Config{
		Endpoints:   []string{endpoint},
		TLS:         tlsConfig,
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: unable to connect to etcd: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	switch action {
	case "ls":
		_, err = listKeys(client, key)
	case "get":
		err = getKey(client, key)
	case "change-monitors-list":
		err = changeMonitorsList(client, flag.Arg(1), flag.Arg(2))
	case "dump":
		err = dump(client)
	case "delete":
		err = deleteKey(client, key)
	case "apply":
		err = applyFile(client, flag.Arg(2), flag.Arg(4))
	default:
		fmt.Fprintf(os.Stderr, "ERROR: invalid action: %s\n", action)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s-ing %s: %v\n", action, key, err)
		os.Exit(1)
	}
}

func changeMonitorsList(client *clientv3.Client, pvName, list string) error {
	decoder := scheme.Codecs.UniversalDeserializer()

	pvKey := fmt.Sprintf("/registry/persistentvolumes/%s", pvName)

	resp, err := clientv3.NewKV(client).Get(context.Background(), pvKey)
	if err != nil {
		fmt.Printf("get key %s %s\n", pvKey, err)
	}

	obj, _, _ := decoder.Decode(resp.Kvs[0].Value, nil, nil)

	pv := obj.(*v1.PersistentVolume)

	monitors := strings.Split(strings.Trim(list, ","), ",")
	pv.Spec.RBD.CephMonitors = monitors

	protoSerializer := protobuf.NewSerializer(scheme.Scheme, scheme.Scheme)
	newObj := new(bytes.Buffer)
	protoSerializer.Encode(obj, newObj)

	_, err = clientv3.NewKV(client).Put(context.Background(), pvKey, newObj.String())
	if err != nil {
		fmt.Printf("put to key %s %s\n", pvKey, err)
	}

	return nil
}

func listKeys(client *clientv3.Client, key string) ([]string, error) {
	var resp *clientv3.GetResponse
	var err error
	if len(key) == 0 {
		resp, err = clientv3.NewKV(client).Get(context.Background(), "/", clientv3.WithFromKey(), clientv3.WithKeysOnly())
	} else {
		resp, err = clientv3.NewKV(client).Get(context.Background(), key, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	}
	if err != nil {
		return []string{""}, err
	}

	keys := []string{}

	for _, kv := range resp.Kvs {
		fmt.Println(string(kv.Key))
		keys = append(keys, string(kv.Key))
	}

	return keys, err
}

func getKey(client *clientv3.Client, key string) error {
	resp, err := clientv3.NewKV(client).Get(context.Background(), key)
	if err != nil {
		return err
	}

	decoder := scheme.Codecs.UniversalDeserializer()
	encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)

	for _, kv := range resp.Kvs {
		obj, gvk, err := decoder.Decode(kv.Value, nil, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: unable to decode %s: %v\n", kv.Key, err)
			continue
		}
		fmt.Println(gvk)
		err = encoder.Encode(obj, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: unable to encode %s: %v\n", kv.Key, err)
			continue
		}
	}

	return nil
}

func dump(client *clientv3.Client) error {
	response, err := clientv3.NewKV(client).Get(context.Background(), "/", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return err
	}

	kvData := []etcd3kv{}
	decoder := scheme.Codecs.UniversalDeserializer()
	encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, false)
	objJSON := &bytes.Buffer{}

	for _, kv := range response.Kvs {
		obj, _, err := decoder.Decode(kv.Value, nil, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: error decoding value %q: %v\n", string(kv.Value), err)
			continue
		}
		objJSON.Reset()
		if err := encoder.Encode(obj, objJSON); err != nil {
			fmt.Fprintf(os.Stderr, "WARN: error encoding object %#v as JSON: %v", obj, err)
			continue
		}
		kvData = append(
			kvData,
			etcd3kv{
				Key:            string(kv.Key),
				Value:          string(objJSON.Bytes()),
				CreateRevision: kv.CreateRevision,
				ModRevision:    kv.ModRevision,
				Version:        kv.Version,
				Lease:          kv.Lease,
			},
		)
	}

	jsonData, err := json.MarshalIndent(kvData, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))

	return nil
}

func deleteKey(client *clientv3.Client, key string) error {
	_, err := clientv3.NewKV(client).Delete(context.Background(), key)
	if err != nil {
		return err
	}

	fmt.Printf("The key %s has been deleted\n", key)
	return nil
}

func applyFile(client *clientv3.Client, key, path string) error {
	contentByte, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = clientv3.NewKV(client).Put(context.Background(), key, string(contentByte))
	if err != nil {
		return err
	}

	fmt.Printf("The key %s has been putted\n", key)
	return nil
}

type etcd3kv struct {
	Key            string `json:"key,omitempty"`
	Value          string `json:"value,omitempty"`
	CreateRevision int64  `json:"create_revision,omitempty"`
	ModRevision    int64  `json:"mod_revision,omitempty"`
	Version        int64  `json:"version,omitempty"`
	Lease          int64  `json:"lease,omitempty"`
}
