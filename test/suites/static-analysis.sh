test_static_analysis_go() {
  if [ -n "${SKIP_STATIC:-}" ]; then
    echo "==> SKIP: Asked to skip static analysis"
    return
  fi

  (
    set -e

    cd ../

    ## Functions starting by empty line
    OUT=$(grep -Rrn --exclude-dir=vendor "^$" -B1 . | grep "func " | grep -v "}$" || true)
    if [ -n "${OUT}" ]; then
      printf "ERROR: Functions must not start with an empty line: \n%s" "${OUT}"
      false
    fi

    ## go vet, if it exists
    if go help vet >/dev/null 2>&1; then
      go vet ./...
    fi

    ## golint
    if which golint >/dev/null 2>&1; then
      golint -set_exit_status ./
    fi

    ## deadcode
    if which deadcode >/dev/null 2>&1; then
      OUT=$(deadcode ./ 2>&1 || true)
      if [ -n "${OUT}" ]; then
        echo "${OUT}" >&2
        false
      fi
    fi

    ## misspell
    if which misspell >/dev/null 2>&1; then
      OUT=$(misspell 2>/dev/null ./ | grep -Ev "^vendor/" || true)
      if [ -n "${OUT}" ]; then
        echo "Found some typos"
        echo "${OUT}"
        exit 1
      fi
    fi

    ## ineffassign
    if which ineffassign >/dev/null 2>&1; then
      ineffassign ./
    fi

    # go fmt
    gofmt -w -s ./

    git add -u :/
    git diff --exit-code
  )
}

test_static_analysis_shell() {
  if [ -n "${SKIP_STATIC:-}" ]; then
    echo "==> SKIP: Asked to skip static analysis"
    return
  fi

  (
    set -e

    cd ../

    # Shell static analysis
    if which shellcheck >/dev/null 2>&1; then
      shellcheck --shell sh test/*.sh test/includes/*.sh test/suites/*.sh
    else
      echo "shellcheck not found, shell static analysis disabled"
    fi

    ## Mixed tabs/spaces in scripts
    OUT=$(grep -Pr '\t' . | grep '\.sh:' || true)
    if [ -n "${OUT}" ]; then
      echo "ERROR: mixed tabs and spaces in script: ${OUT}"
      false
    fi

    ## Trailing whitespace in scripts
    OUT=$(grep -r " $" . | grep '\.sh:' || true)
    if [ -n "${OUT}" ]; then
      echo "ERROR: trailing whitespace in script: ${OUT}"
      false
    fi
  )
}
