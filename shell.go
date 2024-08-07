package main

import "fmt"

var shellFunction string = `
cr() {
    local new_dir
    new_dir=$(/usr/bin/crepo "$@")

    if [[ $? -ne 0 ]]; then
        return 1
    fi

    # return if the output starts with "No repositories found", or "Usage: cr <search_term>"
    if [[ $new_dir == "No repositories found"* ]] || [[ $new_dir =~ Usage:\ cr\ .* ]] || [[ $new_dir =~ cr\(\).* ]]; then
        echo "$new_dir"
            return 1
    fi

    if [[ -d "$new_dir" ]]; then
        cd "$new_dir" || echo "Failed to change directory to $new_dir"
    else
        echo "Invalid directory: $new_dir"
    fi
}
`

func printShellFunction() {
	fmt.Println(shellFunction)
}
