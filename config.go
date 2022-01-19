package main

import "time"

// I know there are libraries to parse config files in yaml/toml/json etc. but I'm just keeping it simple
// with a few constants. In production this is obviously not a good solution.
//
// Regarding my personal style,
// 1. I am very distrustful of code with hardcoded values and tend to end up with quite detailed config files.
//    But I can obviously adapt my tendencies to the team's style.
// 2. If I don't spend time thinking of names that are clear yet concise I tend to err on the side of caution
//    and write rather long names. After all, I have never had trouble understanding a long name.

// Server settings
const PAGES_TO_FETCH = 1
const BASE_URL = "http://localhost:8099"
const RETRIES_PER_PAGE = 100 // Number of times to retry
const API_USERNAME = "api"
const API_PASSWORD = "secret"
const SECONDS_TO_RENEW_TOKEN = 2

// Connection pool settings
const ACTIVE_CONNECTION_TIMEOUT time.Duration = 1 // Timeout values are in seconds
var MAX_CONNECTIONS_PER_HOST = 10
var MAX_CONNECTIONS = 20 // Connections to any host. Not very relevant here.
var WORKERS = 10

// Other
const DEBUG = true
