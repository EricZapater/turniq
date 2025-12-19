#!/bin/sh

# Default to empty if not set
: "${API_URL:=}"

# Recreate config file
echo "window.env = {" > /usr/share/nginx/html/env-config.js
echo "  API_URL: \"$API_URL\"" >> /usr/share/nginx/html/env-config.js
echo "};" >> /usr/share/nginx/html/env-config.js

# Execute the passed command (nginx)
exec "$@"
