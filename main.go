package main

import (
	"bytes"
	"fmt"
	goruntime "runtime"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

var flag string

func main() {
	const cnt = 10000
	podList := corev1.PodList{
		Items: make([]corev1.Pod, cnt),
	}

	var list []runtime.Object
	// 复现问题
	if flag == "leak" {
		list, _ = meta.ExtractList(&podList)
	} else {
		// 问题修复
		list, _ = meta.ExtractListWithAlloc(&podList)
	}
	fmt.Printf("Type: %T\n", list[0])

	str := string(bytes.Repeat([]byte("a"), 20_000))
	for i := 0; i < cnt; i++ {
		podList.Items[i].Annotations = map[string]string{"foo": strings.Clone(str)}
	}
	var memStats goruntime.MemStats
	goruntime.ReadMemStats(&memStats)
	fmt.Println("podList alive:", memStats.HeapInuse)

	pod := &list[100]
	podList = corev1.PodList{}
	goruntime.GC()
	goruntime.ReadMemStats(&memStats)
	fmt.Println("podList dead, pod alive:", memStats.HeapInuse)
	goruntime.KeepAlive(pod)

	pod = nil
	goruntime.GC()
	goruntime.ReadMemStats(&memStats)
	fmt.Println("pod dead:", memStats.HeapInuse)
}
