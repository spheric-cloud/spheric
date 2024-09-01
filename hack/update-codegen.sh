#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
export TERM="xterm-256color"

bold="$(tput bold)"
blue="$(tput setaf 4)"
normal="$(tput sgr0)"

function qualify-gvs() {
  APIS_PKG="$1"
  GROUPS_WITH_VERSIONS="$2"
  join_char=""
  res=""

  for GVs in ${GROUPS_WITH_VERSIONS}; do
    IFS=: read -r G Vs <<<"${GVs}"

    for V in ${Vs//,/ }; do
      res="$res$join_char$APIS_PKG/$G/$V"
      join_char=","
    done
  done

  echo "$res"
}

function qualify-gs() {
  APIS_PKG="$1"
  unset GROUPS
  IFS=' ' read -ra GROUPS <<< "$2"
  join_char=""
  res=""

  for G in "${GROUPS[@]}"; do
    res="$res$join_char$APIS_PKG/$G"
    join_char=","
  done

  echo "$res"
}

MODELS_SCHEMA="${MODELS_SCHEMA:-models-schema}"
CLIENT_GEN="${CLIENT_GEN:-client-gen}"
DEEPCOPY_GEN="${DEEPCOPY_GEN:-deepcopy-gen}"
LISTER_GEN="${LISTER_GEN:-lister-gen}"
INFORMER_GEN="${INFORMER_GEN:-informer-gen}"
DEFAULTER_GEN="${DEFAULTER_GEN:-defaulter-gen}"
CONVERSION_GEN="${CONVERSION_GEN:-conversion-gen}"
OPENAPI_GEN="${OPENAPI_GEN:-openapi-gen}"
APPLYCONFIGURATION_GEN="${APPLYCONFIGURATION_GEN:-applyconfiguration-gen}"

CLIENT_GROUPS="core"
CLIENT_VERSION_GROUPS="core:v1alpha1"
ALL_VERSION_GROUPS="$CLIENT_VERSION_GROUPS"

echo "${bold}Public types${normal}"

echo "Generating ${blue}deepcopy${normal}"
"$DEEPCOPY_GEN" \
  --go-header-file "$SCRIPT_DIR/boilerplate.go.txt" \
  --output-file zz_generated.deepcopy.go \
  "$(qualify-gvs "./api" "$ALL_VERSION_GROUPS")"

echo "Generating ${blue}openapi${normal}"
"$OPENAPI_GEN" \
  --go-header-file "$SCRIPT_DIR/boilerplate.go.txt" \
  --output-pkg "spheric.cloud/spheric/client-go/openapi" \
  --output-dir "./client-go/openapi" \
  --output-file zz_generated.openapi.go \
  --report-filename "./client-go/openapi/api_violations.report" \
  "k8s.io/apimachinery/pkg/apis/meta/v1" \
  "k8s.io/apimachinery/pkg/api/resource" \
  "k8s.io/apimachinery/pkg/runtime" \
  "k8s.io/apimachinery/pkg/version" \
  "$(qualify-gvs "./api" "$ALL_VERSION_GROUPS")"

echo "Generating ${blue}applyconfiguration${normal}"
"$APPLYCONFIGURATION_GEN" \
  --go-header-file "$SCRIPT_DIR/boilerplate.go.txt" \
  --output-pkg "spheric.cloud/spheric/client-go/applyconfigurations" \
  --output-dir "./client-go/applyconfigurations" \
  --openapi-schema <("$MODELS_SCHEMA" --openapi-package "spheric.cloud/spheric/client-go/openapi" --openapi-title "spheric") \
  "$(qualify-gvs "./api" "$ALL_VERSION_GROUPS")"

echo "Generating ${blue}client${normal}"
"$CLIENT_GEN" \
  --go-header-file "$SCRIPT_DIR/boilerplate.go.txt" \
  --output-pkg "spheric.cloud/spheric/client-go" \
  --output-dir "client-go" \
  --apply-configuration-package "spheric.cloud/spheric/client-go/applyconfigurations" \
  --clientset-name "spheric" \
  --input-base "spheric.cloud/spheric/api" \
  --input "core/v1alpha1"

echo "Generating ${blue}lister${normal}"
"$LISTER_GEN" \
  --go-header-file "$SCRIPT_DIR/boilerplate.go.txt" \
  --output-pkg "spheric.cloud/spheric/client-go/listers" \
  --output-dir "client-go/listers" \
  "$(qualify-gvs "./api" "$CLIENT_VERSION_GROUPS")" \

echo "Generating ${blue}informer${normal}"
"$INFORMER_GEN" \
  --go-header-file "$SCRIPT_DIR/boilerplate.go.txt" \
  --output-pkg "spheric.cloud/spheric/client-go/informers" \
  --output-dir "client-go/informers" \
  --listers-package "spheric.cloud/spheric/client-go/listers" \
  --versioned-clientset-package "spheric.cloud/spheric/client-go/spheric" \
  --single-directory \
  "$(qualify-gvs "./api" "$CLIENT_VERSION_GROUPS")"

echo "${bold}Internal types${normal}"

echo "Generating ${blue}deepcopy${normal}"
"$DEEPCOPY_GEN" \
  --go-header-file "$SCRIPT_DIR/boilerplate.go.txt" \
  --output-file zz_generated.deepcopy.go \
  "$(qualify-gs "./internal/apis" "$CLIENT_GROUPS")"

echo "Generating ${blue}defaulter${normal}"
"$DEFAULTER_GEN" \
  --go-header-file "$SCRIPT_DIR/boilerplate.go.txt" \
  --output-file "zz_generated.defaults.go" \
  "$(qualify-gvs "./internal/apis" "$CLIENT_VERSION_GROUPS")"

echo "Generating ${blue}conversion${normal}"
"$CONVERSION_GEN" \
  --go-header-file "$SCRIPT_DIR/boilerplate.go.txt" \
  --output-file "zz_generated.conversion.go" \
  "$(qualify-gs "spheric.cloud/spheric/internal/apis" "$CLIENT_GROUPS")" \
  "$(qualify-gvs "spheric.cloud/spheric/internal/apis" "$CLIENT_VERSION_GROUPS")"
