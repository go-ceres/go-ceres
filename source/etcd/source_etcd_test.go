package etcd

import (
	"log"
	"testing"
	"time"
)

func TestEtcdSource_Read(t *testing.T) {
	etcd := NewSource(
		DefaultConfig().WithEndpoints("203.195.200.40:12379", "203.195.200.40:22379", "203.195.200.40:32379"),
	)
	ds, err := etcd.Read()
	etcd.Watch()
	go func() {
		for range etcd.IsChanged() {
			b, err2 := etcd.Read()
			log.Print(b)
			log.Print("监测到修改了", "json是：", string(b.Data), "-------", err2)
		}
		log.Print("结束了")
	}()
	time.Sleep(1 * time.Minute)
	etcd.UnWatch()
	time.Sleep(1 * time.Minute)
	etcd.Watch()
	select {}

	t.Log(string(ds.Data), err)
}
