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

# Setup env vars
export PROJECT_NAME="knative-gcp-performance"
export USER_NAME="mako-job@knative-gcp-performance.iam.gserviceaccount.com"
export TEST_ROOT_PATH="$GOPATH/src/github.com/google/knative-gcp/test/performance"
export KO_DOCKER_REPO="gcr.io/knative-gcp-performance"
export TEST_CONFIG_VARIANT="continuous"
export TEST_NAMESPACE="default"
export PUBSUB_SERVICE_ACCOUNT="performance-pubsub-test"
export PUBSUB_SERVICE_ACCOUNT_KEY="$(mktemp)"
export PUBSUB_SECRET_NAME="google-cloud-key"

function header() {
  echo "***** $1 *****"
}

# Exit script, dumping current state info.
# Parameters: $1 - error message (optional).
function abort() {
  [[ -n $1 ]] && echo "SCRIPT ERROR: $1"
  exit 1
}

# Creates a new cluster.
# $1 -> name, $2 -> zone/region, $3 -> num_nodes
function create_cluster() {
  header "Creating cluster $1 with $3 nodes in $2"
  gcloud beta container clusters create ${1} \
    --addons=HorizontalPodAutoscaling,HttpLoadBalancing \
    --machine-type=n1-standard-4 \
    --cluster-version=latest --region=${2} \
    --enable-stackdriver-kubernetes --enable-ip-alias \
    --num-nodes=${3} \
    --enable-autorepair \
    --scopes cloud-platform
}

# Create serice account secret on the cluster.
# $1 -> cluster_name, $2 -> cluster_zone
function create_secret() {
  echo "Create service account on cluster $1 in zone $2"
  gcloud container clusters get-credentials $1 --zone=$2 --project=${PROJECT_NAME} || abort "Failed to get cluster creds"
  kubectl create secret generic service-account --from-file=robot.json=${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}
}

# Set up the user credentials for cluster operations.
function setup_user() {
  header "Setup User"

  gcloud config set core/account ${USER_NAME}
  gcloud auth activate-service-account ${USER_NAME} --key-file=${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}
  gcloud config set core/project ${PROJECT_NAME}

  echo "gcloud user is $(gcloud config get-value core/account)"
  echo "Using secret defined in ${PERF_TEST_GOOGLE_APPLICATION_CREDENTIALS}"
}

# Get cluster credentials for GKE cluster
# $1 -> cluster_name, $2 -> cluster_zone
function get_gke_credentials() {
  name=$1
  zone=$2
  echo "Updating cluster with name ${name} in zone ${zone}"
  gcloud container clusters get-credentials ${name} --zone=${zone} --project=${PROJECT_NAME} || abort "Failed to get cluster creds"
}

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

# Install resources required for Mako. Assumes the benchmark is running in
# the kubectl context's default namespace.
function install_mako_resources() {
  echo ">> Setting up config-mako ConfigMap"
  cat <<EOF | kubectl apply -n ${TEST_NAMESPACE} -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: config-mako
data:
  # This should only be used by our performance automation.
  environment: prod
EOF
}

# Delete all the benchmark resources
# $1 -> cluster_name, $2 -> test config variant
function delete_benchmark_resources() {
  name=$1
  variant=$2

  echo ">> Delete all existing jobs and test resources"
  kubectl delete job --all
  ko delete -f "${TEST_ROOT_PATH}/${name}/${variant}/"
}

# Apply all the benchmark resources
# $1 -> cluster_name, $2 -> test config variant
function apply_benchmark_resources() {
  name=$1
  variant=$2

  echo ">> Applying $name benchmark yamls"
  # install the service and cronjob to run the benchmark
  # NOTE: this assumes we have a benchmark with the same name as the cluster
  # If service creation takes long time, we will have some initially unreachable errors in the test
  echo "Using ko version $(ko version)"
  ko apply -f "$TEST_ROOT_PATH/$name/${variant}/" || abort "Failed to apply ${name}/${variant} benchmark yamls"
}

# Update resources installed on the cluster with the up-to-date code. This
# assumes the benchmark is running in the kubectl context's default namespace.
# $1 -> cluster_name
function update_cluster() {
  name=$1

  install_knative_gcp_resources
  enable_cloud_pubsub
  install_mako_resources
  delete_benchmark_resources $name $TEST_CONFIG_VARIANT
  apply_benchmark_resources $name $TEST_CONFIG_VARIANT
}

# Create a new cluster and install serving components and apply benchmark yamls.
# $1 -> cluster_name, $2 -> cluster_zone, $3 -> node_count
function create_new_cluster() {
  name=$1
  zone=$2
  nodes=$3

  # create a new cluster
  create_cluster $name $zone $nodes || abort "Failed to create the new cluster $name"

  # create the secret on the new cluster
  create_secret $name $zone || abort "Failed to create secrets on the new cluster"

  # Setup user credentials to run on GKE for continous runs.
  get_gke_credentials $name $zone || abort "Failed to get GKE credentials for the new cluster"

  # update components on the cluster, e.g. knative-gcp
  update_cluster $name || abort "Failed to update the cluster"
}