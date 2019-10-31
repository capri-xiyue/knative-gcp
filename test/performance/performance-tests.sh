#!/bin/bash

# Copyright 2019 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# performance-tests.sh is added to manage all clusters that run the performance
# benchmarks in eventing repo, it is ONLY intended to be run by Prow, users
# should NOT run it manually.

# Setup env vars to override the default settings
export PROJECT_NAME="knative-eventing-performance"
export BENCHMARK_ROOT_PATH="$GOPATH/src/github.com/google/knative-gcp/test/performance/benchmarks"

source vendor/knative.dev/test-infra/scripts/performance-tests.sh
source ../lib.sh

# Vars used in this script
export TEST_CONFIG_VARIANT="continuous"
export TEST_NAMESPACE="default"
export PUBSUB_SERVICE_ACCOUNT="performance-pubsub-test"
export PUBSUB_SERVICE_ACCOUNT_KEY="$(mktemp)"
export PUBSUB_SECRET_NAME="google-cloud-key"

# Install the knative-gcp resources from the repo
function install_knative_gcp_resources() {
  pushd .
  cd ${GOPATH}/src/github.com/google/knative-gcp

  echo ">> Update knative-gcp core"
  ko apply --selector events.cloud.run/crd-install=true \
    -f config/ || abort "Failed to apply knative-gcp CRDs"

  ko apply \
    -f config/ || abort "Failed to apply knative-gcp resources"

  popd
}

# TODO(chizhg): this function is not needed after we set up the service account.
# Create resources required for Pub/Sub Admin setup
function enable_cloud_pubsub() {
  # When not running on Prow we need to set up a service account for PubSub
  echo "Set up ServiceAccount for Pub/Sub Admin"
  gcloud services enable pubsub.googleapis.com
  gcloud iam service-accounts create ${PUBSUB_SERVICE_ACCOUNT}
  gcloud projects add-iam-policy-binding ${PROJECT_NAME} \
    --member=serviceAccount:${PUBSUB_SERVICE_ACCOUNT}@${PROJECT_NAME}.iam.gserviceaccount.com \
    --role roles/pubsub.editor
  gcloud iam service-accounts keys create ${PUBSUB_SERVICE_ACCOUNT_KEY} \
    --iam-account=${PUBSUB_SERVICE_ACCOUNT}@${PROJECT_NAME}.iam.gserviceaccount.com
  service_account_key="${PUBSUB_SERVICE_ACCOUNT_KEY}"
  kubectl -n ${TEST_NAMESPACE} create secret generic ${PUBSUB_SECRET_NAME} --from-file=key.json=${service_account_key}
}

function update_knative() {
  start_latest_knative_eventing
  install_knative_gcp_resources
  enable_cloud_pubsub
}

function update_benchmark() {
  echo ">> Updating benchmark $1"
  ko delete -f ${BENCHMARK_ROOT_PATH}/$1/${TEST_CONFIG_VARIANT}
  ko apply -f ${BENCHMARK_ROOT_PATH}/$1/${TEST_CONFIG_VARIANT} || abort "failed to apply benchmark $1"
}

main $@