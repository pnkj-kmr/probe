package poll_test

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"probe/apiserver/handler/poll"
// 	M "probe/model"
// 	"probe/safemap"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/stretchr/testify/assert"
// )

// // MockDB implements M.DB interface for testing
// type MockDB struct {
// 	data      map[string][]byte
// 	count     uint64
// 	shouldErr bool
// }

// func NewMockDB() *MockDB {
// 	return &MockDB{
// 		data: make(map[string][]byte),
// 	}
// }

// func (m *MockDB) Create(id string, data []byte) error {
// 	return m.Save(id, data)
// }

// func (m *MockDB) Update(id string, data []byte) error {
// 	return m.Save(id, data)
// }

// func (m *MockDB) Save(id string, data []byte) error {
// 	if m.shouldErr {
// 		return errors.New("mock save error")
// 	}
// 	m.data[id] = data
// 	m.count++
// 	return nil
// }

// func (m *MockDB) Find(id string) ([]byte, error) {
// 	if m.shouldErr {
// 		return nil, errors.New("mock find error")
// 	}
// 	if data, ok := m.data[id]; ok {
// 		return data, nil
// 	}
// 	return nil, nil
// }

// func (m *MockDB) Delete(id string) error {
// 	if m.shouldErr {
// 		return errors.New("mock delete error")
// 	}
// 	if _, ok := m.data[id]; ok {
// 		delete(m.data, id)
// 		m.count--
// 		return nil
// 	}
// 	return errors.New("not found")
// }

// func (m *MockDB) DeleteAll() error {
// 	if m.shouldErr {
// 		return errors.New("mock delete all error")
// 	}
// 	m.data = make(map[string][]byte)
// 	m.count = 0
// 	return nil
// }

// func (m *MockDB) Count() uint64 {
// 	return m.count
// }

// func (m *MockDB) Receive() <-chan []byte {
// 	ch := make(chan []byte)
// 	go func() {
// 		defer close(ch)
// 		for _, data := range m.data {
// 			ch <- data
// 		}
// 	}()
// 	return ch
// }

// func setupTestRouter(t *testing.T) (*poll.R, *safemap.SafeMap[int, M.DB]) {
// 	db := safemap.New[int, M.DB]()
// 	mockDB := NewMockDB()
// 	db.Set(M.ICMP+M.INTERVAL_60, mockDB)
// 	db.Set(M.ICMP+M.INTERVAL_300, mockDB)
// 	router := poll.NewRouter(M.ICMP, db)
// 	return router, db
// }

// func getMockDB(db *safemap.SafeMap[int, M.DB], key int) *MockDB {
// 	if dbObj, ok := db.Get(key); ok {
// 		return dbObj.(*MockDB)
// 	}
// 	return nil
// }

// func TestNewRouter(t *testing.T) {
// 	db := safemap.New[int, M.DB]()
// 	router := poll.NewRouter(M.ICMP, db)

// 	assert.NotNil(t, router)
// 	assert.Equal(t, M.ICMP, router.ID)
// 	assert.NotNil(t, router.Mux())
// 	assert.Equal(t, db, router.DB)
// }

// func TestGetDB(t *testing.T) {
// 	router, _ := setupTestRouter(t)

// 	// Test with string poll period
// 	dbObj, err := router.GetDB("60")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, dbObj)

// 	// Test with invalid string poll period
// 	dbObj, err = router.GetDB("invalid")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, dbObj)

// 	// Test with integer poll period
// 	dbObj, err = router.GetDB(300)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, dbObj)
// }

// func TestDelete(t *testing.T) {
// 	router, db := setupTestRouter(t)
// 	mockDB := getMockDB(db, M.ICMP+M.INTERVAL_60)

// 	// Add test data
// 	testData := []byte(`{"id": "test-id", "data": "test-data"}`)
// 	mockDB.Save("test-id", testData)

// 	// Create a test request
// 	req := httptest.NewRequest("DELETE", "/60/test-id", nil)
// 	rctx := chi.NewRouteContext()
// 	rctx.URLParams.Add("pollPeriod", "60")
// 	rctx.URLParams.Add("id", "test-id")
// 	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

