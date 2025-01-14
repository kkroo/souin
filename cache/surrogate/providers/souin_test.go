package providers

import (
	"net/http"
	"sync"
	"testing"

	"github.com/darkweak/souin/configurationtypes"
	"github.com/darkweak/souin/errors"
	"go.uber.org/zap"
)

func mockSouinProvider() *SouinSurrogateStorage {
	sss := &SouinSurrogateStorage{
		baseStorage: &baseStorage{
			Storage:    &sync.Map{},
			Keys:       make(map[string]configurationtypes.SurrogateKeys),
			keysRegexp: make(map[string]keysRegexpInner),
			dynamic:    true,
			mu:         &sync.Mutex{},
			logger:     zap.NewNop(),
		},
	}

	sss.baseStorage.parent = sss

	return sss
}

func TestSouinSurrogateStorage_Store(t *testing.T) {
	sp := mockSouinProvider()
	res := http.Response{
		Header: http.Header{},
	}
	res.Header.Set(surrogateKey, baseHeaderValue)
	res.Header.Set(surrogateControl, "public, max-age=5")

	var e error
	if e = sp.Store(&res, "stored"); e != nil {
		errors.GenerateError(t, "It should not throw an error while store.")
	}

	if res.Header.Get(surrogateKey) != "test0, test1,   test2,  test3, test4" {
		errors.GenerateError(t, "The response should contains the Surrogate-keys header.")
	}

	if res.Header.Get(surrogateControl) != "public, max-age=5" {
		errors.GenerateError(t, "The response should contains the Surrogate-control header.")
	}
}
