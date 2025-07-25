package crud

import (
	"encoding/json"
	"github.com/Lumicrate/gompose/db"
	"github.com/Lumicrate/gompose/hooks"
	"github.com/Lumicrate/gompose/http"
	"reflect"
	"strconv"
	"strings"
)

func handleGetAll(ctx http.Context, dbAdapter db.DBAdapter, entity any) {
	filters := map[string]any{}
	pagination := db.Pagination{Limit: 10, Offset: 0} // default pagination
	sort := []db.Sort{}

	// parse filters, pagination and sort from query params
	for key, vals := range ctx.QueryParams() {
		val := vals[0]
		switch key {
		case "limit":
			if l, err := strconv.Atoi(val); err == nil {
				pagination.Limit = l
			}
		case "offset":
			if o, err := strconv.Atoi(val); err == nil {
				pagination.Offset = o
			}
		case "sort":
			// example: sort=name,-created_at
			fields := strings.Split(val, ",")
			for _, f := range fields {
				direction := "asc"
				if strings.HasPrefix(f, "-") {
					direction = "desc"
					f = strings.TrimPrefix(f, "-")
				}
				sort = append(sort, db.Sort{Field: f, Direction: direction})
			}
		default:
			filters[key] = val
		}
	}

	result, err := dbAdapter.FindAll(entity, filters, pagination, sort)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, result)
}

func handleGetByID(ctx http.Context, dbAdapter db.DBAdapter, entity any) {
	id := ctx.Param("id")

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	newEntity := reflect.New(t).Interface()

	found, err := dbAdapter.FindByID(id, newEntity)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "entity not found"})
		return
	}

	ctx.JSON(200, found)
}

func handleCreate(ctx http.Context, dbAdapter db.DBAdapter, entity any) {
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	newEntity := reflect.New(t).Interface()

	if err := ctx.Bind(newEntity); err != nil {
		ctx.JSON(400, map[string]string{"error": "invalid input" + err.Error()})
		return
	}

	if hook, ok := newEntity.(hooks.BeforeCreate); ok {
		if err := hook.BeforeCreate(); err != nil {
			ctx.JSON(400, map[string]string{"error": "beforeSave failed: " + err.Error()})
			return
		}
	}

	if err := dbAdapter.Create(newEntity); err != nil {
		ctx.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if hook, ok := newEntity.(hooks.AfterCreate); ok {
		if err := hook.AfterCreate(); err != nil {
			ctx.JSON(400, map[string]string{"error": "afterSave failed: " + err.Error()})
		}
	}

	ctx.JSON(201, newEntity)
}

func handleUpdate(ctx http.Context, dbAdapter db.DBAdapter, entity any) {
	id := ctx.Param("id")
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	updatedEntity := reflect.New(t).Interface()

	if err := ctx.Bind(updatedEntity); err != nil {
		ctx.JSON(400, map[string]string{"error": "invalid input"})
		return
	}

	// Set the ID field in the updated entity to the URL param id if field exists
	setEntityID(updatedEntity, id)

	if hook, ok := updatedEntity.(hooks.BeforeUpdate); ok {
		if err := hook.BeforeUpdate(); err != nil {
			ctx.JSON(400, map[string]string{"error": "beforeUpdate failed: " + err.Error()})
			return
		}
	}

	if err := dbAdapter.Update(updatedEntity); err != nil {
		ctx.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if hook, ok := updatedEntity.(hooks.AfterUpdate); ok {
		if err := hook.AfterUpdate(); err != nil {
			ctx.JSON(400, map[string]string{"error": "afterUpdate failed: " + err.Error()})
			return
		}
	}

	ctx.JSON(200, updatedEntity)
}

func handlePatch(ctx http.Context, dbAdapter db.DBAdapter, entity any) {
	id := ctx.Param("id")
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	existingEntity := reflect.New(t).Interface()

	found, err := dbAdapter.FindByID(id, existingEntity)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "entity not found"})
		return
	}

	patchData := map[string]interface{}{}
	if err := ctx.BindJSON(&patchData); err != nil {
		ctx.JSON(400, map[string]string{"error": "invalid patch data"})
		return
	}

	patchBytes, _ := json.Marshal(patchData)
	if err := json.Unmarshal(patchBytes, &found); err != nil {
		ctx.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if hook, ok := existingEntity.(hooks.BeforePatch); ok {
		if err := hook.BeforePatch(); err != nil {
			ctx.JSON(400, map[string]string{"error": "beforePatch failed: " + err.Error()})
			return
		}
	}

	if err := dbAdapter.Update(found); err != nil {
		ctx.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if hook, ok := existingEntity.(hooks.AfterPatch); ok {
		if err := hook.AfterPatch(); err != nil {
			ctx.JSON(400, map[string]string{"error": "afterPatch failed: " + err.Error()})
			return
		}
	}

	ctx.JSON(200, found)
}

func handleDelete(ctx http.Context, dbAdapter db.DBAdapter, entity any) {
	id := ctx.Param("id")

	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	toDeleteEntity := reflect.New(t).Interface()

	if hook, ok := toDeleteEntity.(hooks.BeforeDelete); ok {
		if err := hook.BeforeDelete(); err != nil {
			ctx.JSON(400, map[string]string{"error": "beforeDelete failed: " + err.Error()})
			return
		}
	}

	if err := dbAdapter.Delete(id, toDeleteEntity); err != nil {
		ctx.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	if hook, ok := toDeleteEntity.(hooks.AfterDelete); ok {
		if err := hook.AfterDelete(); err != nil {
			ctx.JSON(400, map[string]string{"error": "afterDelete failed: " + err.Error()})
			return
		}
	}

	ctx.JSON(204, nil)
}

func setEntityID(entity any, id string) {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByName("ID")
	if !field.IsValid() || !field.CanSet() {
		return
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(id)
	case reflect.Int, reflect.Int64, reflect.Int32:
		if intVal, err := strconv.ParseInt(id, 10, 64); err == nil {
			field.SetInt(intVal)
		}
	default:
		panic("ID type not supported")
	}
}
