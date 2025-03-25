package Server

import (
	Config "GRPC/Config"
	"GRPC/model"
	"GRPC/pb"
	"context"
	"errors"
	"log"
)

type BidServer struct {
	pb.UnimplementedBidServiceServer
}

func (s *BidServer) CreateBid(ctx context.Context, in *pb.Bid) (*pb.BidResponse, error) {
	bid := &model.Bidding{
		CommodityID: int32(in.CommodityId),
		Userid:      int32(in.UserId),
		BidAmount:   float32(in.BidAmount),
		BidStatus:   in.Status,
	}

	res := Config.DB.Create(bid)
	if res.Error != nil {
		return nil, errors.New("error creating bid")
	}

	log.Printf("New Bid Created with ID: %d", bid.Bidid)

	return &pb.BidResponse{
		Bid: &pb.Bid{
			BidId:       bid.Bidid,
			CommodityId: bid.CommodityID,
			UserId:      bid.Userid,
			BidAmount:   bid.BidAmount,
			Status:      bid.BidStatus,
		}, Message: "Bid created..",
	}, nil
}

func (s *BidServer) UpdateBid(ctx context.Context, in *pb.Bid) (*pb.BidResponse, error) {
	var bid model.Bidding

	if err := Config.DB.First(&bid, in.BidId).Error; err != nil {
		return nil, errors.New("bid not found")
	}

	bid.BidAmount = float32(in.BidAmount)
	bid.BidStatus = in.Status

	if err := Config.DB.Save(&bid).Error; err != nil {
		return nil, errors.New("error updating bid")
	}

	return &pb.BidResponse{
		Bid: &pb.Bid{
			BidId:       bid.Bidid,
			CommodityId: bid.CommodityID,
			UserId:      bid.Userid,
			BidAmount:   bid.BidAmount,
			Status:      bid.BidStatus,
		}, Message: "Bid updated..",
	}, nil
}

func (s *BidServer) GetBidByCommodityId(ctx context.Context, in *pb.BidRequest) (*pb.BidList, error) {
	var bids []model.Bidding

	if err := Config.DB.Where("commodity_id = ?", in.CommodityId).Find(&bids).Error; err != nil {
		return nil, errors.New("error fetching bids")
	}

	var bidList []*pb.Bid
	for _, bid := range bids {
		bidList = append(bidList, &pb.Bid{
			BidId:       bid.Bidid,
			CommodityId: bid.CommodityID,
			UserId:      bid.Userid,
			BidAmount:   bid.BidAmount,
			Status:      bid.BidStatus,
		})
	}

	return &pb.BidList{Bids: bidList}, nil
}
func (s *BidServer) BidAccept(ctx context.Context, in *pb.BidRequest) (*pb.BidResponse, error) {
	var bid model.Bidding
	var commodity model.Commodity
	if err := Config.DB.First(&bid, in.BidId).Error; err != nil {
		return nil, errors.New("bid not found")
	}
	if err := Config.DB.First(&bid, in.BidId).Error; err != nil {
		return nil, errors.New("bid not found")
	}

	if err := Config.DB.Model(&bid).Update("bid_status", "accepted").Error; err != nil {
		return nil, errors.New("failed to update bid status")
	}
	if err := Config.DB.Model(&model.Bidding{}).
		Where("commodity_id = ? AND bidid != ?", bid.CommodityID, bid.Bidid).
		Update("bid_status", "rejected").Error; err != nil {
		return nil, errors.New("failed to reject other bids")
	}

	if err := Config.DB.First(&commodity, bid.CommodityID).Error; err != nil {
		return nil, errors.New("commodity not found")
	}

	commodity.Availability = false
	if err := Config.DB.Save(&commodity).Error; err != nil {
		return nil, errors.New("failed to update commodity availability")
	}
	return &pb.BidResponse{
		Bid: &pb.Bid{
			BidId:       bid.Bidid,
			CommodityId: bid.CommodityID,
			UserId:      bid.Userid,
			BidAmount:   bid.BidAmount,
			Status:      "accepted",
		},
		Message: "Bid accepted successfully, commodity marked unavailable",
	}, nil
}

func (s *BidServer) DeleteBid(ctx context.Context, in *pb.BidRequest) (*pb.BidResponse, error) {
	var bid model.Bidding

	if err := Config.DB.Where("bidid = ?", in.BidId).Delete(&bid).Error; err != nil {
		return nil, errors.New("error deleting bid")
	}

	return &pb.BidResponse{
		Message: "Bid deleted successfully",
	}, nil
}
