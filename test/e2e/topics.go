/*
Copyright 2019 Google LLC

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

package e2e

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/pubsub"
	"knative.dev/pkg/test/helpers"

	// The following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func makeTopicOrDie(t *testing.T) (string, func()) {
	ctx := context.Background()
	// Prow sticks the project in this key
	project := os.Getenv(ProwProjectKey)
	if project == "" {
		t.Fatalf("failed to find %q in envvars", ProwProjectKey)
	}
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		t.Fatalf("failed to create pubsub client, %s", err.Error())
	}
	topicName := helpers.AppendRandomString("ps-e2e-test")
	topic := client.Topic(topicName)
	if exists, err := topic.Exists(ctx); err != nil {
		t.Fatalf("failed to verify topic exists, %s", err)
	} else if exists {
		t.Fatalf("topic already exists: %q", topicName)
	} else {
		topic, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			t.Fatalf("failed to create topic, %s", err)
		}
	}
	return topicName, func() {
		deleteTopicOrDie(t, topicName)
	}
}

func makeTopicOrDieWithTopicName(t *testing.T, topicName string) (string, func()) {
	ctx := context.Background()
	// Prow sticks the project in this key
	project := os.Getenv(ProwProjectKey)
	if project == "" {
		t.Fatalf("failed to find %q in envvars", ProwProjectKey)
	}
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		t.Fatalf("failed to create pubsub client, %s", err.Error())
	}
	topic := client.Topic(topicName)
	if exists, err := topic.Exists(ctx); err != nil {
		t.Fatalf("failed to verify topic exists, %s", err)
	} else if exists {
		t.Fatalf("topic already exists: %q", topicName)
	} else {
		topic, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			t.Fatalf("failed to create topic, %s", err)
		}
	}
	return topicName, func() {
		deleteTopicOrDie(t, topicName)
	}
}

func getTopic(t *testing.T, topicName string) *pubsub.Topic {
	ctx := context.Background()
	// Prow sticks the project in this key
	project := os.Getenv(ProwProjectKey)
	if project == "" {
		t.Fatalf("failed to find %q in envvars", ProwProjectKey)
	}
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		t.Fatalf("failed to create pubsub client, %s", err.Error())
	}
	return client.Topic(topicName)
}

func deleteTopicOrDie(t *testing.T, topicName string) {
	ctx := context.Background()
	// Prow sticks the project in this key.
	project := os.Getenv(ProwProjectKey)
	if project == "" {
		t.Fatalf("failed to find %q in envvars", ProwProjectKey)
	}
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		t.Fatalf("failed to create pubsub client, %s", err.Error())
	}
	topic := client.Topic(topicName)
	if exists, err := topic.Exists(ctx); err != nil {
		t.Fatalf("failed to verify topic exists, %s", err)
	} else if exists {
		if err := topic.Delete(ctx); err != nil {
			t.Fatalf("failed to delete topic, %s", err)
		}
	}
}
