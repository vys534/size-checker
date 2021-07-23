# Size Checker

A program to gather basic size information about a directory and its contents, and send it to a Redis cluster. **Does not traverse subdirectories.**

Information that is pushed to the database includes:

- Total size (in bytes) of all files (`sc_total_size`)
- Number of files (`sc_file_count`)
- Time to complete operation (`sc_time_to_complete`)
- Last updated (`sc_last_updated` as unix time)

Note: the values don't expire.

lastly: dawg this shit licensed under da M I T so do fuck all with it idc lol