// 	w := httptest.NewRecorder()
// 	router.Delete(w, req)

// 	assert.Equal(t, http.StatusNoContent, w.Code)
// 	assert.Equal(t, uint64(0), mockDB.Count())
// }

// func TestDeleteAll(t *testing.T) {
// 	router, db := setupTestRouter(t)
// 	mockDB := getMockDB(db, M.ICMP+M.INTERVAL_60)

// 	// Add test data
// 	testData := []byte(`{"id": "test-id", "data": "test-data"}`)
// 	mockDB.Save("test-id", testData)

// 	// Create a test request
// 	req := httptest.NewRequest("DELETE", "/60", nil)
// 	rctx := chi.NewRouteContext()
// 	rctx.URLParams.Add("pollPeriod", "60")
// 	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

// 	w := httptest.NewRecorder()
// 	router.DeleteAll(w, req)

// 	assert.Equal(t, http.StatusNoContent, w.Code)
// 	assert.Equal(t, uint64(0), mockDB.Count())
// }

// func TestCount(t *testing.T) {
// 	router, db := setupTestRouter(t)
// 	mockDB := getMockDB(db, M.ICMP+M.INTERVAL_60)

// 	// Add test data
// 	testData := []byte(`{"id": "test-id", "data": "test-data"}`)
// 	mockDB.Save("test-id", testData)

// 	// Create a test request
// 	req := httptest.NewRequest("GET", "/60", nil)
// 	rctx := chi.NewRouteContext()
// 	rctx.URLParams.Add("pollPeriod", "60")
// 	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

// 	w := httptest.NewRecorder()
// 	router.Count(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var count uint64
// 	err := json.Unmarshal(w.Body.Bytes(), &count)
// 	assert.NoError(t, err)
// 	assert.Equal(t, uint64(1), count)
// }

// func TestGet(t *testing.T) {
// 	router, db := setupTestRouter(t)
// 	mockDB := getMockDB(db, M.ICMP+M.INTERVAL_60)

// 	// Add test data
// 	testData := []byte(`{"id": "test-id", "data": "test-data"}`)
// 	mockDB.Save("test-id", testData)

// 	// Create a test request
// 	req := httptest.NewRequest("GET", "/60/test-id", nil)
// 	rctx := chi.NewRouteContext()
// 	rctx.URLParams.Add("pollPeriod", "60")
// 	rctx.URLParams.Add("id", "test-id")
// 	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

// 	w := httptest.NewRecorder()
// 	router.Get(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Equal(t, testData, w.Body.Bytes())
// }

// func TestGetNotFound(t *testing.T) {
// 	router, _ := setupTestRouter(t)

// 	// Create a test request for non-existent ID
// 	req := httptest.NewRequest("GET", "/60/non-existent", nil)
// 	rctx := chi.NewRouteContext()
// 	rctx.URLParams.Add("pollPeriod", "60")
// 	rctx.URLParams.Add("id", "non-existent")
// 	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

// 	w := httptest.NewRecorder()
// 	router.Get(w, req)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// }

// func TestSave(t *testing.T) {
// 	router, _ := setupTestRouter(t)

// 	// Create test data
// 	testData := map[string]interface{}{
// 		"id":   "test-id",
// 		"data": "test-data",
// 	}
// 	jsonData, _ := json.Marshal(testData)

// 	// Create a test request
// 	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.Save(w, req)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// }

// func TestDBError(t *testing.T) {
// 	router, db := setupTestRouter(t)
// 	mockDB := getMockDB(db, M.ICMP+M.INTERVAL_60)
// 	mockDB.shouldErr = true

// 	// Test Get with DB error
// 	req := httptest.NewRequest("GET", "/60/test-id", nil)
// 	rctx := chi.NewRouteContext()
// 	rctx.URLParams.Add("pollPeriod", "60")
// 	rctx.URLParams.Add("id", "test-id")
// 	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

// 	w := httptest.NewRecorder()
// 	router.Get(w, req)

// 	assert.Equal(t, http.StatusBadRequest, w.Code)
// }
