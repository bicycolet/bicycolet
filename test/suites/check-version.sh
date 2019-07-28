test_check_version() {
    # check the json output
    bicycolet_version=$(bicycolet version --format=json | jq ".client_version")
    bicycolet_major=$(echo "${bicycolet_version}" | xargs printf "%s" | cut -d. -f1)
    [ "${bicycolet_major}" = "0" ]
}
