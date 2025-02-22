/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mockautoscaling

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"k8s.io/klog/v2"
)

func (m *MockAutoscaling) DescribeTagsWithContext(ctx aws.Context, request *autoscaling.DescribeTagsInput, opt ...request.Option) (*autoscaling.DescribeTagsOutput, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	response := &autoscaling.DescribeTagsOutput{}
	for _, g := range m.Groups {
		for _, tag := range g.Tags {
			allFiltersMatch := true
			for _, filter := range request.Filters {
				match := false
				switch aws.StringValue(filter.Name) {
				case "value":
					for _, v := range filter.Values {
						if aws.StringValue(tag.Value) == aws.StringValue(v) {
							match = true
						}
					}

				default:
					klog.Fatalf("Unsupported filter: %v", filter)
				}

				if !match {
					allFiltersMatch = false
					break
				}
			}

			if !allFiltersMatch {
				continue
			}

			response.Tags = append(response.Tags, tag)
		}
	}

	return response, nil
}

func (m *MockAutoscaling) DescribeTags(request *autoscaling.DescribeTagsInput) (*autoscaling.DescribeTagsOutput, error) {
	return m.DescribeTagsWithContext(context.TODO(), request)
}

func (m *MockAutoscaling) DescribeTagsRequest(*autoscaling.DescribeTagsInput) (*request.Request, *autoscaling.DescribeTagsOutput) {
	klog.Fatalf("Not implemented")
	return nil, nil
}

func (m *MockAutoscaling) DescribeTagsPagesWithContext(ctx aws.Context, request *autoscaling.DescribeTagsInput, callback func(*autoscaling.DescribeTagsOutput, bool) bool, options ...request.Option) error {
	// For the mock, we just send everything in one page
	page, err := m.DescribeTagsWithContext(ctx, request)
	if err != nil {
		return err
	}

	callback(page, false)

	return nil
}
