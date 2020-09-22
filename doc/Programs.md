# Programs

  - what is a program?
  - how to switch between them
  - how can 2 programs run at once
  - programs can use bundled, eg: any self driving program will require collision avoidance
  which will run a separate program. a bundle can have priorities changed, depending on the 
  main program requirements.
  - priority based? two programs with the same priority cannot run at once!
      - `0` it will take control of the rc car completely
      and will not allow any other program to run
      - `1` to `9` will take control of the rc car completely, but will allow 
      programs with lower priority to run. the highest priority (lower number) will
      get the control of the vehicle. 
  
 
## 0 - ML Room Image Capture
Used for capturing images and tagging them with metadata required for the ML model
to categorize. 

  1. rotate to north
  2. for every 36deg (while rotating clockwise) take a photo
  3. save the photo in camera data dir and add metadata (in a separate csv file):
        - location_id (user provided at program startup)
        - time of day (collect once)
        - magnetometer data (collect every image from board)
        - wifi data (collect once from board)
        - bluetooth data (collect once from onboard bt)
  4. save the collected data in GCS to develop and train the model in google colab

## 0 - Self recharge
Travel (safely) to the charging station and auto plug. charging stations coming soon (tm)
        
## 1 - Web control
Turn it in a basic RC car with a webcam, controlled from a webinterface

## 1 - ROS Node
Run the car using ROS; make use of all nice UI tools from ROS.

## 1 - Basic autonomy 
Get from one room to another (or even more specific descriptions, eg: my desk, the aquarium, the ferrari :)

## 2 - Collision avoidance

  - the mcu will detect that the distance from TOF sensors is < threshold and will:
    1. stop the car and change a flag in the memory, so the brain can read; refuse to move 
    the motors until the flag is cleared
    2. just change a flag and let the brain decide what to do with the motors; the
    flag will be used only as a notification system
    3. a mix between `1` and `2` where there is a delay changing the flag and stopping 
    the motors; the brain might have enough time to override the motor stop
  - the decision of how to avoid the obstacle should be made by the brain
      
      
## 11 - Room detection
Categorize the room based on information gathered with `0 - ML Room Photo Capture`.

## 12 - Road / Sign recognition
Small traffic signs on a small road in the middle of the living room; experimenting on at smaller scale :)

## 12 - Face detection
  - maybe should be prio 11, model might be too large for the coral with both.
  
## 21 - Human presence
  - get a list of devices from the router and compare with know mac addresses
  - maybe bt?