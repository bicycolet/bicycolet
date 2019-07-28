test_check_unit_tests() {
    (
        set -e

        cd ../
        go test -v ./...
    )
}
