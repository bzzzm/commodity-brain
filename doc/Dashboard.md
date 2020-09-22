# Dashboard

Multiple apps (separate single-page apps) that will make use of the api, ws or mjpeg stream.
The program selector should be in the top bar and always changable (like project in gcp ui). 

## KPI
Uses the websockets to update information in a grafana style dashboard

## RC Driving
A dashboard that will be used for driving the car using the onboard camera and keyboard controls.

## Switches and knobs
Uses the API to switch different values from the database, restart routines, etc.

## ML Room Image Capture
Uses program [0 - ML Room Image Capture](Programs.md#), similar to driving, but has a UI made to create new
photos the ML process.