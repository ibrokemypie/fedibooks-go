package db

import "github.com/hashicorp/go-memdb"

func InitialiseDB() (*memdb.MemDB, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"status": {
				Name: "status",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"author_id": {
						Name:    "author_id",
						Unique:  false,
						Indexer: &SubFieldIndexer{Fields: []Field{{Struct: "Account", Sub: "ID"}}},
					},
					"text": {
						Name:    "text",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Text"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}
