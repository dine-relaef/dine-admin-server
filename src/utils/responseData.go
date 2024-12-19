package utils

import (
	models "menu-server/src/models"
)

func RestaurantResponse(restaurant []models.Restaurant) []models.ResponseRestaurantData {
	var res_restaurants []models.ResponseRestaurantData
	for _, r := range restaurant {
		res_restaurants = append(res_restaurants, models.ResponseRestaurantData{
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
