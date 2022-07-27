package vimeo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestVideo_GetID(t *testing.T) {
	v := &Video{Name: "Test", URI: "/videos/1"}

	if id := v.GetID(); id != 1 {
		t.Errorf("Video.GetID returned %+v, want %+v", id, 1)
	}
}

func TestVideosService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormURLValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `{"data": [{"name": "Test"}]}`)
	})

	videos, _, err := client.Videos.List(OptPage(1), OptPerPage(2))
	if err != nil {
		t.Errorf("Videos.List returned unexpected error: %v", err)
	}

	want := []*Video{{Name: "Test"}}
	if !reflect.DeepEqual(videos, want) {
		t.Errorf("Videos.List returned %+v, want %+v", videos, want)
	}
}

func TestVideosService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name": "Test"}`)
	})

	video, _, err := client.Videos.Get(1)
	if err != nil {
		t.Errorf("Videos.Get returned unexpected error: %v", err)
	}

	want := &Video{Name: "Test"}
	if !reflect.DeepEqual(video, want) {
		t.Errorf("Videos.Get returned %+v, want %+v", video, want)
	}
}

func TestVideosService_Edit(t *testing.T) {
	setup()
	defer teardown()

	input := &VideoRequest{
		Name:        "name",
		Description: "desc",
	}

	mux.HandleFunc("/videos/1", func(w http.ResponseWriter, r *http.Request) {
		v := &VideoRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.Edit returned unexpected error: %v", err)
		}

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.Edit body is %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"name": "name"}`)
	})

	video, _, err := client.Videos.Edit(1, input)
	if err != nil {
		t.Errorf("Videos.Edit returned unexpected error: %v", err)
	}

	want := &Video{Name: "name"}
	if !reflect.DeepEqual(video, want) {
		t.Errorf("Videos.Edit returned %+v, want %+v", video, want)
	}
}

func TestVideosService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Videos.Delete(1)
	if err != nil {
		t.Errorf("Videos.Delete returned unexpected error: %v", err)
	}
}

func TestVideosService_ListCategory(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/categories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormURLValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `{"data": [{"name": "Test"}]}`)
	})

	categories, _, err := client.Videos.ListCategory(1, OptPage(1), OptPerPage(2))
	if err != nil {
		t.Errorf("Videos.ListCategory returned unexpected error: %v", err)
	}

	want := []*Category{{Name: "Test"}}
	if !reflect.DeepEqual(categories, want) {
		t.Errorf("Videos.ListCategory returned %+v, want %+v", categories, want)
	}
}

func TestVideosService_ListComment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormURLValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `{"data": [{"text": "Test"}]}`)
	})

	comments, _, err := client.Videos.ListComment(1, OptPage(1), OptPerPage(2))
	if err != nil {
		t.Errorf("Videos.ListComment returned unexpected error: %v", err)
	}

	want := []*Comment{{Text: "Test"}}
	if !reflect.DeepEqual(comments, want) {
		t.Errorf("Videos.ListComment returned %+v, want %+v", comments, want)
	}
}

func TestVideosService_GetComment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormURLValues(t, r, values{
			"fields": "name",
		})
		fmt.Fprint(w, `{"text": "Test"}`)
	})

	comment, _, err := client.Videos.GetComment(1, 1, OptFields([]string{"name"}))
	if err != nil {
		t.Errorf("Videos.GetComment returned unexpected error: %v", err)
	}

	want := &Comment{Text: "Test"}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Videos.GetComment returned %+v, want %+v", comment, want)
	}
}

func TestVideosService_AddComment(t *testing.T) {
	setup()
	defer teardown()

	input := &CommentRequest{
		Text: "name",
	}

	mux.HandleFunc("/videos/1/comments", func(w http.ResponseWriter, r *http.Request) {
		v := &CommentRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.AddComment returned unexpected error: %v", err)
		}

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.AddComment body is %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"text": "name"}`)
	})

	comment, _, err := client.Videos.AddComment(1, input)
	if err != nil {
		t.Errorf("Videos.AddComment returned unexpected error: %v", err)
	}

	want := &Comment{Text: "name"}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Videos.AddComment returned %+v, want %+v", comment, want)
	}
}

