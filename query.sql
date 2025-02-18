-- name: GetMediaByTag :many
SELECT
  mp.file
FROM media_parts mp
  INNER JOIN media_items mi ON mp.media_item_id = mi.id
  INNER JOIN taggings tg ON mi.metadata_item_id = tg.metadata_item_id
  INNER JOIN tags t ON tg.tag_id = t.id
WHERE
  t.tag_type = 314 AND t.tag = ?;
