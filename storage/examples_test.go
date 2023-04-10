package storage

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	v1 "github.com/vine-io/apimachinery/apis/meta/v1"
	"github.com/vine-io/apimachinery/runtime"
	"github.com/vine-io/apimachinery/schema"
	"github.com/vine-io/apimachinery/storage/dao"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GroupName is the group name for this API
const GroupName = "auth"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

var (
	Builder      = NewFactoryBuilder(addKnownFactory)
	AddToBuilder = Builder.AddToFactory
	sets         = make([]Storage, 0)
)

func addKnownFactory(tx *gorm.DB, f Factory) error {
	return f.AddKnownStorages(tx, SchemeGroupVersion, sets...)
}

var fb Factory

func init() {
	tx := &gorm.DB{}
	fb = NewStorageFactory()

	sets = append(sets, &TestStorage{})

	if err := AddToBuilder(tx, fb); err != nil {
		// handle error
	}
}

func Test_GetStorage(t *testing.T) {
	o := new(TestObj)
	o.GetObjectKind().SetGroupVersionKind(SchemeGroupVersion.WithKind("TestObj"))

	s, err := fb.NewStorage(nil, o)
	if err != nil {
		t.Fatal(err)
	}

	s.Target()
}

type TestObj struct {
	v1.TypeMeta
}

func (t *TestObj) DeepCopyObject() runtime.Object {
	out := new(TestObj)
	*out = *t
	return out
}

func (t *TestObj) DeepFromObject(out runtime.Object) {
	out = t.DeepCopyObject()
}

type TestObjList struct {
	v1.TypeMeta
	Items []*TestObj
}

func (t *TestObjList) DeepCopyObject() runtime.Object {
	out := new(TestObjList)
	*out = *t
	return out
}

func (t *TestObjList) DeepFromObject(out runtime.Object) {
	out = t.DeepCopyObject()
}

// TestStorage the Storage for Test
type TestStorage struct {
	tx    *gorm.DB            `json:"-" dao:"-"`
	exprs []clause.Expression `json:"-" dao:"-"`
}

func (m *TestStorage) FindPk(ctx context.Context, pk any) (runtime.Object, error) {
	//TODO implement me
	panic("implement me")
}

func (m *TestStorage) AutoMigrate(tx *gorm.DB) error {
	return nil
}

func (m *TestStorage) Load(tx *gorm.DB, object runtime.Object) error {
	in := object.(*TestObj)
	_ = in
	return nil
}

func (m *TestStorage) WithTx(tx *gorm.DB) *TestStorage {
	m.tx = tx
	return m
}

func (m *TestStorage) ToTest() *TestObj {
	return &TestObj{}
}

func (TestStorage) TableName() string {
	return "tests"
}

func (m TestStorage) PrimaryKey() (string, interface{}, bool) {
	return "uid", "", true
}

func (m *TestStorage) FindPage(ctx context.Context, page, size int32) (runtime.Object, error) {
	pk, _, _ := m.PrimaryKey()

	m.exprs = append(m.exprs,
		clause.OrderBy{Columns: []clause.OrderByColumn{{Column: clause.Column{Table: m.TableName(), Name: pk}, Desc: true}}},
		dao.Cond().Build("inner_deletion_timestamp", 0),
	)

	//total, err := m.Count(ctx)
	//if err != nil {
	//	return nil,  err
	//}

	m.exprs = append(m.exprs, clause.Limit{Offset: int((page - 1) * size)})
	data, err := m.findAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *TestStorage) FindAll(ctx context.Context) (runtime.Object, error) {
	m.exprs = append(m.exprs, dao.Cond().Build("inner_deletion_timestamp", 0))
	return m.findAll(ctx)
}

func (m *TestStorage) FindPureAll(ctx context.Context) (runtime.Object, error) {
	return m.findAll(ctx)
}

func (m *TestStorage) findAll(ctx context.Context) (runtime.Object, error) {
	dest := make([]*TestObj, 0)
	tx := m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx)

	clauses := append(m.extractClauses(tx), m.exprs...)
	if err := tx.Clauses(clauses...).Find(&dest).Error; err != nil {
		return nil, err
	}

	out := &TestObjList{}
	out.Items = dest

	return out, nil
}

