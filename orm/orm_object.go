// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/zhaocloud/beego for the canonical source repository
// @license     http://github.com/zhaocloud/beego/blob/master/LICENSE
// @authors     slene

package orm

import (
	"fmt"
	"reflect"
)

// an insert queryer struct
type insertSet struct {
	mi     *modelInfo
	orm    *orm
	stmt   stmtQuerier
	closed bool
}

var _ Inserter = new(insertSet)

// insert model ignore it's registered or not.
func (o *insertSet) Insert(md interface{}) (int64, error) {
	if o.closed {
		return 0, ErrStmtClosed
	}
	val := reflect.ValueOf(md)
	ind := reflect.Indirect(val)
	typ := ind.Type()
	name := getFullName(typ)
	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<Inserter.Insert> cannot use non-ptr model struct `%s`", name))
	}
	if name != o.mi.fullName {
		panic(fmt.Errorf("<Inserter.Insert> need model `%s` but found `%s`", o.mi.fullName, name))
	}
	id, err := o.orm.alias.DbBaser.InsertStmt(o.stmt, o.mi, ind, o.orm.alias.TZ)
	if err != nil {
		return id, err
	}
	if id > 0 {
		if o.mi.fields.pk.auto {
			if o.mi.fields.pk.fieldType&IsPostiveIntegerField > 0 {
				ind.Field(o.mi.fields.pk.fieldIndex).SetUint(uint64(id))
			} else {
				ind.Field(o.mi.fields.pk.fieldIndex).SetInt(id)
			}
		}
	}
	return id, nil
}

// close insert queryer statement
func (o *insertSet) Close() error {
	if o.closed {
		return ErrStmtClosed
	}
	o.closed = true
	return o.stmt.Close()
}

// create new insert queryer.
func newInsertSet(orm *orm, mi *modelInfo) (Inserter, error) {
	bi := new(insertSet)
	bi.orm = orm
	bi.mi = mi
	st, query, err := orm.alias.DbBaser.PrepareInsert(orm.db, mi)
	if err != nil {
		return nil, err
	}
	if Debug {
		bi.stmt = newStmtQueryLog(orm.alias, st, query)
	} else {
		bi.stmt = st
	}
	return bi, nil
}
