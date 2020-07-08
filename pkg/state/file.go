package state

// Crappy file-based implementation of the verb.StateManager interface.

import(
	"bytes"
	"context"
	"encoding/gob"
	"io/ioutil"
)

type FileStateManager struct {
	Dir string
}

func (fsm FileStateManager)ReadState(ctx context.Context, name string, ptr interface{}) error {
	if data,err := ioutil.ReadFile(fsm.Dir+"/"+name); err != nil {
		return err
	} else {
		buf := bytes.NewBuffer(data)
		if err := gob.NewDecoder(buf).Decode(ptr); err != nil {
			return err
		}
	}
	return nil
}

func (fsm FileStateManager)WriteState(ctx context.Context, name string, ptr interface{}) error {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(ptr); err != nil {
		return err
	} else if err := ioutil.WriteFile(fsm.Dir+"/"+name, buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}