func TestVideosService_EditComment(t *testing.T) {
	setup()
	defer teardown()

	input := &CommentRequest{
		Text: "name",
	}

	mux.HandleFunc("/videos/1/comments/1", func(w http.ResponseWriter, r *http.Request) {
		v := &CommentRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.EditComment returned unexpected error: %v", err)
		}

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.EditComment body is %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"text": "name"}`)
	})

	comment, _, err := client.Videos.EditComment(1, 1, input)
	if err != nil {
		t.Errorf("Videos.EditComment returned unexpected error: %v", err)
	}

	want := &Comment{Text: "name"}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Videos.EditComment returned %+v, want %+v", comment, want)
	}
}

func TestVideosService_DeleteComment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Videos.DeleteComment(1, 1)
	if err != nil {
		t.Errorf("Videos.DeleteComment returned unexpected error: %v", err)
	}
}

func TestVideosService_ListReplies(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/comments/1/replies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormURLValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `{"data": [{"text": "Test"}]}`)
	})

	replies, _, err := client.Videos.ListReplies(1, 1, OptPage(1), OptPerPage(2))
	if err != nil {
		t.Errorf("Videos.ListReplies returned unexpected error: %v", err)
	}

	want := []*Comment{{Text: "Test"}}
	if !reflect.DeepEqual(replies, want) {
		t.Errorf("Videos.ListReplies returned %+v, want %+v", replies, want)
	}
}

func TestVideosService_AddReplies(t *testing.T) {
	setup()
	defer teardown()

	input := &CommentRequest{
		Text: "name",
	}

	mux.HandleFunc("/videos/1/comments/1/replies", func(w http.ResponseWriter, r *http.Request) {
		v := &CommentRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.AddReplies returned unexpected error: %v", err)
		}

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.AddReplies body is %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"text": "name"}`)
	})

	replies, _, err := client.Videos.AddReplies(1, 1, input)
	if err != nil {
		t.Errorf("Videos.AddReplies returned unexpected error: %v", err)
	}

	want := &Comment{Text: "name"}
	if !reflect.DeepEqual(replies, want) {
		t.Errorf("Videos.AddReplies returned %+v, want %+v", replies, want)
	}
}

func TestVideosService_ListCredit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/credits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormURLValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `{"data": [{"name": "Test"}]}`)
	})

	credits, _, err := client.Videos.ListCredit(1, OptPage(1), OptPerPage(2))
	if err != nil {
		t.Errorf("Videos.ListCredit returned unexpected error: %v", err)
	}

	want := []*Credit{{Name: "Test"}}
	if !reflect.DeepEqual(credits, want) {
		t.Errorf("Videos.ListCredit returned %+v, want %+v", credits, want)
	}
}

func TestVideosService_GetCredit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/credits/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormURLValues(t, r, values{
			"fields": "name",
		})
		fmt.Fprint(w, `{"name": "Test"}`)
	})

	credit, _, err := client.Videos.GetCredit(1, 1, OptFields([]string{"name"}))
	if err != nil {
		t.Errorf("Videos.GetCredit returned unexpected error: %v", err)
	}

	want := &Credit{Name: "Test"}
	if !reflect.DeepEqual(credit, want) {
		t.Errorf("Videos.GetCredit returned %+v, want %+v", credit, want)
	}
}

func TestVideosService_AddCredit(t *testing.T) {
	setup()
	defer teardown()

	input := &CreditRequest{
		Name: "name",
	}

	mux.HandleFunc("/videos/1/credits", func(w http.ResponseWriter, r *http.Request) {
		v := &CreditRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.AddCredit returned unexpected error: %v", err)
		}

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.AddCredit body is %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"name": "name"}`)
	})

	credit, _, err := client.Videos.AddCredit(1, input)
	if err != nil {
		t.Errorf("Videos.AddCredit returned unexpected error: %v", err)
	}

	want := &Credit{Name: "name"}
	if !reflect.DeepEqual(credit, want) {
		t.Errorf("Videos.AddCredit returned %+v, want %+v", credit, want)
	}
}

func TestVideosService_EditCredit(t *testing.T) {
	setup()
	defer teardown()

	input := &CreditRequest{
		Name: "name",
	}

	mux.HandleFunc("/videos/1/credits/1", func(w http.ResponseWriter, r *http.Request) {
		v := &CreditRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.EditCredit returned unexpected error: %v", err)
		}

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.EditCredit body is %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"name": "name"}`)
	})

	credit, _, err := client.Videos.EditCredit(1, 1, input)
	if err != nil {
		t.Errorf("Videos.EditCredit returned unexpected error: %v", err)
	}

	want := &Credit{Name: "name"}
	if !reflect.DeepEqual(credit, want) {
		t.Errorf("Videos.EditCredit returned %+v, want %+v", credit, want)
	}
}

