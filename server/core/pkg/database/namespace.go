package database

// GetNamespace will get a namespace from the database by its name
func (db *Database) GetNamespace(name string) (*Namespace, error) {
	namespace := &Namespace{}

	err := db.conn.Model(namespace).
		Where(`"namespace"."name" = ?`, name).
		Relation("ObjectStore").
		Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return namespace, nil
}

// CreateNamespace will create a namespace in the database
func (db *Database) CreateNamespace(name string, description string,
	objectStoreName string, bucket string) error {

	objStore, err := db.GetObjectStore(objectStoreName)
	if err != nil {
		return ConvertError(err)
	}

	ns := &Namespace{Name: name,
		Description:   description,
		ObjectStoreId: objStore.Id,
		BucketName:    bucket}

	_, err = db.conn.Model(ns).Insert()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// DeleteNamespace will delete a namespace from the database by its name
func (db *Database) DeleteNamespace(name string) error {
	namespace := &Namespace{}

	_, err := db.conn.Model(namespace).Where(`"namespace"."name" = ?`, name).Delete()
	if err != nil {
		return ConvertError(err)
	}

	return nil
}

// ListNamespaces will list all namespaces with limit/offset for pagination
func (db *Database) ListNamespaces(limit int, offset int) ([]*Namespace, error) {
	var namespaces []*Namespace

	err := db.conn.Model(&namespaces).Order("id ASC").Limit(limit).Offset(offset).Relation("ObjectStore").Select()
	if err != nil {
		return nil, ConvertError(err)
	}

	return namespaces, nil
}
