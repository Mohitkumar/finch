package standalone

import (
	"context"
	"encoding/json"

	"github.com/mohitkumar/finch/model"
	"github.com/mohitkumar/finch/persistence"
	"github.com/mohitkumar/finch/storage"
	"github.com/mohitkumar/finch/storage/standalone_storage"
)

var _ persistence.WorkflowDao = new(BadgerWorkflowDaoImpl)

type Config struct {
	DBPath string
}
type BadgerWorkflowDaoImpl struct {
	Config
	storage standalone_storage.StandAloneStorage
}

func NewBadgerWorkflowDaoImpl(conf Config) (*BadgerWorkflowDaoImpl, error) {
	co := &standalone_storage.Config{
		DBPath: conf.DBPath,
	}
	st := standalone_storage.NewStandAloneStorage(co)
	dao := &BadgerWorkflowDaoImpl{
		Config:  conf,
		storage: *st,
	}
	return dao, nil
}

func (dao *BadgerWorkflowDaoImpl) Save(wf model.Workflow) (bool, error) {
	wfBytes, err := json.Marshal(wf)
	if err != nil {
		return false, err
	}
	ctx := context.Background()
	mod := storage.Put{
		Key:   []byte(getWFKey(wf.Name)),
		Cf:    persistence.METADATA_CF,
		Value: wfBytes,
	}
	dao.storage.Write(ctx, []storage.Modify{{Data: mod}})
	return true, nil
}

func (dao *BadgerWorkflowDaoImpl) Delete(name string) (bool, error) {
	ctx := context.Background()
	mod := storage.Delete{
		Key: []byte(getWFKey(name)),
		Cf:  persistence.METADATA_CF,
	}
	dao.storage.Write(ctx, []storage.Modify{{Data: mod}})
	return true, nil
}

func (dao *BadgerWorkflowDaoImpl) Get(name string) (*model.Workflow, error) {
	ctx := context.Background()
	reader, err := dao.storage.Reader(ctx)
	if err != nil {
		return nil, err
	}
	data, err := reader.GetCF(persistence.METADATA_CF, []byte(getWFKey(name)))

	if err != nil {
		return nil, err
	}
	var wf model.Workflow
	err = json.Unmarshal(data, &wf)
	if err != nil {
		return nil, err
	}
	return &wf, nil
}

func getWFKey(key string) string {
	return persistence.WF_PREFIX + key
}
