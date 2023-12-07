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
	GrowthBookService struct {
		endpoint string
		features growthbook.FeatureMap
	}
)

func NewGrowthBookService(endpoint string) *GrowthBookService {
	return &GrowthBookService{
		endpoint: endpoint,
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
	if s.features == nil {
		featureMap, err := s.fetchFeatureMap()
		if err != nil {
			return nil, err
		}
		s.features = growthbook.ParseFeatureMap(featureMap)
	}

	return s.features, nil
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

	s.features = growthbook.ParseFeatureMap(featureMap)
	fmt.Println("features loaded successfully:", string(featureMap))
	return nil
}