func (m *TestStorage) FindEntitiesPage(ctx context.Context, page, size int) ([]*TestObj, int64, error) {
	pk, _, _ := m.PrimaryKey()

	m.exprs = append(m.exprs,
		clause.OrderBy{Columns: []clause.OrderByColumn{{Column: clause.Column{Table: m.TableName(), Name: pk}, Desc: true}}},
		dao.Cond().Build("inner_deletion_timestamp", 0),
	)

	total, err := m.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	m.exprs = append(m.exprs, clause.Limit{Offset: int((page - 1) * size)})
	data, err := m.findAllEntities(ctx)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (m *TestStorage) FindAllEntities(ctx context.Context) ([]*TestObj, error) {
	m.exprs = append(m.exprs, dao.Cond().Build("inner_deletion_timestamp", 0))
	return m.findAllEntities(ctx)
}

func (m *TestStorage) FindPureAllEntities(ctx context.Context) ([]*TestObj, error) {
	return m.findAllEntities(ctx)
}

func (m *TestStorage) findAllEntities(ctx context.Context) ([]*TestObj, error) {
	dest := make([]*TestStorage, 0)
	tx := m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx)

	clauses := append(m.extractClauses(tx), m.exprs...)
	if err := tx.Clauses(clauses...).Find(&dest).Error; err != nil {
		return nil, err
	}

	outs := make([]*TestObj, len(dest))
	for i := range dest {
		outs[i] = dest[i].ToTest()
	}

	return outs, nil
}

func (m *TestStorage) Count(ctx context.Context) (total int64, err error) {
	tx := m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx)

	clauses := append(m.extractClauses(tx), dao.Cond().Build("inner_deletion_timestamp", 0))
	clauses = append(clauses, m.exprs...)

	err = tx.Clauses(clauses...).Count(&total).Error
	return
}

func (m *TestStorage) FindOne(ctx context.Context) (runtime.Object, error) {
	tx := m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx)

	clauses := append(m.extractClauses(tx), dao.Cond().Build("inner_deletion_timestamp", 0))
	clauses = append(clauses, m.exprs...)

	if err := tx.Clauses(clauses...).First(&m).Error; err != nil {
		return nil, err
	}

	return m.ToTest(), nil
}

func (m *TestStorage) FindPureOne(ctx context.Context) (runtime.Object, error) {
	tx := m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx)

	if err := tx.Clauses(m.exprs...).First(&m).Error; err != nil {
		return nil, err
	}

	return m.ToTest(), nil
}

func (m *TestStorage) Cond(exprs ...clause.Expression) Storage {
	m.exprs = append(m.exprs, exprs...)
	return m
}

func (m *TestStorage) Target() reflect.Type {
	return reflect.TypeOf(new(TestObj))
}

func (m *TestStorage) extractClauses(tx *gorm.DB) []clause.Expression {
	exprs := make([]clause.Expression, 0)
	return exprs
}

func (m *TestStorage) Create(ctx context.Context) (runtime.Object, error) {
	tx := m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx)

	if err := tx.Create(m).Error; err != nil {
		return nil, err
	}

	return m.ToTest(), nil
}

func (m *TestStorage) BatchUpdates(ctx context.Context) error {
	tx := m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx)
	values := map[string]interface{}{}
	return tx.Clauses(m.exprs...).Updates(values).Error
}

func (m *TestStorage) Updates(ctx context.Context) (runtime.Object, error) {
	pk, pkv, isNil := m.PrimaryKey()
	if isNil {
		return nil, errors.New("missing primary key")
	}

	values := map[string]interface{}{}
	tx := m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx).Where(pk+" = ?", pkv)

	if err := tx.Updates(values).Error; err != nil {
		return nil, err
	}

	err := m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx).Where(pk+" = ?", pkv).First(m).Error
	if err != nil {
		return nil, err
	}
	return m.ToTest(), nil
}

func (m *TestStorage) BatchDelete(ctx context.Context, soft bool) error {
	tx := m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx)
	clauses := append(m.exprs, m.extractClauses(tx)...)

	if soft {
		return tx.Clauses(clauses...).Updates(map[string]interface{}{"inner_deletion_timestamp": time.Now().UnixNano()}).Error
	}
	return tx.Clauses(clauses...).Delete(&TestStorage{}).Error
}

func (m *TestStorage) Delete(ctx context.Context, soft bool) error {
	pk, pkv, isNil := m.PrimaryKey()
	if isNil {
		return errors.New("missing primary key")
	}

	tx := m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx)

	if soft {
		return tx.Where(pk+" = ?", pkv).Updates(map[string]interface{}{"inner_deletion_timestamp": time.Now().UnixNano()}).Error
	}
	return tx.Where(pk+" = ?", pkv).Delete(&TestStorage{}).Error
}

func (m *TestStorage) Tx(ctx context.Context) *gorm.DB {
	return m.tx.Session(&gorm.Session{}).Table(m.TableName()).WithContext(ctx).Clauses(m.exprs...)
}