func TestVideosService_DeleteCredit(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/credits/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Videos.DeleteCredit(1, 1)
	if err != nil {
		t.Errorf("Videos.DeleteCredit returned unexpected error: %v", err)
	}
}

func TestVideosService_ListPictures(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/pictures", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data": [{"uri": "Test"}]}`)
	})

	pictures, _, err := client.Videos.ListPictures(1)
	if err != nil {
		t.Errorf("Videos.ListPictures returned unexpected error: %v", err)
	}

	want := []*Pictures{{URI: "Test"}}
	if !reflect.DeepEqual(pictures, want) {
		t.Errorf("Videos.ListPictures returned %+v, want %+v", pictures, want)
	}
}

func TestVideosService_CreatePictures(t *testing.T) {
	setup()
	defer teardown()

	input := &PicturesRequest{
		Active: true,
	}

	mux.HandleFunc("/videos/1/pictures", func(w http.ResponseWriter, r *http.Request) {
		v := &PicturesRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.CreatePictures returned unexpected error: %v", err)
		}

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.CreatePictures body is %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"uri": "name"}`)
	})

	pictures, _, err := client.Videos.CreatePictures(1, input)
	if err != nil {
		t.Errorf("Videos.CreatePictures returned unexpected error: %v", err)
	}

	want := &Pictures{URI: "name"}
	if !reflect.DeepEqual(pictures, want) {
		t.Errorf("Videos.CreatePictures returned %+v, want %+v", pictures, want)
	}
}

func TestVideosService_GetPictures(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/pictures/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"uri": "Test"}`)
	})

	pictures, _, err := client.Videos.GetPictures(1, 1)
	if err != nil {
		t.Errorf("Videos.GetPictures returned unexpected error: %v", err)
	}

	want := &Pictures{URI: "Test"}
	if !reflect.DeepEqual(pictures, want) {
		t.Errorf("Videos.GetPictures returned %+v, want %+v", pictures, want)
	}
}

func TestVideosService_EditPictures(t *testing.T) {
	setup()
	defer teardown()

	input := &PicturesRequest{
		Active: true,
	}

	mux.HandleFunc("/videos/1/pictures/1", func(w http.ResponseWriter, r *http.Request) {
		v := &PicturesRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.EditPictures returned unexpected error: %v", err)
		}

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.EditPictures body is %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"uri": "name"}`)
	})

	pictures, _, err := client.Videos.EditPictures(1, 1, input)
	if err != nil {
		t.Errorf("Videos.EditPictures returned unexpected error: %v", err)
	}

	want := &Pictures{URI: "name"}
	if !reflect.DeepEqual(pictures, want) {
		t.Errorf("Videos.EditPictures returned %+v, want %+v", pictures, want)
	}
}

func TestVideosService_DeletePictures(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/pictures/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Videos.DeletePictures(1, 1)
	if err != nil {
		t.Errorf("Videos.DeletePictures returned unexpected error: %v", err)
	}
}

func TestVideosService_GetPreset(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/presets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name": "Test"}`)
	})

	preset, _, err := client.Videos.GetPreset(1, 1)
	if err != nil {
		t.Errorf("Videos.GetPreset returned unexpected error: %v", err)
	}

	want := &Preset{Name: "Test"}
	if !reflect.DeepEqual(preset, want) {
		t.Errorf("Videos.GetPreset returned %+v, want %+v", preset, want)
	}
}

func TestVideosService_AssignPreset(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/presets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	_, err := client.Videos.AssignPreset(1, 1)
	if err != nil {
		t.Errorf("Videos.AssignPreset returned unexpected error: %v", err)
	}
}

func TestVideosService_UnassignPreset(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/presets/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Videos.UnassignPreset(1, 1)
	if err != nil {
		t.Errorf("Videos.UnassignPreset returned unexpected error: %v", err)
	}
}

func TestVideosService_ListDomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/privacy/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data": [{"uri": "Test"}]}`)
	})

	domains, _, err := client.Videos.ListDomain(1)
	if err != nil {
		t.Errorf("Videos.ListDomain returned unexpected error: %v", err)
	}

	want := []*Domain{{URI: "Test"}}
	if !reflect.DeepEqual(domains, want) {
		t.Errorf("Videos.ListDomain returned %+v, want %+v", domains, want)
	}
}

