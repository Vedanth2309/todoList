package services

import (
	"errors"
	"log"

	"tracker/config"
	"tracker/models"
)

const idxAll = "tracker_tasks,tracker_notes,tracker_diary"

// IndexDoc creates or updates (PUT is an upsert) the Elasticsearch copy of a document.
func IndexDoc(col string, d models.Doc) {
	s := d.Search()
	if s == nil || !config.ESEnabled() {
		return
	}
	b := d.B()
	body := map[string]any{"userId": b.UserID.Hex(), "type": s.Type, "title": s.Title, "content": s.Content, "tags": s.Tags, "date": s.Date}
	if _, code, err := config.ESRequest("PUT", "/tracker_"+col+"/_doc/"+b.ID.Hex(), body); err != nil || code >= 300 {
		log.Printf("elasticsearch index (%s) failed, status %d: %v", col, code, err)
	}
}

func DeleteDoc(col, id string) {
	if _, ok := map[string]bool{"tasks": true, "notes": true, "diary": true}[col]; !ok || !config.ESEnabled() {
		return
	}
	if _, code, err := config.ESRequest("DELETE", "/tracker_"+col+"/_doc/"+id, nil); err != nil || (code >= 300 && code != 404) {
		log.Printf("elasticsearch delete (%s) failed, status %d: %v", col, code, err)
	}
}

func Search(userID, q, typ, tag, from, to string) ([]map[string]any, error) {
	if !config.ESEnabled() {
		return nil, errors.New("Search is not configured (ELASTICSEARCH_URL is empty)")
	}
	filter := []any{map[string]any{"term": map[string]any{"userId.keyword": userID}}}
	if typ != "" {
		filter = append(filter, map[string]any{"term": map[string]any{"type.keyword": typ}})
	}
	if tag != "" {
		filter = append(filter, map[string]any{"term": map[string]any{"tags.keyword": tag}})
	}
	if from != "" || to != "" {
		r := map[string]any{}
		if from != "" {
			r["gte"] = from
		}
		if to != "" {
			r["lte"] = to
		}
		filter = append(filter, map[string]any{"range": map[string]any{"date": r}})
	}
	query := map[string]any{"bool": map[string]any{"filter": filter, "must": []any{
		map[string]any{"multi_match": map[string]any{"query": q, "fields": []string{"title^2", "content", "tags"}, "fuzziness": "AUTO"}}}}}
	res, code, err := config.ESRequest("POST", "/"+idxAll+"/_search?ignore_unavailable=true&allow_no_indices=true", map[string]any{"query": query, "size": 50})
	if err != nil || code >= 300 {
		return nil, errors.New("Search service unavailable")
	}
	out := []map[string]any{}
	hits, _ := res["hits"].(map[string]any)
	list, _ := hits["hits"].([]any)
	for _, x := range list {
		h, _ := x.(map[string]any)
		s, _ := h["_source"].(map[string]any)
		out = append(out, map[string]any{"id": h["_id"], "type": s["type"], "title": s["title"], "content": s["content"], "tags": s["tags"], "date": s["date"]})
	}
	return out, nil
}
