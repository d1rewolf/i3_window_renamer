i3_window_renamer simply listens in the background to i3 window title change events, renames 
window titles to include the X Window instance name. This is useful to me because I often have
the same application open multiple times (for example, multiple profiles of Firefox running), 
and the instance name lets me easily distinguish between them.

To build, go get and go build, and then "exec --no-startup-id /path/to/build/executable" in your i3 config.