func TestVideosService_AllowDomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/privacy/domains/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	_, err := client.Videos.AllowDomain(1, "1")
	if err != nil {
		t.Errorf("Videos.AllowDomain returned unexpected error: %v", err)
	}
}

func TestVideosService_DisallowDomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/privacy/domains/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Videos.DisallowDomain(1, "1")
	if err != nil {
		t.Errorf("Videos.DisallowDomain returned unexpected error: %v", err)
	}
}

func TestVideosService_ListUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/privacy/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data": [{"name": "Test"}]}`)
	})

	users, _, err := client.Videos.ListUser(1)
	if err != nil {
		t.Errorf("Videos.ListUser returned unexpected error: %v", err)
	}

	want := []*User{{Name: "Test"}}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("Videos.ListUser returned %+v, want %+v", users, want)
	}
}

func TestVideosService_AllowUsers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/privacy/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	_, err := client.Videos.AllowUsers(1)
	if err != nil {
		t.Errorf("Videos.AllowUsers returned unexpected error: %v", err)
	}
}

func TestVideosService_AllowUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/privacy/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	_, err := client.Videos.AllowUser(1, "1")
	if err != nil {
		t.Errorf("Videos.AllowDomain returned unexpected error: %v", err)
	}
}

func TestVideosService_DisallowUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/privacy/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Videos.DisallowUser(1, "1")
	if err != nil {
		t.Errorf("Videos.DisallowUser returned unexpected error: %v", err)
	}
}

func TestVideosService_ListTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/tags", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data": [{"uri": "Test"}]}`)
	})

	tags, _, err := client.Videos.ListTag(1)
	if err != nil {
		t.Errorf("Videos.ListTag returned unexpected error: %v", err)
	}

	want := []*Tag{{URI: "Test"}}
	if !reflect.DeepEqual(tags, want) {
		t.Errorf("Videos.ListTag returned %+v, want %+v", tags, want)
	}
}

func TestVideosService_GetTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/tags/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name": "Test"}`)
	})

	tag, _, err := client.Videos.GetTag(1, "1")
	if err != nil {
		t.Errorf("Videos.GetTag returned unexpected error: %v", err)
	}

	want := &Tag{Name: "Test"}
	if !reflect.DeepEqual(tag, want) {
		t.Errorf("Videos.GetTag returned %+v, want %+v", tag, want)
	}
}

func TestVideosService_AssignTagList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/tags", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		want := `[{"name":"1"},{"name":"2"}]`
		body, _ := ioutil.ReadAll(r.Body)
		if body := strings.TrimSpace(string(body)); body != want {
			t.Errorf("NewRequest Body is %v, want %v", body, want)
		}
	})

	_, err := client.Videos.AssignTagList(1, []string{"1", "2"})
	if err != nil {
		t.Errorf("Videos.AssignTagList returned unexpected error: %v", err)
	}
}

func TestVideosService_AssignTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/tags/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	_, err := client.Videos.AssignTag(1, "1")
	if err != nil {
		t.Errorf("Videos.AssignTag returned unexpected error: %v", err)
	}
}

func TestVideosService_UnassignTag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/tags/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Videos.UnassignTag(1, "1")
	if err != nil {
		t.Errorf("Videos.UnassignTag returned unexpected error: %v", err)
	}
}

func TestVideosService_ListTextTrack(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/texttracks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"data": [{"name": "Test"}]}`)
	})

	textTrack, _, err := client.Videos.ListTextTrack(1)
	if err != nil {
		t.Errorf("Videos.ListTextTrack returned unexpected error: %v", err)
	}

	want := []*TextTrack{{Name: "Test"}}
	if !reflect.DeepEqual(textTrack, want) {
		t.Errorf("Videos.ListTextTrack returned %+v, want %+v", textTrack, want)
	}
}

func TestVideosService_AddTextTrack(t *testing.T) {
	setup()
	defer teardown()

	input := &TextTrackRequest{
		Name: "name",
	}

	mux.HandleFunc("/videos/1/texttracks", func(w http.ResponseWriter, r *http.Request) {
		v := &TextTrackRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.AddTextTrack returned unexpected error: %v", err)
		}

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.AddTextTrack body is %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"name": "name"}`)
	})

	textTrack, _, err := client.Videos.AddTextTrack(1, input)
	if err != nil {
		t.Errorf("Videos.AddTextTrack returned unexpected error: %v", err)
	}

	want := &TextTrack{Name: "name"}
	if !reflect.DeepEqual(textTrack, want) {
		t.Errorf("Videos.AddTextTrack returned %+v, want %+v", textTrack, want)
	}
}

