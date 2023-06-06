package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"online_fashion_shop/api/model/rating"
	"online_fashion_shop/api/repository"
	"strconv"
	"time"
)

type RatingService interface {
	Create(ratingInfo *rating.Rating) error
	Get(id string) (*rating.Rating, error)
	List(productIds string) (*rating.Rating, error)
	Delete(id string) (*rating.Rating, error)
	Update(rating *rating.Rating) error
	GetProductAvrRate(productId string) (int, error)
	GetMultiProductAvrRate(productIds []string) (map[string]int, error)
}

type RatingServiceImpl struct {
	Repo  repository.ProductRatingRepository
	Cache *redis.Client
}

func (s *RatingServiceImpl) Create(ratingInfo *rating.Rating) error {
	return s.Repo.InsertOne(ratingInfo)
}

func (s *RatingServiceImpl) List(productIds string) (*rating.Rating, error) {
	//TODO implement me
	panic("implement me")
}

func NewRatingServiceImpl(repo repository.ProductRatingRepository, cache *redis.Client) RatingService {
	return &RatingServiceImpl{
		Repo:  repo,
		Cache: cache
	}

}

func (s *RatingServiceImpl) Get(id string) (*rating.Rating, error) {
	return s.Repo.Get(id)
}

func (s *RatingServiceImpl) List(productIds []string) ([]*rating.Rating, error) {
	option := rating.RateSearchOption{RateFor: productIds}
	return s.Repo.List(option)
}

func (s *RatingServiceImpl) Delete(id string) (*rating.Rating, error) {
	rate, err := s.Repo.Get(id)
	if err != nil {
		return nil, err
	}
	err = s.Repo.DeleteOne(id)
	if err == nil {
		s.invalidateCache(rate.RateFor)
	}
	return rate, err
}

func (s *RatingServiceImpl) Update(rating *rating.Rating) error {
	err := s.Repo.Update(rating)
	if err == nil {
		s.invalidateCache(rating.RateFor)
	}
	return err
}

func (s *RatingServiceImpl) GetProductAvrRate(productId string) (int, error) {
	cacheKey := fmt.Sprintf("rating.avr.%s", productId)

	// Try to get the average rating from the cache first
	avrRate, err := s.getAvrRateFromCache(cacheKey)
	if err == nil {
		return avrRate, nil
	}

	// If the average rating is not in the cache, get it from the repository
	avrRate, err = s.Repo.GetAvr(productId)
	if err != nil {
		return 0, err
	}

	// Cache the average rating with a TTL of 1 hour
	err = s.Cache.Set(context.TODO(), cacheKey, avrRate, time.Hour).Err()
	if err != nil {
		log.Printf("Failed to cache average rating for product %s: %v", productId, err)
	}

	return avrRate, nil
}

func (s *RatingServiceImpl) GetMultiProductAvrRate(productIds []string) (map[string]int, error) {
	avrRates := make(map[string]int)

	// Try to get the average ratings from the cache first
	cacheKeys := make([]string, len(productIds))
	for i, productId := range productIds {
		cacheKeys[i] = fmt.Sprintf("rating.avr.%s", productId)
	}

	cacheResults, err := s.Cache.MGet(context.TODO(), cacheKeys...).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to get average ratings from cache: %v", err)
	}

	// Iterate through the product IDs and cache results
	for i, productId := range productIds {
		// Try to get the average rating from the cache
		cacheResult := cacheResults[i]
		if cacheResult != nil {
			strAvrRate, ok := cacheResult.(string)
			if !ok {
				log.Printf("Invalid cache value for product %s: %v", productId, cacheResult)
			} else {
				avrRate, err := strconv.Atoi(strAvrRate)
				if err != nil {
					log.Printf("Invalid cache value for product %s: %v", productId, strAvrRate)
				} else {
					avrRates[productId] = avrRate
					continue
				}
			}
		}

		// If the average rating is not in the cache, get it from the repository
		avrRate, err := s.Repo.GetAvr(productId)
		if err != nil {
			return nil, err
		}

		// Cache the average rating with a TTL of 1 hour
		cacheKey := fmt.Sprintf("rating.avr.%s", productId)
		err = s.Cache.Set(context.TODO(), cacheKey, avrRate, time.Hour).Err()
		if err != nil {
			log.Printf("Failed to cache average rating for product %s: %v", productId, err)
		}

		avrRates[productId] = avrRate
	}

	return avrRates, nil
}

func (s *RatingServiceImpl) getAvrRateFromCache(cacheKey string) (int, error) {
	// Try to get the average rating from the cache
	strAvrRate, err := s.Cache.Get(context.TODO(), cacheKey).Result()
	if err == redis.Nil {
		// The average rating is not in the cache
		return 0, fmt.Errorf("average rating not found in cache for key %s", cacheKey)
	} else if err != nil {
		// There was an error retrieving the average rating from the cache
		return 0, fmt.Errorf("failed to get average rating from cache: %v", err)
	}

	// Convert the average rating to an integer
	avrRate, err := strconv.Atoi(strAvrRate)
	if err != nil {
		// The cached value could not be converted to an integer
		return 0, fmt.Errorf("invalid average rating value in cache for key %s: %v", cacheKey, err)
	}

	return avrRate, nil
}

func (s *RatingServiceImpl) invalidateCache(productId string) {
	cacheKey := fmt.Sprintf("rating.avr.%s", productId)
	err := s.Cache.Del(context.TODO(), cacheKey).Err()
	if err != nil {
		log.Printf("Failed to invalidate cache for product %s: %v", productId, err)
	}
}
