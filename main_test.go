package main

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"plex-query/database"
)

// Rename a closed fixture to a name with URI delimiters. The prefix also exists,
// so opening the wrong database cannot pass just because Ping succeeded.
func TestPlexDBPathAndActiveMedia(t *testing.T) {
	dir := t.TempDir()
	plain := filepath.Join(dir, "library")
	db, err := sql.Open("sqlite", plain)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(string(schema)); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`
		INSERT INTO tags(id,tag_type,tag) VALUES(1,314,'imdb://tt123');
		INSERT INTO taggings(metadata_item_id,tag_id) VALUES(10,1);
		INSERT INTO media_items(id,metadata_item_id,deleted_at) VALUES(1,10,NULL),(2,10,123);
		INSERT INTO media_parts(media_item_id,file,deleted_at) VALUES
			(1,'/active.mkv',NULL),(1,'/deleted-part.mkv',123),(2,'/deleted-item.mkv',NULL);
	`)
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
	name := "library#1 ?100%.db"
	if err = os.Rename(plain, filepath.Join(dir, name)); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(plain, nil, 0600); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	ro, err := openPlexDB(name)
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	files, err := database.New(ro).GetMediaByTag(context.Background(), sql.NullString{String: "imdb://tt123", Valid: true})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(files, []sql.NullString{{String: "/active.mkv", Valid: true}}) {
		t.Fatalf("unexpected files: %v", files)
	}
	if _, err = ro.Exec("DELETE FROM tags"); err == nil {
		t.Fatal("read-only database accepted a write")
	}
}
