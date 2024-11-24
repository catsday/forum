package handlers

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDatabase struct {
	GetDataFunc           func(id int) (string, error)
	InsertDataFunc        func(data string) error
	GetCategoryFunc       func(categoryID int) (string, error)
	AddPostToCategoryFunc func(postID, categoryID int) error
}

func (m *MockDatabase) GetData(id int) (string, error) {
	return m.GetDataFunc(id)
}

func (m *MockDatabase) InsertData(data string) error {
	return m.InsertDataFunc(data)
}

func (m *MockDatabase) GetCategory(categoryID int) (string, error) {
	return m.GetCategoryFunc(categoryID)
}

func (m *MockDatabase) AddPostToCategory(postID, categoryID int) error {
	return m.AddPostToCategoryFunc(postID, categoryID)
}

// handler for fetching data by id
func GetHandler(db *MockDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := 1
		data, err := db.GetData(id)
		if err != nil {
			http.Error(w, "Error fetching data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	}
}

// handler for inserting new data
func PostHandler(db *MockDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		if len(body) == 0 {
			http.Error(w, "Empty body", http.StatusBadRequest)
			return
		}
		err := db.InsertData(string(body))
		if err != nil {
			http.Error(w, "Error inserting data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

// test for valid response from get handler
func TestGetHandler_ValidResponse(t *testing.T) {
	mockDB := &MockDatabase{
		GetDataFunc: func(id int) (string, error) {
			return "Mocked Data", nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/data", nil)
	rec := httptest.NewRecorder()

	handler := GetHandler(mockDB)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Mocked Data", rec.Body.String())
	assert.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
}

// test for handling empty data in get handler
func TestGetHandler_NotFound(t *testing.T) {
	mockDB := &MockDatabase{
		GetDataFunc: func(id int) (string, error) {
			return "", nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/data", nil)
	rec := httptest.NewRecorder()

	handler := GetHandler(mockDB)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "", rec.Body.String())
}

// test for valid data insertion in post handler
func TestPostHandler_ValidData(t *testing.T) {
	mockDB := &MockDatabase{
		InsertDataFunc: func(data string) error {
			assert.Equal(t, "New Data", data)
			return nil
		},
	}

	body := bytes.NewBufferString("New Data")
	req := httptest.NewRequest(http.MethodPost, "/data", body)
	rec := httptest.NewRecorder()

	handler := PostHandler(mockDB)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

// test for handling empty body in post handler
func TestPostHandler_EmptyBody(t *testing.T) {
	mockDB := &MockDatabase{
		InsertDataFunc: func(data string) error {
			t.FailNow()
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/data", nil)
	rec := httptest.NewRecorder()

	handler := PostHandler(mockDB)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Empty body")
}

// test for multiple requests to get handler
func TestGetHandler_MultipleRequests(t *testing.T) {
	mockDB := &MockDatabase{
		GetDataFunc: func(id int) (string, error) {
			return "Mocked Data", nil
		},
	}

	handler := GetHandler(mockDB)

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/data", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Mocked Data", rec.Body.String())
	}
}

// test for handling multiple post requests with valid data
func TestPostHandler_MultipleRequests(t *testing.T) {
	mockDB := &MockDatabase{
		InsertDataFunc: func(data string) error {
			return nil
		},
	}

	handler := PostHandler(mockDB)

	for i := 0; i < 5; i++ {
		body := bytes.NewBufferString("Data " + strconv.Itoa(i))
		req := httptest.NewRequest(http.MethodPost, "/data", body)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

// test for handling empty string in post handler
func TestPostHandler_EmptyString(t *testing.T) {
	mockDB := &MockDatabase{
		InsertDataFunc: func(data string) error {
			assert.Equal(t, "", data)
			return nil
		},
	}

	body := bytes.NewBufferString("")
	req := httptest.NewRequest(http.MethodPost, "/data", body)
	rec := httptest.NewRecorder()

	handler := PostHandler(mockDB)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Empty body")
}

// test for checking headers set in post handler
func TestPostHandler_CheckHeaders(t *testing.T) {
	mockDB := &MockDatabase{
		InsertDataFunc: func(data string) error {
			return nil
		},
	}

	body := bytes.NewBufferString("Test Data")
	req := httptest.NewRequest(http.MethodPost, "/data", body)
	rec := httptest.NewRecorder()

	handler := PostHandler(mockDB)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

// test for linking a post to a category
func AddPostToCategoryHandler(db *MockDatabase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postID := 1
		categoryID := 2

		err := db.AddPostToCategoryFunc(postID, categoryID)
		if err != nil {
			http.Error(w, "Error linking post to category", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}
}

// test for successful post-category linking
func TestAddPostToCategory_Success(t *testing.T) {
	mockDB := &MockDatabase{
		AddPostToCategoryFunc: func(postID, categoryID int) error {
			assert.Equal(t, 1, postID)
			assert.Equal(t, 2, categoryID)
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/add-to-category", nil)
	rec := httptest.NewRecorder()

	handler := AddPostToCategoryHandler(mockDB)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

// test for failure during post-category linking
func TestAddPostToCategory_Failure(t *testing.T) {
	mockDB := &MockDatabase{
		AddPostToCategoryFunc: func(postID, categoryID int) error {
			return errors.New("linking error")
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/add-to-category", nil)
	rec := httptest.NewRecorder()

	handler := AddPostToCategoryHandler(mockDB)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error linking post to category")
}

// test for fetching a category successfully
func TestGetCategory_Success(t *testing.T) {
	mockDB := &MockDatabase{
		GetCategoryFunc: func(categoryID int) (string, error) {
			assert.Equal(t, 1, categoryID)
			return "Technology", nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/category?id=1", nil)
	rec := httptest.NewRecorder()

	handler := func(w http.ResponseWriter, r *http.Request) {
		categoryID := 1
		data, err := mockDB.GetCategory(categoryID)
		if err != nil {
			http.Error(w, "Error fetching category", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
	}
	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Technology", rec.Body.String())
}

// test for handling category not found error
func TestGetCategory_NotFound(t *testing.T) {
	mockDB := &MockDatabase{
		GetCategoryFunc: func(categoryID int) (string, error) {
			return "", errors.New("not found")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/category?id=99", nil)
	rec := httptest.NewRecorder()

	handler := func(w http.ResponseWriter, r *http.Request) {
		categoryID := 99
		_, err := mockDB.GetCategory(categoryID)
		if err != nil {
			http.Error(w, "Error fetching category", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	handler(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error fetching category")
}
