package featureflag

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/growthbook/growthbook-golang"
	"log"
	"net/http"
)

type (
	Cache interface {
		Get(ctx context.Context, key string) ([]byte, error)
		Set(ctx context.Context, key string, value []byte) error
	}

	GrowthBookService struct {
		endpoint string
		cache    Cache
	}
)

func NewGrowthBookService(endpoint string, cache Cache) *GrowthBookService {
	return &GrowthBookService{
		endpoint: endpoint,
		cache:    cache,
	}
}

func (s *GrowthBookService) fetchFeatureMap() ([]byte, error) {
	resp, err := http.Get(s.endpoint)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body := struct {
		Features json.RawMessage `json:"features"`
	}{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return body.Features, nil
}

func (s *GrowthBookService) getFeatures() (growthbook.FeatureMap, error) {
	featureMap, err := s.cache.Get(context.TODO(), "features")
	if err != nil {
		fmt.Println("while fetching cache:", err)
		featureMap, err = s.fetchFeatureMap()
		if err != nil {
			return nil, err
		}
	}

	return growthbook.ParseFeatureMap(featureMap), nil
}

func (s *GrowthBookService) Eval(ctx context.Context, name string) bool {
	features, err := s.getFeatures()
	if err != nil {
		log.Println("while getting feature map:", err)
		return false
	}

	growthBookContext := growthbook.NewContext().
		WithFeatures(features).
		WithAttributes(growthbook.Attributes{
			"client-id": ctx.Value("client-id").(string),
		})

	gb := growthbook.New(growthBookContext)
	return gb.EvalFeature(name).On
}

func (s *GrowthBookService) Refresh() error {
	featureMap, err := s.fetchFeatureMap()
	if err != nil {
		return err
	}

	if err := s.cache.Set(context.TODO(), "features", featureMap); err != nil {
		return err
	}

	fmt.Println("features loaded successfully:", string(featureMap))
	return nil
}
