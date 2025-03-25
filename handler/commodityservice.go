package Server

import (
	Config "GRPC/Config"
	"GRPC/model"
	"GRPC/pb"
	"context"
	"database/sql"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommodityServer struct {
	pb.UnimplementedCommodityServiceServer
	db *sql.DB
}

func (s *CommodityServer) CreateCommodity(ctx context.Context, in *pb.Commodity) (*pb.CommodityResponse, error) {
	commodity := &model.Commodity{
		CommodityID:  int32(in.Id),
		ProductName:  in.ProductName,
		FarmerID:     int32(in.FarmerId),
		Quantity:     int32(in.Quantity),
		BasePrice:    float64(in.BasePrice),
		Availability: in.Availability,
		CategoryID:   int32(in.CategoryId),
	}
	err := Config.DB.Create(&commodity)
	if err.Error != nil {
		return nil, errors.New("error saving commodity details..")
	}
	res2 := Config.DB.Preload("Category").First(&commodity, commodity.CommodityID).Error
	if res2 != nil {
		return nil, res2
	}
	log.Printf("The commodity created with ID %v", commodity.CommodityID)
	response := &pb.CommodityResponse{
		Commodity: &pb.Commodity{
			Id:           commodity.CommodityID,
			Availability: commodity.Availability,
			BasePrice:    float32(commodity.BasePrice),
			ProductName:  commodity.ProductName,
			FarmerId:     commodity.FarmerID,
			Quantity:     commodity.Quantity,
			CategoryId:   commodity.CategoryID,
		},
		Category: &pb.Category{
			Id:           commodity.CategoryID,
			CategoryName: commodity.Category.CategoryName,
		},
	}

	return response, nil
}

func (s *CommodityServer) GetCommodities(ctx context.Context, in *pb.Empty) (*pb.CommodityList, error) {
	var commodities []model.Commodity
	err := Config.DB.Find(&commodities)
	if err.Error != nil {
		return nil, err.Error
	}
	var pbcommodity []*pb.Commodity
	for _, com := range commodities {
		pbcommodity = append(pbcommodity, &pb.Commodity{Id: com.CommodityID, ProductName: com.ProductName, FarmerId: com.FarmerID, Quantity: com.Quantity, BasePrice: float32(com.BasePrice), Availability: com.Availability, CategoryId: com.CategoryID})
	}
	return &pb.CommodityList{Commodities: pbcommodity}, nil
}

func (s *CommodityServer) DeleteCommodity(ctx context.Context, in *pb.CommodityRequest) (*pb.DeleteResponse, error) {
	var commodity model.Commodity
	res := Config.DB.First(&commodity, in.Id)
	if res.Error != nil {
		log.Println("Commodity not found:", res.Error)
		return nil, errors.New("commodity not found")
	}
	if err := Config.DB.Delete(&commodity).Error; err != nil {
		return nil, err
	}
	return &pb.DeleteResponse{Message: "Commodity deleted successfully"}, nil
}

func (s *CommodityServer) UpdateCommodity(ctx context.Context, in *pb.UpdateCommodityReq) (*pb.CommodityResponse, error) {
	var commodity model.Commodity

	res := Config.DB.Preload("Category").First(&commodity, in.Id)
	if res.Error != nil {
		log.Println("Commodity not found:", res.Error)
		return nil, errors.New("commodity not found")
	}

	commodity.Availability = in.Availability
	commodity.BasePrice = float64(in.BasePrice)
	commodity.ProductName = in.ProductName
	commodity.Quantity = in.Quantity
	if in.CategoryName != "" {
		commodity.Category.CategoryName = in.CategoryName
	}

	if err := Config.DB.Save(&commodity).Error; err != nil {
		return nil, err
	}
	if err := Config.DB.Save(&commodity.Category).Error; err != nil {
		return nil, err
	}

	response := &pb.CommodityResponse{
		Commodity: &pb.Commodity{
			Id:           commodity.CommodityID,
			Availability: commodity.Availability,
			BasePrice:    float32(commodity.BasePrice),
			ProductName:  commodity.ProductName,
			FarmerId:     commodity.FarmerID,
			Quantity:     commodity.Quantity,
			CategoryId:   commodity.CategoryID,
		},
		Category: &pb.Category{
			Id:           commodity.CategoryID,
			CategoryName: commodity.Category.CategoryName,
		},
	}

	return response, nil
}

func (s *CommodityServer) GetCommodityByID(ctx context.Context, in *pb.CommodityRequest) (*pb.CommodityResponse, error) {
	var commodity model.Commodity
	err := Config.DB.Preload("Category").First(&commodity, in.Id).Error
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Farmer not found")
	}
	result := &pb.CommodityResponse{
		Commodity: &pb.Commodity{
			Id:           commodity.CommodityID,
			ProductName:  commodity.ProductName,
			FarmerId:     commodity.FarmerID,
			Quantity:     commodity.Quantity,
			BasePrice:    float32(commodity.BasePrice),
			Availability: commodity.Availability,
			CategoryId:   commodity.CategoryID,
		},
		Category: &pb.Category{
			CategoryName: commodity.Category.CategoryName,
		},
	}

	return result, nil
}
