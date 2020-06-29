// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang/protobuf/ptypes"
	empty "github.com/golang/protobuf/ptypes/empty"
	v1 "github.com/srinandan/sample-apps/tracking/pkg/api/v1"
	"go.opencensus.io/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Tracking struct {
	TrackingId      string    `json:"tracking_id,omitempty"`
	Status          string    `json:"status,omitempty"`
	CreateTime      time.Time `json:"create_time,omitempty"`
	UpdateTime      time.Time `json:"update_time,omitempty"`
	Signed          string    `json:"signed,omitempty"`
	Weight          string    `json:"weight,omitempty"`
	EstDeliveryTime time.Time `json:"est_delivery_time,omitempty"`
	Carrier         string    `json:"carrier,omitempty"`
}

var trackings = []Tracking{}

func ReadTrackingFile() error {
	trackingBytes, err := ioutil.ReadFile("tracking.json")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(trackingBytes, &trackings); err != nil {
		return err
	}

	return nil
}

// server is used to implement TrackingServer
type ShipmentServer struct {
	v1.UnimplementedShipmentServer
}

func NewShipmentService() (v1.ShipmentServer, error) {
	err := ReadTrackingFile()
	if err != nil {
		return &ShipmentServer{}, err
	}

	return &ShipmentServer{}, err
}

func (s *ShipmentServer) GetTracking(ctx context.Context, req *v1.GetTrackingRequest) (*v1.Tracking, error) {
	ctx, span := trace.StartSpan(ctx, "GetTracking")
	defer span.End()

	for _, tracking := range trackings {
		if tracking.TrackingId == req.TrackingId {
			return &v1.Tracking{
				TrackingId:      tracking.TrackingId,
				Status:          tracking.Status,
				CreateTime:      getTimestamp(tracking.CreateTime),
				UpdateTime:      getTimestamp(tracking.UpdateTime),
				Weight:          tracking.Weight,
				EstDeliveryTime: getTimestamp(tracking.EstDeliveryTime),
				Carrier:         tracking.Carrier,
			}, nil
		}
	}

	return &v1.Tracking{}, fmt.Errorf("tracking item not found")
}

func (s *ShipmentServer) ListTracking(ctx context.Context, empty *empty.Empty) (*v1.ListTrackingResponse, error) {
	ctx, span := trace.StartSpan(ctx, "ListTracking")
	defer span.End()

	listTrackingResponse := v1.ListTrackingResponse{}

	if len(trackings) == 0 {
		return &listTrackingResponse, fmt.Errorf("tracking items not found")
	}

	for _, tracking := range trackings {
		Tracking := v1.Tracking{}
		Tracking.TrackingId = tracking.TrackingId
		Tracking.Status = tracking.Status
		Tracking.CreateTime, _ = ptypes.TimestampProto(tracking.CreateTime)
		Tracking.UpdateTime, _ = ptypes.TimestampProto(tracking.UpdateTime)
		Tracking.Weight = tracking.Weight
		Tracking.EstDeliveryTime, _ = ptypes.TimestampProto(tracking.EstDeliveryTime)
		Tracking.Carrier = tracking.Carrier
		listTrackingResponse.Trackings = append(listTrackingResponse.Trackings, &Tracking)
	}
	return &listTrackingResponse, nil
}

func getTimestamp(t time.Time) *timestamppb.Timestamp {
	ts, _ := ptypes.TimestampProto(t)
	return ts
}
