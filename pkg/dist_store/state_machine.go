package dist_store

import (
	"github.com/fbundle/paxos/pkg/local_store"
	"github.com/fbundle/paxos/pkg/paxos"
	"github.com/google/uuid"
)

type Entry struct {
	Key string `json:"key"`
	Val string `json:"val"`
	Ver uint64 `json:"ver"`
}

type Cmd struct {
	Uuid    uuid.UUID `json:"uuid"`
	Entries []Entry   `json:"entries"`
}

func makeCmd(entries []Entry) Cmd {
	return Cmd{
		Uuid:    uuid.New(),
		Entries: entries,
	}
}

func (cmd Cmd) Equal(other Cmd) bool {
	return cmd.Uuid == other.Uuid
}

type stateMachine struct {
	store local_store.MemStore[string, Entry]
}

func newStateMachine() *stateMachine {
	return &stateMachine{
		store: local_store.NewMemStore[string, Entry](),
	}
}

func (sm *stateMachine) Get(key string) Entry {
	return sm.store.Update(func(txn local_store.Txn[string, Entry]) any {
		return getDefaultEntry(txn, key)
	}).(Entry)
}

func (sm *stateMachine) Keys() []string {
	return sm.store.Update(func(txn local_store.Txn[string, Entry]) any {
		return sm.store.Keys()
	}).([]string)
}

func (sm *stateMachine) Apply(logId paxos.LogId, cmd Cmd) {
	sm.store.Update(func(txn local_store.Txn[string, Entry]) any {
		for _, entry := range cmd.Entries {
			oldEntry := getDefaultEntry(txn, entry.Key)
			if entry.Ver <= oldEntry.Ver {
				continue // ignore update
			}
			if len(entry.Val) == 0 {
				txn.Del(entry.Key)
			} else {
				txn.Set(entry.Key, entry)
			}
		}
		return nil
	})
}
