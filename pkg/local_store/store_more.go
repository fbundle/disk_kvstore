package local_store

import "encoding/json"

func zero[T any]() T {
	var v T
	return v
}

type storeKV[K comparable, V any] struct {
	ss StringStore
}

func (s *storeKV[K, V]) Update(update func(txn Txn[K, V]) any) any {
	var out any
	s.ss.Update(func(txn Txn[string, string]) any {
		out = update(&txnKV[K, V]{txn: txn})
		return nil
	})
	return out
}

type txnKV[K comparable, V any] struct {
	txn Txn[string, string]
}

func (t *txnKV[K, V]) Get(k K) (v V, ok bool) {
	kb, err := json.Marshal(k)
	if err != nil {
		panic(err)
	}
	ks := string(kb)
	vs, ok := t.txn.Get(ks)
	if !ok {
		return zero[V](), false
	}
	err = json.Unmarshal([]byte(vs), &v)
	if err != nil {
		panic(err)
	}
	return v, true
}

func (t *txnKV[K, V]) Set(k K, v V) {
	kb, err := json.Marshal(k)
	if err != nil {
		panic(err)
	}
	ks := string(kb)
	vb, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	vs := string(vb)
	t.txn.Set(ks, vs)
}

func (t *txnKV[K, V]) Del(k K) {
	kb, err := json.Marshal(k)
	if err != nil {
		panic(err)
	}
	ks := string(kb)
	t.txn.Del(ks)
}
