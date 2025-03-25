package Server

import (
	Config "GRPC/Config"
	"GRPC/model"
	"GRPC/pb"
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FarmerServer struct {
	pb.UnimplementedFarmerServiceServer
}

func (s *FarmerServer) CreateFarmer(ctx context.Context, in *pb.Farmer) (*pb.FarmerResponse, error) {
	farmer := &model.Farmer{
		Userid:      in.GetId(),
		DigitalId:   in.GetDigitalId(),
		LandHectare: in.GetLandInHectares(),
	}
	res := Config.DB.Create(farmer)

	if res.Error != nil {
		return nil, errors.New("error saving farmer details")
	}
	log.Printf("New Farmer Created with ID: %d", farmer.FarmerID)
	res1 := Config.DB.Preload("User").First(&farmer, farmer.FarmerID)
	if res1.Error != nil {
		return nil, res1.Error
	}
	result := &pb.FarmerResponse{
		Farmer: &pb.Farmer{
			FarmerId:       farmer.FarmerID,
			Id:             farmer.Userid,
			DigitalId:      farmer.DigitalId,
			LandInHectares: farmer.LandHectare,
		},
		User: &pb.User{
			Id:          int32(farmer.User.Userid),
			UserName:    farmer.User.Name,
			Email:       farmer.User.Email,
			PhoneNumber: farmer.User.Phone,
			Password:    farmer.User.Password,
			Role:        farmer.User.Role,
			Location:    farmer.User.Location,
		},
	}
	return result, nil
}
func (s *FarmerServer) UpdateFarmer(ctx context.Context, in *pb.UpdateFarmerRequest) (*pb.FarmerResponse, error) {
	var farmer model.Farmer

	res := Config.DB.Preload("User").First(&farmer, in.FarmerId)
	if res.Error != nil {
		log.Println("Farmer not found:", res.Error)
		return nil, errors.New("farmer not found")
	}

	farmer.DigitalId = in.DigitalId
	farmer.LandHectare = in.LandInHectares

	if in.UserName != "" {
		farmer.User.Name = in.UserName
	}
	if in.Email != "" {
		farmer.User.Email = in.Email
	}
	if in.PhoneNumber != "" {
		farmer.User.Phone = in.PhoneNumber
	}
	if in.Password != "" {
		farmer.User.Password = in.Password
	}
	if in.Location != "" {
		farmer.User.Location = in.Location
	}

	if err := Config.DB.Save(&farmer).Error; err != nil {
		return nil, err
	}
	if err := Config.DB.Save(&farmer.User).Error; err != nil {
		return nil, err
	}

	result := &pb.FarmerResponse{
		Farmer: &pb.Farmer{
			FarmerId:       farmer.FarmerID,
			Id:             farmer.Userid,
			DigitalId:      farmer.DigitalId,
			LandInHectares: farmer.LandHectare,
		},
		User: &pb.User{
			Id:          int32(farmer.User.Userid),
			UserName:    farmer.User.Name,
			Email:       farmer.User.Email,
			PhoneNumber: farmer.User.Phone,
			Password:    farmer.User.Password,
			Role:        farmer.User.Role,
			Location:    farmer.User.Location,
		},
	}
	return result, nil
}

func (s *FarmerServer) GetFarmerByID(ctx context.Context, in *pb.FarmerRequest) (*pb.FarmerResponse, error) {
	var farmer model.Farmer
	err := Config.DB.Preload("User").First(&farmer, in.FarmerId).Error
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Farmer not found")
	}
	result := &pb.FarmerResponse{
		Farmer: &pb.Farmer{
			FarmerId:       farmer.FarmerID,
			Id:             farmer.Userid,
			DigitalId:      farmer.DigitalId,
			LandInHectares: float32(farmer.LandHectare),
		},
		User: &pb.User{
			Id:          farmer.User.Userid,
			UserName:    farmer.User.Name,
			Email:       farmer.User.Email,
			PhoneNumber: farmer.User.Phone,
			Password:    farmer.User.Password,
			Role:        farmer.User.Role,
			Location:    farmer.User.Location,
		},
	}
	return result, nil
}

func (s *FarmerServer) DeleteFarmer(ctx context.Context, in *pb.FarmerRequest) (*pb.DeleteResponse, error) {
	var farmer model.Farmer
	data := Config.DB.First(&farmer, in.FarmerId)
	if data.Error != nil {
		return nil, errors.New("Farmer not found")
	}
	err := Config.DB.Delete(&farmer)
	if err.Error != nil {
		return nil, err.Error
	}
	return &pb.DeleteResponse{Message: "Farmer Record deleted successfully..!"}, nil

}
