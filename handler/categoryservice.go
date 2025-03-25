package Server

import (
	Config "GRPC/Config"
	"GRPC/model"
	"GRPC/pb"
	"context"
	"errors"
	"log"
)

type CategoryServer struct {
	pb.UnimplementedCategoryServiceServer
}

func (s *CategoryServer) CreateCategory(ctx context.Context, in *pb.Category) (*pb.CategoryResponse, error) {
	category := &model.Category{
		ID:           int32(in.Id),
		CategoryName: in.CategoryName,
	}
	err := Config.DB.Create(&category)

	if err.Error != nil {
		return nil, errors.New("error saving category details ")
	}
	log.Printf("The Category Created with ID : %d", category.ID)

	return &pb.CategoryResponse{Category: in}, nil

}

func (s *CategoryServer) GetCategories(ctx context.Context, in *pb.Empty) (*pb.CategoryList, error) {
	var categories []model.Category
	err := Config.DB.Find(&categories)
	if err.Error != nil {
		return nil, err.Error
	}
	var pbcategory []*pb.Category
	for _, cat := range categories {
		pbcategory = append(pbcategory, &pb.Category{Id: int32(cat.ID), CategoryName: cat.CategoryName})
	}
	return &pb.CategoryList{Categories: pbcategory}, nil
}

func (s *CategoryServer) DeleteCategory(ctx context.Context, in *pb.CategoryRequest) (*pb.DeleteResponse, error) {
	var category model.Category
	res := Config.DB.First(&category, in.Id)
	if res.Error != nil {
		return nil, res.Error
	}
	err := Config.DB.Delete(&category)
	if err.Error != nil {
		return nil, err.Error
	}
	return &pb.DeleteResponse{Message: "Category record deleted successfully.."}, nil
}
