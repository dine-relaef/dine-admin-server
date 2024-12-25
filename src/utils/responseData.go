package utils

import (
	models_restaurant "dine-server/src/models/restaurants"
)

func RestaurantResponse(restaurant []models_restaurant.Restaurant) []models_restaurant.ResponseRestaurantData {
	var res_restaurants []models_restaurant.ResponseRestaurantData
	for _, r := range restaurant {
		res_restaurants = append(res_restaurants, models_restaurant.ResponseRestaurantData{
			ID:             r.ID,
			Name:           r.Name,
			PureVeg:        r.PureVeg,
			Description:    r.Description,
			Location:       r.Location,
			BannerImageUrl: r.BannerImageUrl,
			LogoImageUrl:   r.LogoImageUrl,
			Phone:          r.Phone,
			Email:          r.Email,
		})
	}
	return res_restaurants
}
