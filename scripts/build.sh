#! /bin/bash

echo -e "Start running the script..."
cd ../

echo -e "Start building the app..."
wails build --clean -tags "production"

echo -e "End running the script!"