func TestVideosService_GetTextTrack(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/texttracks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"uri": "Test"}`)
	})

	textTrack, _, err := client.Videos.GetTextTrack(1, 1)
	if err != nil {
		t.Errorf("Videos.GetTextTrack returned unexpected error: %v", err)
	}

	want := &TextTrack{URI: "Test"}
	if !reflect.DeepEqual(textTrack, want) {
		t.Errorf("Videos.GetTextTrack returned %+v, want %+v", textTrack, want)
	}
}

func TestVideosService_EditTextTrack(t *testing.T) {
	setup()
	defer teardown()

	input := &TextTrackRequest{
		Name: "name",
	}

	mux.HandleFunc("/videos/1/texttracks/1", func(w http.ResponseWriter, r *http.Request) {
		v := &TextTrackRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.EditTextTrack returned unexpected error: %v", err)
		}

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.EditTextTrack body is %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"name": "name"}`)
	})

	textTrack, _, err := client.Videos.EditTextTrack(1, 1, input)
	if err != nil {
		t.Errorf("Videos.EditTextTrack returned unexpected error: %v", err)
	}

	want := &TextTrack{Name: "name"}
	if !reflect.DeepEqual(textTrack, want) {
		t.Errorf("Videos.EditTextTrack returned %+v, want %+v", textTrack, want)
	}
}

func TestVideosService_DeleteTextTrack(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/texttracks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Videos.DeleteTextTrack(1, 1)
	if err != nil {
		t.Errorf("Videos.DeleteTextTrack returned unexpected error: %v", err)
	}
}

func TestVideosService_ListRelatedVideo(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/videos/1/videos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormURLValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `{"data": [{"name": "Test"}]}`)
	})

	videos, _, err := client.Videos.ListRelatedVideo(1, OptPage(1), OptPerPage(2))
	if err != nil {
		t.Errorf("Videos.ListRelatedVideo returned unexpected error: %v", err)
	}

	want := []*Video{{Name: "Test"}}
	if !reflect.DeepEqual(videos, want) {
		t.Errorf("Videos.ListRelatedVideo returned %+v, want %+v", videos, want)
	}
}

func TestVideosService_getUploadVideo(t *testing.T) {
	setup()
	defer teardown()

	input := &UploadVideoRequest{
		Upload: &Upload{
			Approach: "tus",
			Size:     10,
		},
	}

	mux.HandleFunc("/me/videos", func(w http.ResponseWriter, r *http.Request) {
		v := &UploadVideoRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.getUploadVideo returned unexpected error: %v", err)
		}

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.getUploadVideo body is %+v, want %+v", v, input)
		}
	})

	uploadVideo, _, err := getUploadVideo(client, "POST", "/me/videos", input)
	if err != nil {
		t.Errorf("Videos.getUploadVideo returned unexpected error: %v", err)
	}

	want := &Video{}
	if !reflect.DeepEqual(uploadVideo, want) {
		t.Errorf("Videos.getUploadVideo returned %+v, want %+v", uploadVideo.URI, want.URI)
	}
}

func TestVideosService_uploadVideoByURL(t *testing.T) {
	setup()
	defer teardown()

	videoURL := "http://video.com/1.mp4"
	input := &UploadVideoRequest{
		Upload: &Upload{
			Approach: "pull",
			Link:     videoURL,
		},
	}

	mux.HandleFunc("/me/videos", func(w http.ResponseWriter, r *http.Request) {
		v := &UploadVideoRequest{
			Upload: &Upload{
				Approach: "pull",
				Link:     videoURL,
			},
		}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatalf("Videos.getUploadVideo returned unexpected error: %v", err)
		}

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Videos.getUploadVideo body is %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"name": "Test"}`)
	})

	video, _, err := uploadVideoByURL(client, "/me/videos", videoURL)
	if err != nil {
		t.Errorf("Videos.Get returned unexpected error: %v", err)
	}

	want := &Video{Name: "Test"}
	if !reflect.DeepEqual(video, want) {
		t.Errorf("Videos.Get returned %+v, want %+v", video, want)
	}
}
