# EAFK(Etcd Artifact For Kubernetes)

The EAFK(Etcd Artifact For Kubernetes) is a tool for manipulating etcd data.

It refer to two different versions of etcdhelper, one [etcdhelper](https://github.com/openshift/origin/tree/master/tools/etcdhelper) is developed by OpenShift, another [etcdhelper](https://github.com/flant/examples/tree/master/2020/04-etcdhelper) is modified by Flant.

The EAFK supports `apply, get, list, delete` sub command to manipulate the etcd data directly, rather than using kubectl. The advantage it brings is that it can bypass some of the limitations of k8s, which can be useful in some cases.

### How to use

##### 1. eafk ls

`eafk ls [<key>]` can list all `etcd` keys.

```shell
eafk --cacert /etc/ssl/etcd/ssl/ca.pem --cert /etc/ssl/etcd/ssl/node-master1.pem --key /etc/ssl/etcd/ssl/node-master1-key.pem --endpoint https://127.0.0.1:2379 ls

eafk --cacert /etc/ssl/etcd/ssl/ca.pem --cert /etc/ssl/etcd/ssl/node-master1.pem --key /etc/ssl/etcd/ssl/node-master1-key.pem --endpoint https://127.0.0.1:2379 ls /registry/services/
```

result:

```shell
/registry/services/endpoints/cdi/cdi-api
/registry/services/endpoints/cdi/cdi-prometheus-metrics
/registry/services/endpoints/cdi/cdi-uploadproxy
/registry/services/endpoints/cdi/cdi-uploadproxy-nodeport
...
```

##### 2. eafk get

`eafk get <key>` can get the key's gvk and the value with json format.

```shell
eafk --cacert /etc/ssl/etcd/ssl/ca.pem --cert /etc/ssl/etcd/ssl/node-master1.pem --key /etc/ssl/etcd/ssl/node-master1-key.pem --endpoint https://127.0.0.1:2379 get /registry/services/specs/kube-system/coredns
```

result:

```shell
/v1, Kind=Service
{
  "kind": "Service",
  "apiVersion": "v1",
  "metadata": {
    "name": "coredns",
    "namespace": "kube-system",
    "uid": "16b9bccd-faf5-4952-bb16-b9eed607b5ba",
    "creationTimestamp": "2021-10-09T06:37:00Z",
    "labels": {
      "addonmanager.kubernetes.io/mode": "Reconcile",
      "k8s-app": "kube-dns",
      "kubernetes.io/cluster-service": "true",
      "kubernetes.io/name": "coredns"
    }
  },
  "spec": {
    "ports": [
      {
        "name": "dns",
        "protocol": "UDP",
        "port": 53,
        "targetPort": 53
      },
      {
        "name": "dns-tcp",
        "protocol": "TCP",
        "port": 53,
        "targetPort": 53
      },
      {
        "name": "metrics",
        "protocol": "TCP",
        "port": 9153,
        "targetPort": 9153
      }
    ],
    "selector": {
      "k8s-app": "kube-dns"
    },
    "clusterIP": "100.105.0.3",
    "type": "ClusterIP",
    "sessionAffinity": "None"
  },
  "status": {
    "loadBalancer": {

    }
  }
}
```

### 3. eafk apply

`eafk apply` can `create` or `update` the `key`'s value by specific file.

```json
{
  "kind": "Endpoints",
  "apiVersion": "v1",
  "metadata": {
    "name": "karmada-apiserver",
    "namespace": "default",
    "uid": "fcde5443-8e53-4110-b6ba-5004104beee5",
    "creationTimestamp": "2021-10-12T03:18:04Z",
    "labels": {
      "app": "karmada-apiserver"
    }
  }
}
```

```shell
eafk --cacert /etc/ssl/etcd/ssl/ca.pem --cert /etc/ssl/etcd/ssl/node-master1.pem --key /etc/ssl/etcd/ssl/node-master1-key.pem --endpoint https://127.0.0.1:2379 apply key /registry/services/endpoints/default/karmada-apiserver -f test.json
```

result:

```shell
The key /registry/services/endpoints/default/karmada-apiserver has been putted
```

### 4. eafk delete

`eafk delete` can delete a `key` and its `value` from the etcd.

```shell
eafk --cacert /etc/ssl/etcd/ssl/ca.pem --cert /etc/ssl/etcd/ssl/node-master1.pem --key /etc/ssl/etcd/ssl/node-master1-key.pem --endpoint https://127.0.0.1:2379 delete /registry/services/endpoints/default/karmada-apiserver
```

result:

```shell
The key /registry/services/endpoints/default/karmada-apiserver has been deleted
```

##### 5. eafk dump

`eafk dump` can get all key and value, and print them all by json format.

```shell
eafk --cacert /etc/ssl/etcd/ssl/ca.pem --cert /etc/ssl/etcd/ssl/node-master1.pem --key /etc/ssl/etcd/ssl/node-master1-key.pem --endpoint https://127.0.0.1:2379 dump
```

result:

```shell
...
  {
    "key": "/registry/pods/kube-system/node-problem-detector-s2cj6",
    "value": "{\"kind\":\"Pod\",\"apiVersion\":\"v1\",\"metadata\":{\"name\":\"node-problem-detector-s2cj6\",\"generateName\":\"node-problem-detector-\",\"namespace\":\"kube-system\",\"uid\":\"b6c15277-8adf-47b5-bfd7-96fa3a9a416f\",\"creationTimestamp\":\"2021-10-09T06:38:17Z\",\"labels\":{\"app\":\"node-problem-detector\",\"controller-revision-hash\":\"5f7865559\",\"pod-template-generation\":\"1\"},\"ownerReferences\":[{\"apiVersion\":\"apps/v1\",\"kind\":\"DaemonSet\",\"name\":\"node-problem-detector\",\"uid\":\"1e0da6d8-6e12-46fc-a6db-279499ffb7eb\",\"controller\":true,\"blockOwnerDeletion\":true}]},\"spec\":{\"volumes\":[{\"name\":\"log\",\"hostPath\":{\"path\":\"/var/log/\",\"type\":\"\"}},{\"name\":\"kmsg\",\"hostPath\":{\"path\":\"/dev/kmsg\",\"type\":\"\"}},{\"name\":\"localtime\",\"hostPath\":{\"path\":\"/etc/localtime\",\"type\":\"\"}},{\"name\":\"machine-id\",\"hostPath\":{\"path\":\"/etc/machine-id\",\"type\":\"File\"}},{\"name\":\"systemd\",\"hostPath\":{\"path\":\"/run/systemd/system/\",\"type\":\"Directory\"}},{\"name\":\"dbus\",\"hostPath\":{\"path\":\"/var/run/dbus/\",\"type\":\"Directory\"}},{\"name\":\"docker-sock\",\"hostPath\":{\"path\":\"/var/run/docker.sock\",\"type\":\"Socket\"}},{\"name\":\"check-script\",\"configMap\":{\"name\":\"node-problem-detector-config-sh\",\"items\":[{\"key\":\"check_ntp.sh\",\"path\":\"check_ntp.sh\",\"mode\":511},{\"key\":\"network_problem.sh\",\"path\":\"network_problem.sh\",\"mode\":511}],\"defaultMode\":420}},{\"name\":\"config\",\"configMap\":{\"name\":\"node-problem-detector-config\",\"items\":[{\"key\":\"known-modules.json\",\"path\":\"known-modules.json\"},{\"key\":\"abrt-adaptor.json\",\"path\":\"abrt-adaptor.json\"},{\"key\":\"docker-monitor.json\",\"path\":\"docker-monitor.json\"},{\"key\":\"kernel-monitor.json\",\"path\":\"kernel-monitor.json\"},{\"key\":\"kubelet-log-monitor.json\",\"path\":\"kubelet-log-monitor.json\"},{\"key\":\"systemd-monitor.json\",\"path\":\"systemd-monitor.json\"},{\"key\":\"custom-plugin-monitor.json\",\"path\":\"custom-plugin-monitor.json\"},{\"key\":\"docker-monitor-counter.json\",\"path\":\"docker-monitor-counter.json\"},{\"key\":\"health-checker-docker.json\",\"path\":\"health-checker-docker.json\"},{\"key\":\"health-checker-kubelet.json\",\"path\":\"health-checker-kubelet.json\"},{\"key\":\"kernel-monitor-counter.json\",\"path\":\"kernel-monitor-counter.json\"},{\"key\":\"network-problem-monitor.json\",\"path\":\"network-problem-monitor.json\"},{\"key\":\"systemd-monitor-counter.json\",\"path\":\"systemd-monitor-counter.json\"},{\"key\":\"net-cgroup-system-stats-monitor.json\",\"path\":\"net-cgroup-system-stats-monitor.json\"},{\"key\":\"system-stats-monitor.json\",\"path\":\"system-stats-monitor.json\"}],\"defaultMode\":420}},{\"name\":\"node-problem-detector-token-6mnb5\",\"secret\":{\"secretName\":\"node-problem-detector-token-6mnb5\",\"defaultMode\":420}}],\"containers\":[{\"name\":\"node-problem-detector\",\"image\":\"registry-jinan-lab.inspurcloud.cn/library/cke/node-problem-detector:0.8.9-20210819_190210\",\"command\":[\"/node-problem-detector\",\"--logtostderr\",\"--config.system-log-monitor=/config/abrt-adaptor.json,/config/docker-monitor.json,/config/kernel-monitor.json,/config/kubelet-log-monitor.json,/config/systemd-monitor.json\",\"--config.custom-plugin-monitor=/config/custom-plugin-monitor.json,/config/docker-monitor-counter.json,/config/health-checker-docker.json,/config/health-checker-kubelet.json,/config/kernel-monitor-counter.json,/config/network-problem-monitor.json,/config/systemd-monitor-counter.json\",\"--config.system-stats-monitor=/config/net-cgroup-system-stats-monitor.json,/config/system-stats-monitor.json\"],\"env\":[{\"name\":\"NODE_NAME\",\"valueFrom\":{\"fieldRef\":{\"apiVersion\":\"v1\",\"fieldPath\":\"spec.nodeName\"}}}],\"resources\":{\"limits\":{\"cpu\":\"50m\",\"memory\":\"80Mi\"},\"requests\":{\"cpu\":\"50m\",\"memory\":\"80Mi\"}},\"volumeMounts\":[{\"name\":\"log\",\"mountPath\":\"/var/log\"},{\"name\":\"kmsg\",\"readOnly\":true,\"mountPath\":\"/dev/kmsg\"},{\"name\":\"localtime\",\"readOnly\":true,\"mountPath\":\"/etc/localtime\"},{\"name\":\"config\",\"readOnly\":true,\"mountPath\":\"/config\"},{\"name\":\"check-script\",\"readOnly\":true,\"mountPath\":\"/checkscript\"},{\"name\":\"machine-id\",\"readOnly\":true,\"mountPath\":\"/etc/machine-id\"},{\"name\":\"systemd\",\"mountPath\":\"/run/systemd/system\"},{\"name\":\"docker-sock\",\"mountPath\":\"/var/run/docker.sock\"},{\"name\":\"dbus\",\"mountPath\":\"/var/run/dbus/\",\"mountPropagation\":\"Bidirectional\"},{\"name\":\"node-problem-detector-token-6mnb5\",\"readOnly\":true,\"mountPath\":\"/var/run/secrets/kubernetes.io/serviceaccount\"}],\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\",\"imagePullPolicy\":\"IfNotPresent\",\"securityContext\":{\"privileged\":true}}],\"restartPolicy\":\"Always\",\"terminationGracePeriodSeconds\":30,\"dnsPolicy\":\"ClusterFirst\",\"nodeSelector\":{\"node-role.kubernetes.io/node\":\"true\"},\"serviceAccountName\":\"node-problem-detector\",\"serviceAccount\":\"node-problem-detector\",\"nodeName\":\"slave1\",\"hostNetwork\":true,\"securityContext\":{},\"affinity\":{\"nodeAffinity\":{\"requiredDuringSchedulingIgnoredDuringExecution\":{\"nodeSelectorTerms\":[{\"matchFields\":[{\"key\":\"metadata.name\",\"operator\":\"In\",\"values\":[\"slave1\"]}]}]}}},\"schedulerName\":\"default-scheduler\",\"tolerations\":[{\"operator\":\"Exists\",\"effect\":\"NoSchedule\"},{\"operator\":\"Exists\",\"effect\":\"NoExecute\"},{\"key\":\"node-role.kubernetes.io/master\",\"operator\":\"Exists\",\"effect\":\"NoSchedule\"},{\"key\":\"node-role.kubernetes.io/master\",\"operator\":\"Exists\",\"effect\":\"NoExecute\"},{\"key\":\"node.kubernetes.io/not-ready\",\"operator\":\"Exists\",\"effect\":\"NoExecute\"},{\"key\":\"node.kubernetes.io/unreachable\",\"operator\":\"Exists\",\"effect\":\"NoExecute\"},{\"key\":\"node.kubernetes.io/disk-pressure\",\"operator\":\"Exists\",\"effect\":\"NoSchedule\"},{\"key\":\"node.kubernetes.io/memory-pressure\",\"operator\":\"Exists\",\"effect\":\"NoSchedule\"},{\"key\":\"node.kubernetes.io/pid-pressure\",\"operator\":\"Exists\",\"effect\":\"NoSchedule\"},{\"key\":\"node.kubernetes.io/unschedulable\",\"operator\":\"Exists\",\"effect\":\"NoSchedule\"},{\"key\":\"node.kubernetes.io/network-unavailable\",\"operator\":\"Exists\",\"effect\":\"NoSchedule\"}],\"priorityClassName\":\"system-node-critical\",\"priority\":2000001000,\"enableServiceLinks\":true,\"preemptionPolicy\":\"PreemptLowerPriority\"},\"status\":{\"phase\":\"Running\",\"conditions\":[{\"type\":\"Initialized\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"2021-10-09T06:38:16Z\"},{\"type\":\"Ready\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"2021-10-09T06:38:36Z\"},{\"type\":\"ContainersReady\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"2021-10-09T06:38:36Z\"},{\"type\":\"PodScheduled\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"2021-10-09T06:38:17Z\"}],\"hostIP\":\"192.168.122.109\",\"podIP\":\"192.168.122.109\",\"podIPs\":[{\"ip\":\"192.168.122.109\"}],\"startTime\":\"2021-10-09T06:38:16Z\",\"containerStatuses\":[{\"name\":\"node-problem-detector\",\"state\":{\"running\":{\"startedAt\":\"2021-10-09T06:38:35Z\"}},\"lastState\":{},\"ready\":true,\"restartCount\":0,\"image\":\"registry-jinan-lab.inspurcloud.cn/library/cke/node-problem-detector:0.8.9-20210819_190210\",\"imageID\":\"docker-pullable://registry-jinan-lab.inspurcloud.cn/library/cke/node-problem-detector@sha256:c56307dcaacb71261f5cb5fae750dbd2c0b7b8d9f69b1ae0e4e41495f6ffbb97\",\"containerID\":\"docker://c893debda70d19d15b58cce5021d727856fa6381a2e221f02c9d3f049d10595a\",\"started\":true}],\"qosClass\":\"Guaranteed\"}}\n",
    "create_revision": 3590,
    "mod_revision": 4709,
    "version": 4
  },
...
```

### TODO

1. Optimize subcommands and parameters.

2. Add eafk version to print current version info.

3. Add eafk init to generate the etcd's config file, and read the config to reduce the command's etcd parameter.

4. Add eafk edit and eafk patch command.

4. Encapsulate some commonly used functions, such as update the pv's structure.


